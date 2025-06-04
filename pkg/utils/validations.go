package utils

import (
	"net/mail"
	"strings"
)

func IsValidEmail(email string) bool {
	if add, err := mail.ParseAddress(email); err != nil {
		return false
	} else if add.Name != "" {
		return false
	}

	// mail.ParseAddress accepts quoted strings for the address which can lead to sending
	// to the wrong email address check for multiple '@' symbols and invalidate
	if strings.Count(email, "@") > 1 {
		return false
	}

	return true
}
