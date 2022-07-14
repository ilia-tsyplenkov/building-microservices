package files

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
)

func setupLocal(t *testing.T) (*Local, string, func()) {
	d, err := ioutil.TempDir("", "files")
	if err != nil {
		t.Fatal(err)
	}
	l, err := NewLocal(d, 100)
	if err != nil {
		t.Fatal(err)
	}

	return l, d, func() { os.RemoveAll(d) }
}

func TestSaveContentOfReader(t *testing.T) {
	savePath := "/1/test.txt"
	fileContent := "Hello World"
	l, dir, cleanup := setupLocal(t)
	defer cleanup()
	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContent)))
	assert.NoError(t, err)

	// check file creating
	f, err := os.Open(path.Join(dir, savePath))
	assert.NoError(t, err)

	defer f.Close()
	d, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, string(d))

}

func TestSaveContentOfReaderAndReadIt(t *testing.T) {
	savePath := "/1/test.txt"
	fileContent := "Hello World"
	l, _, cleanup := setupLocal(t)
	defer cleanup()
	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContent)))
	assert.NoError(t, err)

	r, err := l.Get(savePath)
	assert.NoError(t, err)
	defer r.Close()
	d, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, string(d))

}

func TestSaveLargeFileReturnsError(t *testing.T) {
	savePath := "/1/test.txt"
	fileContent := "Hello World, a very very very very long file."
	l, _, cleanup := setupLocal(t)
	l.maxFileSize = len(fileContent) - 3
	defer cleanup()
	err := l.Save(savePath, bytes.NewBuffer([]byte(fileContent)))
	assert.Error(t, err)
}
