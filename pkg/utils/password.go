package utils

import (
	"errors"
	"fmt"
	"strings"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"golang.org/x/crypto/bcrypt"
)

type InvalidPassword struct {
	ID     string
	Err    string
	Params map[string]any
}

func (ip InvalidPassword) Error() string {
	return ip.Err
}

// IsValidPassword validate a password with a given ConfigPassword,
//
// idPrefix is optional, E,g if not passed, min length error will be "password.min_length"
func IsValidPassword(pass string, settings *pb.ConfigPassword, idPrefix string) *InvalidPassword {
	id := "password."
	isErr := false
	var params map[string]any

	if idPrefix != "" {
		id = idPrefix
	}

	if len(pass) < int(settings.GetMinimumLength()) {
		id += "min_length"
		isErr = true
		params = map[string]any{"Min": settings.GetMinimumLength()}
	} else if len(pass) > int(settings.GetMaximumLength()) {
		id += "max_length"
		isErr = true
		params = map[string]any{"Max": settings.GetMaximumLength()}
	} else if settings.GetLowercase() && !strings.ContainsAny(pass, LowercaseLetters) {
		id += "lowercase"
		isErr = true
	} else if settings.GetUppercase() && !strings.ContainsAny(pass, UppercaseLetters) {
		id += "uppercase"
		isErr = true
	} else if settings.GetNumber() && !strings.ContainsAny(pass, Numbers) {
		id += "numbers"
		isErr = true
	} else if settings.GetSymbol() && !strings.ContainsAny(pass, Symbols) {
		id += "symbols"
		isErr = true
	}

	if isErr {
		return &InvalidPassword{
			ID:     id,
			Err:    fmt.Sprintf("invalid password: %s, err: %s ", pass, id),
			Params: params,
		}
	}

	return nil
}

// PasswordHash generates a hash using the bcrypt.GenerateFromPassword
func PasswordHash(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// PasswordCheck returns error if the hash does not match a given password
func PasswordCheck(hash string, password string) error {
	if hash == "" || password == "" {
		return errors.New("empty password or hash")
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
