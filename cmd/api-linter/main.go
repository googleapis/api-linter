package main

import (
	"log"
	"os"
)

func main() {
	if err := runCLI(rules(), configs(), os.Args); err != nil {
		log.Fatal(err)
	}
}
