package main

import (
	"github.com/AWildBeard/errata/flags"
	"os"
)

func main() {
	file, err := os.OpenFile("config.yaml",  os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	err = flags.WriteExampleConfig(file)
	if err != nil {
		panic(err)
	}
}

