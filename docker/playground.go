package docker

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func PlayCode(code string) (string, error) {
	_, deferFunc, err := setupTempDir()
	if err != nil {
		return "", err
	}

	defer deferFunc()

	err = ioutil.WriteFile("code.piece", []byte(code), 0666)
	if err != nil {
		return "", err
	}

	return PlayFile("code.piece")
}

func PlayFile(file string) (string, error) {
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	cmd := exec.Command(
		"docker", "run", "--rm",
		"-v", path+":/code/piece",
		"-t", "flwos/playground",
		"--file", "/code/piece",
	)
	output, err := cmd.Output()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}

func setupTempDir() (string, func(), error) {
	tmpdir, err := ioutil.TempDir("", "flw-playground-")
	if err != nil {
		return "", nil, err
	}

	oldwd, err := os.Getwd()
	if err != nil {
		os.RemoveAll(tmpdir)
		return "", nil, err
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
