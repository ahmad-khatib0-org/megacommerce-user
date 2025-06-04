package utils

import (
	"fmt"
	"strings"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"google.golang.org/grpc/codes"
)

func IsValidPassword(pass string, settings *pb.ConfigPassword, where, idPrefix string) *AppError {
	id := "user.create.password."
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
		return NewAppError(where, id, nil, fmt.Sprintf("invalid password: %s, err: %s ", pass, id), int(codes.Internal))
	}

	return nil
}
