package docker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
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

	codeVolume.hostPath = "/tmp"
	codeVolume.dockerPath = "/tmp"
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
	err = CopyFile(file, newFile)
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

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

	var cmd []string
	var mountPoint string
	if isCodePice {
		mountPoint = "/code/piece"
		cmd = []string{"--lang", lang, "--code", mountPoint}
	} else {
		mountPoint = "/code/main.go"
		cmd = []string{"--lang", lang, "--file", mountPoint}
	}

	resp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: "flwos/playground",
			Tty:   true,
			Cmd:   cmd,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: path,
					Target: mountPoint,
				},
			},
		},
		nil, "",
	)
	if err != nil {
		return "", err
	}

	containerID := resp.ID

	defer cli.ContainerRemove(
		context.Background(),
		containerID,
		types.ContainerRemoveOptions{Force: true},
	)

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	io.Copy(buf, out)

	return buf.String(), nil
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
