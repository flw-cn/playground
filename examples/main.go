package main

import (
	"log"

	"github.com/flw-cn/playground"
)

func main() {
	output, err := playground.PlayGo(`fmt.Printf("Hello, world hhhh!\n")`)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	sep := "-------------------"

	log.Println("output:")
	log.Println(sep)
	log.Println(output)
	log.Println(sep)
}
