// Package utils provide common utils
package utils

import (
	"strings"

	"github.com/oklog/ulid/v2"
)

func NewIDPointer() *string {
	id := ulid.Make().String()
	return &id
}

func NewID() string {
	return ulid.Make().String()
}

func NewPointer[T any](v T) *T {
	return &v
}

// GetMetadataValue convert the []string metadata value to a map[string]string
func GetMetadataValue(value []string) map[string]string {
	if len(value) == 0 {
		return map[string]string{}
	}

	result := make(map[string]string, len(value))
	for _, v := range value {
		parts := strings.SplitN(v, ":", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}
