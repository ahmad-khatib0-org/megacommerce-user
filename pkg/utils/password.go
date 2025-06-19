package utils

import (
	"fmt"
	"strings"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"golang.org/x/crypto/bcrypt"
)

type InvalidPassword struct {
	Id  string
	Err string
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

	if idPrefix != "" {
		id = idPrefix
	}

	if len(pass) < int(settings.GetMinimumLength()) {
		id += "min_length"
		isErr = true
	}

	if len(pass) < int(settings.GetMaximumLenght()) {
		id += "min_length"
		isErr = true
	}

	if settings.GetLowercase() && !strings.ContainsAny(pass, LowercaseLetters) {
		id += "lowercase"
		isErr = true
	}

	if settings.GetUppercase() && !strings.ContainsAny(pass, UppercaseLetters) {
		id += "uppercase"
		isErr = true
	}

	if settings.GetNumber() && !strings.ContainsAny(pass, Numbers) {
		id += "numbers"
		isErr = true
	}

	if settings.GetSymbol() && !strings.ContainsAny(pass, Symbols) {
		id += "symbols"
		isErr = true
	}

	if isErr {
		return &InvalidPassword{
			Id:  id,
			Err: fmt.Sprintf("invalid password: %s, err: %s ", pass, id),
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
