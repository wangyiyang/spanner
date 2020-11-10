package file

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"path"
	"path/filepath"
	"strings"
)

const twBufSize = 32 * 1024

type (
	Archive       io.ReadCloser
	ArchiveReader io.Reader
	Compression   int
	TarOptions    struct {
		Includes    []string
		Excludes    []string
		Compression Compression
		NoLchown    bool
	}
)

const (
	Uncompressed Compression = iota
	Bzip2
	Gzip
	Xz
)

func (compression *Compression) Extension() string {
	switch *compression {
	case Uncompressed:
		return "tar"
	case Bzip2:
		return "tar.bz2"
	case Gzip:
		return "tar.gz"
	case Xz:
		return "tar.xz"
	}
	return ""
}

type nopWriteCloser struct {
	io.Writer
}

func (w *nopWriteCloser) Close() error { return nil }

func NopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{w}
}

func CompressStream(dest io.WriteCloser, compression Compression) (io.WriteCloser, error) {
	switch compression {
	case Uncompressed:
		return NopWriteCloser(dest), nil
	case Gzip:
		return gzip.NewWriter(dest), nil
	case Bzip2, Xz:
		// archive/bzip2 does not support writing, and there is no xz support at all
		// However, this is not a problem as docker only currently generates gzipped tars
		return nil, fmt.Errorf("Unsupported compression format %s", (&compression).Extension())
	default:
		return nil, fmt.Errorf("Unsupported compression format %s", (&compression).Extension())
	}
}

func GetFileName(filePath string) (name string) {
	filenameWithSuffix := path.Base(filePath)
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix)
	name = strings.TrimSuffix(filenameWithSuffix, fileSuffix)
	return
}

// Matches returns true if relFilePath matches any of the patterns
func matches(relFilePath string, patterns []string) (bool, error) {
	for _, exclude := range patterns {
		matched, err := filepath.Match(exclude, relFilePath)
		if err != nil {
			log.Fatalf("Error matching: %s (pattern: %s)", relFilePath, exclude)
			return false, err
		}
		if matched {
			if filepath.Clean(relFilePath) == "." {
				log.Fatalf("Can't exclude whole path, excluding pattern: %s", exclude)
				continue
			}
			log.Fatalf("Skipping excluded path: %s", relFilePath)
			return true, nil
		}
	}
	return false, nil
}
