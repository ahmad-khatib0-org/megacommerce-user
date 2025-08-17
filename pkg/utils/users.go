package utils

import (
	"math/rand"
	"strings"
)

// RandomUserName generates a random username that matches the
// ValidUserNameRegex. The length is between minLen and maxLen (inclusive).
func RandomUserName(min, max int) string {
	const validUsernameChars = LowercaseLetters + UppercaseLetters + Numbers + "._-"
	if min < 1 || max < min {
		panic("invalid min/max length")
	}

	length := rand.Intn(max-min+1) + min
	username := make([]byte, length)
	for i := range username {
		username[i] = validUsernameChars[rand.Intn(len(validUsernameChars))]
	}

	return string(username)
}

// IsValidUsernameChars checks if the username characters are allowed
func IsValidUsernameChars(un string) bool {
	return ValidUserNameRegex.MatchString(un)
}

// SanitizeUnicode will remove undesirable Unicode characters from a string.
func SanitizeUnicode(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", "")
	return strings.Map(filterBlocklist, s)
}

// filterBlocklist returns `r` if it is not in the blocklist, otherwise drop (-1).
// Blocklist is taken from https://www.w3.org/TR/unicode-xml/#Charlist
func filterBlocklist(r rune) rune {
	const drop = -1

	switch r {
	case '\u0340', '\u0341': // clones of grave and acute; deprecated in Unicode
		return drop
	case '\u17A3', '\u17D3': // obsolete characters for Khmer; deprecated in Unicode
		return drop
	case '\u2028', '\u2029': // line and paragraph separator
		return drop
	case '\u202A', '\u202B', '\u202C', '\u202D', '\u202E': // BIDI embedding controls
		return drop
	case '\u206A', '\u206B': // activate/inhibit symmetric swapping; deprecated in Unicode
		return drop
	case '\u206C', '\u206D': // activate/inhibit Arabic form shaping; deprecated in Unicode
		return drop
	case '\u206E', '\u206F': // activate/inhibit national digit shapes; deprecated in Unicode
		return drop
	case '\uFFF9', '\uFFFA', '\uFFFB': // interlinear annotation characters
		return drop
	case '\uFEFF': // byte order mark
		return drop
	case '\uFFFC': // object replacement character
		return drop
	}

	// Scoping for musical notation
	if r >= 0x0001D173 && r <= 0x0001D17A {
		return drop
	}

	// Language tag code points
	if r >= 0x000E0000 && r <= 0x000E007F {
		return drop
	}

	return r
}
