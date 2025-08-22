package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/fatih/color"
	"github.com/stretchr/testify/mock"
)

func ReadTestFile(name string) ([]byte, error) {
	testsDir, err := FindDir("test_files")
	if err != nil {
		return nil, err
	}

	return os.ReadFile(filepath.Join(testsDir, name))
}

func ReadTestImage(name string) (string, error) {
	data, err := ReadTestFile(name)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// WithMockDebug wraps a matcher with logging. It logs every argument
// received and whether it matched.
func WithMockDebug[T any](label string, matcher func(T) bool) any {
	logged := false
	return mock.MatchedBy(func(arg T) bool {
		matched := matcher(arg)

		if !logged {
			// Use colored logging for clarity
			status := color.New(color.FgGreen).Sprint("MATCH")
			if !matched {
				status = color.New(color.FgRed).Sprint("NO MATCH")
			}

			fmt.Printf("\n🧪 [%s] %s\n", label, status)
			fmt.Printf("    Type:  %s\n", reflect.TypeOf(arg))
			fmt.Printf("    Value: %#v\n\n", arg)

			logged = true
		}

		return matched
	})
}
