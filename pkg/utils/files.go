package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FindDir walks up directories until it finds dir.
func FindDir(dir string) (string, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		path := filepath.Join(rootDir, dir)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path, nil
		}
		parent := filepath.Dir(rootDir)
		if parent == rootDir {
			return "", fmt.Errorf("could not find dir: %s ", dir)
		}
		rootDir = parent
	}
}

// CleanBase64 strips prefixes and commas
func CleanBase64(input string) string {
	re := regexp.MustCompile(`^data:image\/[a-zA-Z]+;base64,`)
	cleaned := re.ReplaceAllString(input, "")

	// Remove commas (not valid in base64)
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	return strings.TrimSpace(cleaned)
}
