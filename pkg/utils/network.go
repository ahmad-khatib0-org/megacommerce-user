package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"
	"syscall"

	"google.golang.org/grpc/metadata"
)

func ProcessAcceptedLanguage(header string, availableLangs []string, defaultLang string) string {
	if slices.Contains(availableLangs, header) {
		return header
	}
	return defaultLang
}

func GetAcceptedLanguageFromGrpcCtx(ctx context.Context, md metadata.MD, availableLangs []string, defaultLang string) string {
	lang := defaultLang
	if al := md.Get("accept-language"); len(al) > 0 {
		lang = ProcessAcceptedLanguage(al[0], availableLangs, defaultLang)
	}

	return lang
}

type ErrorSeverity int

const (
	SeverityRecoverable ErrorSeverity = iota
	SeverityFatal
	SeverityUnknown
)

func ClassifyHTTPServerError(err error) ErrorSeverity {
	if err == nil {
		return SeverityRecoverable
	}

	// Check for graceful shutdown first
	if errors.Is(err, http.ErrServerClosed) {
		return SeverityRecoverable
	}

	// 1. Net OpError classification
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		switch opErr.Op {
		case "listen", "bind":
			return SeverityFatal // Startup errors are fatal

		case "read", "write", "accept":
			// Connection-level errors - usually recoverable
			if isConnectionErrorFatal(opErr) {
				return SeverityFatal
			}
			return SeverityRecoverable

		case "dial":
			// Outbound connection errors - usually recoverable
			return SeverityRecoverable
		}
	}

	// 2. Syscall error classification (wrapped in os.SyscallError)
	var syscallErr *os.SyscallError
	if errors.As(err, &syscallErr) {
		return classifySyscallError(syscallErr.Err)
	}

	// 3. Direct syscall errors (unwrapped syscall.EXXX errors)
	if severity := classifyDirectSyscallError(err); severity != SeverityUnknown {
		return severity
	}

	// 4. TLS errors (usually fatal)
	var tlsErr *tls.CertificateVerificationError
	if errors.As(err, &tlsErr) {
		return SeverityFatal
	}

	// 5. DNS errors (usually recoverable)
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		if dnsErr.IsNotFound {
			return SeverityFatal // Unknown host is fatal for servers
		}
		return SeverityRecoverable // Temporary DNS issues
	}

	// 6. AddrError (usually fatal)
	var addrErr *net.AddrError
	if errors.As(err, &addrErr) {
		return SeverityFatal
	}

	// 7. Invalid host/port errors
	if isInvalidAddressError(err) {
		return SeverityFatal
	}

	// 8. File system errors for TLS certificates
	if isCertificateFileError(err) {
		return SeverityFatal
	}

	// 9. String pattern matching (fallback)
	return classifyByStringPattern(err)
}

func isConnectionErrorFatal(opErr *net.OpError) bool {
	// Check if this connection error is actually fatal
	var syscallErr *os.SyscallError
	if errors.As(opErr.Err, &syscallErr) {
		// Recoverable connection errors
		switch {
		case errors.Is(syscallErr.Err, syscall.ECONNABORTED),
			errors.Is(syscallErr.Err, syscall.ECONNRESET),
			errors.Is(syscallErr.Err, syscall.EPIPE),
			errors.Is(syscallErr.Err, syscall.EAGAIN),
			errors.Is(syscallErr.Err, syscall.EWOULDBLOCK),
			errors.Is(syscallErr.Err, syscall.EINTR),
			errors.Is(syscallErr.Err, syscall.ETIMEDOUT):
			return false // These are recoverable
		}

		// Fatal connection errors
		switch {
		case errors.Is(syscallErr.Err, syscall.EBADF), // Bad file descriptor
			errors.Is(syscallErr.Err, syscall.ENOTSOCK),     // Not a socket
			errors.Is(syscallErr.Err, syscall.EPROTO),       // Protocol error
			errors.Is(syscallErr.Err, syscall.ENOPROTOOPT),  // Protocol not available
			errors.Is(syscallErr.Err, syscall.EAFNOSUPPORT): // Address family not supported
			return true // These are fatal
		}
	}

	// Default: Assume connection errors are recoverable
	// (Better to keep serving than crash on unknown errors)
	return false
}

