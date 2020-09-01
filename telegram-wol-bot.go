package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"path/filepath"
)

func fatalError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	var opts struct {
		ConfigPath string `short:"c" long:"config" description:"Path of the config file" required:"true"`
	}

	_, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}

	actualPath, err := filepath.Abs(opts.ConfigPath)
	fatalError(err)

	fmt.Printf("配置文件位于 %s\n", actualPath)
}
