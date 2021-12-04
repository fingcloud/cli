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
	defaultIgnores = []string{".git", "node_modules", "*.*~", "bower_components", ".*"}
)

func readIgnoreFile(fpath string) ([]string, error) {
	bytes, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(bytes), "\n")

	patterns := funk.FilterString(lines, func(s string) bool {
		if strings.HasPrefix(s, "#") {
			return false
		}
		if strings.TrimSpace(s) == "" {
			return false
		}
		return true
	})

	return patterns, nil
}

func loadIgnorefiles(projectPath, path string) []string {

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

	ignorePatterns = funk.Map(ignorePatterns, func(pattern string) string {
		pattern = strings.ReplaceAll(pattern, "\\", "/")

		suffux := ""
		if strings.HasSuffix(pattern, "/") {
			suffux = "/"
		}

		absolute := ""
		if strings.HasPrefix(pattern, "!") {
			if string(pattern[1]) == "/" {
				absolute = "/"
			}
			return "!" + absolute + filepath.Join(projectPath, path, pattern[1:]) + suffux
		}

		if string(pattern[0]) == "/" {
			absolute = "/"
		}
		return absolute + filepath.Join(projectPath, path, pattern) + suffux
	}).([]string)

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
