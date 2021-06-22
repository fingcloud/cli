package fileutils

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/thoas/go-funk"
)

var (
	defaultIngnores = []string{".git", "node_modules", "*.*~", "bower_components", ".*"}
)

func readIgnoreFile(fpath string) ([]string, error) {
	bytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bytes), "\n")

	patterns := funk.FilterString(lines, func(s string) bool {
		return !strings.HasPrefix(s, "#") || strings.TrimSpace(s) != ""
	})

	return patterns, nil
}

func loadIgnorefiles(path string) []string {

	ignorePatterns := make([]string, 0)

	if patterns, err := readIgnoreFile(filepath.Join(path, ".gitignore")); err == nil {
		ignorePatterns = append(ignorePatterns, patterns...)
	}

	if patterns, err := readIgnoreFile(filepath.Join(path, ".dockerignore")); err == nil {
		ignorePatterns = append(ignorePatterns, patterns...)
	}

	if patterns, err := readIgnoreFile(filepath.Join(path, ".fingignore")); err == nil {
		ignorePatterns = append(ignorePatterns, patterns...)
	}

	return ignorePatterns
}

func fileChecksum(fpath string) (string, error) {
	hasher := sha256.New()

	bytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return "", err
	}

	hasher.Write(bytes)

	checksum := hex.EncodeToString(hasher.Sum(nil))

	return checksum, nil
}
