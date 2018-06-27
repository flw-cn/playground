package playground

import (
	"io/ioutil"
	"os/exec"
)

func PlayGoString(code string) (string, error) {
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

	return PlayGoFile("main.go")
}

func PlayGoCode(file string) (string, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	return PlayGoString(string(content))
}

func PlayGoFile(file string) (string, error) {
	cmd := exec.Command("go", "run", file)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", &PlayError{
			err:    err,
			output: string(output),
		}
	}

	return string(output), nil
}
