package helpers

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/fingcloud/cli/api"
	ignore "github.com/sabhiram/go-gitignore"
)

func GetFiles(path string) ([]*api.FileInfo, error) {
	patterns := loadIgnorefiles(path)
	patterns = append(patterns, defaultIngnores...)

	ignore := ignore.CompileIgnoreLines(patterns...)

	files := make([]*api.FileInfo, 0)

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if ignore.MatchesPath(path) {
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
			Path:     path,
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

func Compress(files []*api.FileInfo, writer io.Writer) error {
	gzw := gzip.NewWriter(writer)
	tw := tar.NewWriter(gzw)

	for _, file := range files {
		info, err := os.Stat(file.Path)
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
			f, err := os.Open(file.Path)
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
