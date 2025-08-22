package utils

import "slices"

func ProcessAcceptedLanguage(header string, availableLangs []string, defaultLang string) string {
	if slices.Contains(availableLangs, header) {
		return header
	}
	return defaultLang
}
