package main

import (
	"log"

	"github.com/flw-cn/playground/docker"
)

func main() {
	code := `
	cmd := exec.Command("uname", "-a")
	output, _ := cmd.Output()
	fmt.Printf("%s", output)
`
	output, err := docker.PlayCode(code)
	if err != nil {
		log.Printf("Error: %s\n%s", err, output)
		return
	}

	sep := "-------------------"

	log.Println("output:")
	log.Println(sep)
	log.Println(output)
	log.Println(sep)
}
