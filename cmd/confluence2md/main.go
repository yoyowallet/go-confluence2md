package main

import (
	"log"
	"os"

	confluence2md "github.com/yoyowallet/go-confluence2md"
)

func main() {
	err := confluence2md.Convert(os.Stdin, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
