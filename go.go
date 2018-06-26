package playground

import (
	"io/ioutil"
	"os/exec"
)

func PlayGo(code string) (string, error) {
	_, deferFunc, err := setupTempDir()
	if err != nil {
		return "", err
	}

	defer deferFunc()

	data := []byte{}

	data = append(data, []byte("package main\n")...)
	data = append(data, []byte("func main() {\n")...)
	data = append(data, []byte(code)...)
	data = append(data, []byte("}\n")...)

	err = ioutil.WriteFile("main.go", data, 0666)
	if err != nil {
		return "", err
	}

	cmd := exec.Command("goimports", "-w", "main.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", &PlayError{
			err:    err,
			output: string(output),
		}
	}

	cmd = exec.Command("go", "run", "main.go")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return "", &PlayError{
			err:    err,
			output: string(output),
		}
	}

	return string(output), nil
}

func PlayGoFile(file string) (string, error) {
	code, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return PlayGo(string(code))
}
