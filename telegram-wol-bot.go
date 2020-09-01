package main

import (
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	Token string
	ChatID float64
	Computers []Computers
}

type Computers struct {
	Name string
	Mac string
}

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

	configPath, err := filepath.Abs(opts.ConfigPath)
	fatalError(err)

	jsonFile, err := os.Open(configPath)
	if err != nil {
		fmt.Printf("Error opening JSON file: %s\n", err.Error())
		os.Exit(1)
	}
	defer jsonFile.Close()

	content, err := ioutil.ReadAll(jsonFile)
	fatalError(err)

	var config Config
	err = json.Unmarshal(content, &config)
	fatalError(err)

	fmt.Println(config.Token)
	fmt.Println(int64(config.ChatID))
	fmt.Println(config.Computers[0])
}
