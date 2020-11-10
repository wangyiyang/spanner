package file

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func addTarFile(path, name string, tw *tar.Writer, twBuf *bufio.Writer) error {
	fi, err := os.Lstat(path)
	if err != nil {
		return err
	}
	link := ""
	if fi.Mode()&os.ModeSymlink != 0 {
		if link, err = os.Readlink(path); err != nil {
			return err
		}
	}
	hdr, err := tar.FileInfoHeader(fi, link)
	if err != nil {
		return err
	}
	if fi.IsDir() && !strings.HasSuffix(name, "/") {
		name = name + "/"
	}

	hdr.Name = name
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if hdr.Typeflag == tar.TypeReg {
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		twBuf.Reset(tw)
		_, err = io.Copy(twBuf, file)
		file.Close()
		if err != nil {
			return err
		}
		err = twBuf.Flush()
		if err != nil {
			return err
		}
		twBuf.Reset(nil)
	}
	return nil
}

func TarWithOptions(srcPath string, options *TarOptions) (io.ReadCloser, error) {
	pipeReader, pipeWriter := io.Pipe()

	compressWriter, err := CompressStream(pipeWriter, options.Compression)
	if err != nil {
		return nil, err
	}

	tw := tar.NewWriter(compressWriter)

	go func() {
		// In general we logging errors here but ignore them because
		// during e.g. a diff operation the container can continue
		// mutating the filesystem and we can see transient errors
		// from this

		if options.Includes == nil {
			options.Includes = []string{"."}
		}

		twBuf := bufio.NewWriterSize(nil, twBufSize)

		for _, include := range options.Includes {
			filepath.Walk(filepath.Join(srcPath, include), func(filePath string, f os.FileInfo, err error) error {
				if err != nil {
					return nil
				}

				relFilePath, err := filepath.Rel(srcPath, filePath)
				if err != nil {
					return nil
				}

				skip, err := matches(relFilePath, options.Excludes)
				if err != nil {
					return err
				}

				if skip {
					if f.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}

				if err := addTarFile(filePath, relFilePath, tw, twBuf); err != nil {
					return err
				}
				return nil
			})
		}

		// Make sure to check the error on Close.
		if err := tw.Close(); err != nil {
			log.Print(err)
			return

		}
		if err := compressWriter.Close(); err != nil {
			log.Print(err)
			return
		}
		if err := pipeWriter.Close(); err != nil {
			log.Print(err)
			return
		}
	}()

	return pipeReader, nil
}

func BuildTar(root string) (io.ReadCloser, error) {
	filename := path.Join(root, "Dockerfile")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("no Dockerfile found in %s", root)
	}
	var excludes []string
	ignore, err := ioutil.ReadFile(path.Join(root, ".dockerignore"))
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("Error reading .dockerignore: '%s'", err)
	}
	for _, pattern := range strings.Split(string(ignore), "\n") {
		ok, err := filepath.Match(pattern, "Dockerfile")
		if err != nil {
			return nil, fmt.Errorf("Bad .dockerignore pattern: '%s', error: %s", pattern, err)
		}
		if ok {
			return nil, fmt.Errorf("Dockerfile was excluded by .dockerignore pattern '%s'", pattern)
		}
		excludes = append(excludes, pattern)
	}

	options := &TarOptions{
		Compression: Uncompressed,
		Excludes:    excludes,
	}

	context, err := TarWithOptions(root, options)

	if err != nil {
		return nil, err
	}
	return context, nil
}
