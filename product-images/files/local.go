package files

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/xerrors"
)

type Local struct {
	maxFileSize int
	basePath    string
}

func NewLocal(basePath string, maxSize int) (*Local, error) {
	p, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}
	return &Local{maxFileSize: maxSize, basePath: p}, nil
}

func (l *Local) fullPath(p string) string {
	return path.Join(l.basePath, p)
}

func (l *Local) Get(path string) (*os.File, error) {
	fp := l.fullPath(path)
	f, err := os.Open(fp)
	if err != nil {
		return nil, xerrors.Errorf("Unable to open file: %w", err)
	}
	return f, nil
}

func (l *Local) Save(path string, content io.Reader) error {
	fp := l.fullPath(path)
	d := filepath.Dir(fp)
	err := os.MkdirAll(d, os.ModePerm)
	if err != nil {
		return xerrors.Errorf("Unable to create directory: %w", err)
	}

	_, err = os.Stat(fp)
	if err == nil || !os.IsNotExist(err) {
		if err := os.Remove(fp); err != nil {
			return xerrors.Errorf("Unable to delete file: %w", err)
		}

	}

	f, err := os.Create(fp)
	if err != nil {
		return xerrors.Errorf("Unable to create file: %w", err)
	}
	defer f.Close()

	b := make([]byte, 0, 1024)
	n := 0
	numberOfRead := 0
	var readErr error
	var writeErr error

	for n, readErr = content.Read(b[len(b):cap(b)]); numberOfRead <= l.maxFileSize; n, readErr = content.Read(b[len(b):cap(b)]) {
		if readErr != nil && readErr != io.EOF {
			break
		}
		numberOfRead += n
		_, writeErr := f.Write(b[:len(b)+n])
		if writeErr != nil {
			break
		}
		if readErr == io.EOF {
			break
		}
		b = b[:0]
	}
	if numberOfRead > l.maxFileSize {
		return xerrors.Errorf("Size of received file more than allowed maximum - %d", l.maxFileSize)
	}
	if readErr != nil && readErr != io.EOF {
		err = readErr
	} else if writeErr != nil {
		err = writeErr
	}
	if err != nil {
		return xerrors.Errorf("Unable to save file: %w", err)
	}

	return nil
}
