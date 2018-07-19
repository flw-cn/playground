package docker

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var supportedLang map[string]string

var codeVolume struct {
	hostPath   string
	dockerPath string
}

func init() {
	supportedLang = map[string]string{
		"go":     "go",
		"golang": "go",
	}

	codeVolume.hostPath = ""
	codeVolume.dockerPath = ""
}

func Boarding(hostPath, dockerPath string) error {
	if hostPath == "" {
		return errors.New("must provide playground path")
	} else if !filepath.IsAbs(hostPath) {
		return fmt.Errorf("%s is not an absolute path", hostPath)
	} else {
		codeVolume.hostPath = hostPath
	}

	if dockerPath == "" {
		codeVolume.dockerPath = hostPath
	} else if !filepath.IsAbs(dockerPath) {
		return fmt.Errorf("%s is not an absolute path", dockerPath)
	} else {
		codeVolume.dockerPath = dockerPath
	}

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
	if lang = supportedLang[lang]; lang == "" {
		return "", fmt.Errorf("Unsupported language: %s", lang)
	}

	tmpdir, deferFunc, err := setupTempDir()
	if err != nil {
		return "", err
	}

	defer deferFunc()

	baseName := filepath.Base(file)
	newFile := filepath.Join(tmpdir, baseName)
	err = os.Rename(file, newFile)
	if err != nil {
		return "", err
	}

	return play(lang, newFile, false)
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
		"create", "--rm",
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker create failed: %s", err)
	}

	container := string(output)
	fmt.Printf("container: %s", container)
	if len(container) < 12 {
		return "", errors.New("docker create failed")
	}

	container = container[0:12]
	cmd = exec.CommandContext(ctx, "docker", "start", "--attach", container)
	output, err = cmd.CombinedOutput()
	if err != nil {
		exec.Command("docker", "kill", container).Run()
		if ctx.Err() != nil {
			err = ctx.Err()
		}
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
