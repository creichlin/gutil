package io

import (
	"io/ioutil"
	"os"
	"path"
)

// WriteFileAtomic creates a temp file and renames it to to
// desired filename using a series of operations:
// - create temp file
// - write
// - sync
// - close
// - chmod
// - rename
// if in any of the operations after create an error occures
// the tempFile will be deleted (if it works) and the error returned
//
// this will result in an atomically written file, either it's the old one or
// the new content but not half of the change
func WriteFileAtomic(filename string, data []byte, perm os.FileMode) (err error) {
	dir, name := path.Split(filename)
	tmpFile, err := ioutil.TempFile(dir, name)
	if err != nil {
		return
	}

	// if an error happens, try to remove the created tmpFile
	defer func() {
		if err != nil {
			os.Remove(tmpFile.Name())
		}
	}()

	if _, err = tmpFile.Write(data); err != nil {
		return
	}

	if err = tmpFile.Sync(); err != nil {
		return
	}

	if err = tmpFile.Close(); err != nil {
		return
	}

	if err = os.Chmod(tmpFile.Name(), perm); err != nil {
		return
	}

	return os.Rename(tmpFile.Name(), filename)
}
