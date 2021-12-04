package fileutils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fingcloud/cli/pkg/api"
	ignore "github.com/sabhiram/go-gitignore"
)

func GetFiles(projectPath string) ([]*api.FileInfo, error) {
	patterns := make([]string, 0)
	patterns = append(patterns, defaultIgnores...)

	cachedDirs := make(map[string]bool, 0)
	ig := ignore.CompileIgnoreLines(patterns...)

	files := make([]*api.FileInfo, 0)

	err := filepath.Walk(projectPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// has ignorefile ? add to patterns
		if info.IsDir() {
			if _, ok := cachedDirs[path]; !ok {
				dirPatterns := loadIgnorefiles(projectPath, path)
				patterns = append(patterns, dirPatterns...)
				ig = ignore.CompileIgnoreLines(patterns...)
				cachedDirs[path] = true
			}
		}

		if ig.MatchesPath(path) {
			return nil
		}

		filePath, _ := filepath.Rel(projectPath, path)
		if filePath == "" || filePath == "." {
			return nil
		}

		var checksum string
		if !info.IsDir() {
			checksum, err = fileChecksum(path)
			if err != nil {
				return fmt.Errorf("can't resolve checksum: %v", err)
			}
		}

		files = append(files, &api.FileInfo{
			Path:     filePath,
			Size:     info.Size(),
			Dir:      info.IsDir(),
			Checksum: checksum,
			Mode:     fileMode(info),
		})

		return nil
	})

	return files, err
}

func fileMode(info fs.FileInfo) fs.FileMode {
	if info.IsDir() || info.Mode()&0111 != 0 {
		return fs.FileMode(0755)
	}

	return fs.FileMode(0644)
}

func isExecutable(mode fs.FileMode) bool {
	return mode&0111 != 0
}

func Compress(projectPath string, files []*api.FileInfo, writer io.Writer) error {
	gzw := gzip.NewWriter(writer)
	tw := tar.NewWriter(gzw)

	for _, file := range files {
		path := filepath.Join(projectPath, file.Path)
		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, file.Path)
		if err != nil {
			return err
		}

		// must provide real name
		// (see https://golang.org/src/archive/tar/common.go?#L626)
		header.Name = filepath.ToSlash(file.Path)
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			if _, err := io.Copy(tw, f); err != nil {
				return err
			}
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}

	if err := gzw.Close(); err != nil {
		return err
	}

	return nil
}

func Decompress(dst string, reader io.Reader) error {
	zr, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	tr := tar.NewReader(zr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(dst, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			f.Close()
		}
	}

	return nil
}
