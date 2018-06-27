package docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var supportedLang map[string]string

func init() {
	supportedLang = map[string]string{
		"go":     "go",
		"golang": "go",
	}
}

func PlayCode(lang, code string) (string, error) {
	if lang = supportedLang[lang]; lang == "" {
		return "", fmt.Errorf("Unsupported language: %s", lang)
	}

	_, deferFunc, err := setupTempDir()
	if err != nil {
		return "", err
	}

	defer deferFunc()

	err = ioutil.WriteFile("code.piece", []byte(code), 0666)
	if err != nil {
		return "", err
	}

	return play(lang, "code.piece", true)
}

func PlayFile(lang, file string) (string, error) {
	return play(lang, file, false)
}

func play(lang, file string, isCodePice bool) (string, error) {
	if lang = supportedLang[lang]; lang == "" {
		return "", fmt.Errorf("Unsupported language: %s", lang)
	}

	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	args := []string{
		"run", "--rm",
	}

	if isCodePice {
		args = append(args,
			"-v", path+":/code/piece",
			"-t", "flwos/playground",
			"--lang", lang,
			"--code", "/code/piece",
		)
	} else {
		args = append(args,
			"-v", path+":/code/main.go",
			"-t", "flwos/playground",
			"--lang", lang,
			"--file", "/code/main.go",
		)
	}

	fmt.Printf("RUN docker %v\n", args)
	cmd := exec.Command("docker", args...)
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
