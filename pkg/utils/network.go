package utils

import (
	"context"
	"slices"

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
