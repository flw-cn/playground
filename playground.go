package playground

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type PlayError struct {
	err    error
	output string
}

func (e *PlayError) Error() string {
	return e.err.Error()
}

func (e *PlayError) Output() string {
	return e.output
}

func setupTempDir(file ...string) (string, func(), error) {
	tmpdir, err := ioutil.TempDir("", "flw-playground-")
	if err != nil {
		return "", nil, err
	}

	oldwd, err := os.Getwd()
	if err != nil {
		os.RemoveAll(tmpdir)
		return "", nil, err
	}

	for _, f := range file {
		os.Rename(f, filepath.Join(tmpdir, filepath.Base(f)))
	}

	err = os.Chdir(tmpdir)
	if err != nil {
		return "", nil, err
	}

	return tmpdir, func() {
		os.RemoveAll(tmpdir)
		os.Chdir(oldwd)
	}, nil
}
