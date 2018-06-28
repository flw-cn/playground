package docker

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var supportedLang map[string]string

var codeVolume struct {
	dockerPath string
	hostPath   string
}

func init() {
	supportedLang = map[string]string{
		"go":     "go",
		"golang": "go",
	}

	codeVolume.dockerPath = "/"
	codeVolume.hostPath = "/"
}

func Boarding(dockerPath, hostPath string) error {
	if !filepath.IsAbs(dockerPath) {
		return fmt.Errorf("%s is not an absolute path", dockerPath)
	}

	if !filepath.IsAbs(hostPath) {
		return fmt.Errorf("%s is not an absolute path", hostPath)
	}

	codeVolume.dockerPath = dockerPath
	codeVolume.hostPath = hostPath

	return nil
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

	path, err := boarding(file)
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

	cmd := exec.Command("docker", args...)
	output, err := cmd.Output()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}

func boarding(file string) (string, error) {
	if !filepath.IsAbs(file) {
		var err error
		file, err = filepath.Abs(file)
		if err != nil {
			return "", fmt.Errorf(
				`filepath.Abs("%s") returns error: %v`,
				file, err,
			)
		}
	}

	relPath, err := filepath.Rel(codeVolume.dockerPath, file)
	if err != nil {
		return "", fmt.Errorf(
			`filepath.Rel("%s", "%s") returns error: %v`,
			codeVolume.dockerPath, file, err,
		)
	}

	return filepath.Join(codeVolume.hostPath, relPath), nil
}

func setupTempDir() (string, func(), error) {
	tmpdir, err := ioutil.TempDir(codeVolume.dockerPath, "flw-playground-")
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
