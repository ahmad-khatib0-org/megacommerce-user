package utils

import "github.com/oklog/ulid/v2"

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
