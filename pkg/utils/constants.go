package utils

import "regexp"

const (
	LowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	UppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Numbers          = "0123456789"
	Symbols          = " !\"\\#$%&'()*+,-./:;<=>?@[]^_`|~"
	MaxPropSizeBytes = 1024 * 1024
)

var ValidUserNameRegex = regexp.MustCompile(`^[a-zA-Z0-9\.\-_]+$`)
