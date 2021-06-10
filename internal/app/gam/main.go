package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

func ErrorMessage(message string) {
	panic(message)
}

func main() {
	// Get home dir and set app dir
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	GamAppDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam-app"

	// Set gam dir
	switch runtime.GOOS {
	case "windows":
		GamDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam"
	case "linux":
		GamDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam"
	default:
		panic("Unsupported platform")
	}

	fmt.Println("Current platform: ", CurrentPlatform)

	// Create folders
	os.MkdirAll(GamDir, 0755)
	os.MkdirAll(GamAppDir, 0755)

	// Get args
	argsWithoutProg := os.Args[1:]

	// Do action
	switch argsWithoutProg[0] {
	case "init":
		if runtime.GOOS == "windows" {
			source, _ := os.Open(os.Args[0])
			destination, err := os.Create(GamDir + "/gam.exe")
			if err != nil {
				panic(err)
			}
			io.Copy(destination, source)
		}
		if runtime.GOOS == "linux" {
			source, _ := os.Open(os.Args[0])
			destination, err := os.Create(GamDir + "/gam")
			if err != nil {
				panic(err)
			}
			io.Copy(destination, source)
			os.Chmod(GamDir+"/gam", 0755)

			// Create script
			d1 := []byte("#!/bin/bash\n~/.gam/gam \"$@\"")
			ioutil.WriteFile("/usr/local/bin/gam", d1, 0755)
		}
	case "install":
		shell_install(argsWithoutProg[1])
	case "service":
		switch argsWithoutProg[1] {
		case "daemon":
			server_start("127.0.0.1:14393")
		}
	case "process":
		switch argsWithoutProg[1] {
		case "list":
			fmt.Println("SAS")
		}
	case "run":
		fmt.Println("FUCK")
	}
}