func classifySyscallError(err error) ErrorSeverity {
	if err == nil {
		return SeverityUnknown
	}

	switch {
	// Fatal syscall errors
	case errors.Is(err, syscall.EADDRINUSE): // Address already in use
		return SeverityFatal
	case errors.Is(err, syscall.EACCES): // Permission denied
		return SeverityFatal
	case errors.Is(err, syscall.EAFNOSUPPORT): // Address family not supported
		return SeverityFatal
	case errors.Is(err, syscall.EINVAL): // Invalid argument
		return SeverityFatal
	case errors.Is(err, syscall.ENFILE): // File table overflow
		return SeverityFatal
	case errors.Is(err, syscall.EMFILE): // Too many open files
		return SeverityFatal
	case errors.Is(err, syscall.ENOMEM): // Out of memory
		return SeverityFatal
	case errors.Is(err, syscall.ENOBUFS): // No buffer space available
		return SeverityFatal
	case errors.Is(err, syscall.EPROTONOSUPPORT): // Protocol not supported
		return SeverityFatal
	case errors.Is(err, syscall.ENOPROTOOPT): // Protocol option not available
		return SeverityFatal

	// Recoverable syscall errors
	case errors.Is(err, syscall.ECONNABORTED): // Connection aborted
		return SeverityRecoverable
	case errors.Is(err, syscall.ECONNRESET): // Connection reset by peer
		return SeverityRecoverable
	case errors.Is(err, syscall.EPIPE): // Broken pipe
		return SeverityRecoverable
	case errors.Is(err, syscall.EAGAIN): // Resource temporarily unavailable
		return SeverityRecoverable
	case errors.Is(err, syscall.EWOULDBLOCK): // Operation would block
		return SeverityRecoverable
	case errors.Is(err, syscall.EINTR): // Interrupted system call
		return SeverityRecoverable
	case errors.Is(err, syscall.ETIMEDOUT): // Connection timed out
		return SeverityRecoverable
	case errors.Is(err, syscall.ECONNREFUSED): // Connection refused
		return SeverityRecoverable
	}

	return SeverityUnknown
}

func classifyDirectSyscallError(err error) ErrorSeverity {
	// Check if it's a direct syscall error (not wrapped)
	switch err {
	case
		syscall.EADDRINUSE,
		syscall.EACCES,
		syscall.EAFNOSUPPORT,
		syscall.EINVAL,
		syscall.ENFILE,
		syscall.EMFILE,
		syscall.ENOMEM,
		syscall.ENOBUFS,
		syscall.EPROTONOSUPPORT,
		syscall.ENOPROTOOPT:
		return SeverityFatal

	case
		syscall.ECONNABORTED,
		syscall.ECONNRESET,
		syscall.EPIPE,
		syscall.EWOULDBLOCK,
		syscall.EINTR,
		syscall.ETIMEDOUT,
		syscall.ECONNREFUSED:
		return SeverityRecoverable
	}

	return SeverityUnknown
}

func isInvalidAddressError(err error) bool {
	errorMsg := strings.ToLower(err.Error())
	invalidPatterns := []string{
		"missing port in address",
		"too many colons in address",
		"invalid port",
		"invalid address",
		"unknown port",
		"missing address",
	}

	for _, pattern := range invalidPatterns {
		if strings.Contains(errorMsg, pattern) {
			return true
		}
	}
	return false
}

func isCertificateFileError(err error) bool {
	errorMsg := strings.ToLower(err.Error())
	certPatterns := []string{
		"tls:",
		"x509:",
		"certificate",
		"private key",
		"pem",
		"failed to load",
		"no such file",
		"permission denied",
		"read error",
	}

	for _, pattern := range certPatterns {
		if strings.Contains(errorMsg, pattern) {
			return true
		}
	}
	return false
}

func classifyByStringPattern(err error) ErrorSeverity {
	errorMsg := strings.ToLower(err.Error())

	// Fatal error patterns
	fatalPatterns := []string{
		"address already in use",
		"permission denied",
		"bind:",
		"listen:",
		"invalid configuration",
		"required environment variable",
		"missing port",
		"unknown port",
		"tls handshake failure",
		"certificate",
		"private key",
		"no such file",
		"protocol not available",
		"address family not supported",
		"too many open files",
		"out of memory",
		"no buffer space",
	}

	for _, pattern := range fatalPatterns {
		if strings.Contains(errorMsg, pattern) {
			return SeverityFatal
		}
	}

	// Recoverable error patterns
	recoverablePatterns := []string{
		"connection reset",
		"broken pipe",
		"connection aborted",
		"timeout",
		"eof",
		"temporary failure",
		"network is unreachable",
		"no route to host",
		"connection refused",
		"operation timed out",
	}

	for _, pattern := range recoverablePatterns {
		if strings.Contains(errorMsg, pattern) {
			return SeverityRecoverable
		}
	}

	return SeverityUnknown
}

// IsUnrecoverableHTTPServerError check if an HTTP server error is recoverable or not
func IsUnrecoverableHTTPServerError(err error) bool {
	return ClassifyHTTPServerError(err) == SeverityFatal
}
