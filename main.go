package main

import (
	"credondocr-dd-technical-test/config"
	"credondocr-dd-technical-test/server"
	"flag"
	"fmt"
	"os"
)

func main() {
	environment := flag.String("e", "development", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()
	config.Init(*environment)
	server.Init()
}
