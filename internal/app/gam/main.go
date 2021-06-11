package gam

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

// Ssageo
func ErrorMessage(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func loadConfig() {
	f, err := os.Open("config.json")
	if err != nil {
		return
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&Config)
	if err != nil {
		return
	}
}

func init() {
	// Prepare
	loadConfig()
}

func Start(version string) {
	// Set version
	Config.Version = "v" + version

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

	// fmt.Println("Current platform: ", CurrentPlatform)

	// Create folders
	os.MkdirAll(GamDir, 0755)
	os.MkdirAll(GamAppDir, 0755)

	if len(os.Args) <= 1 {
		fmt.Println("No params")
		return
	}

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
		case "kill":
			killDaemon()
		}

	case "process":
		switch argsWithoutProg[1] {
		case "list":
			pl := processList()
			for _, p := range pl {
				if p.Name == "app.exe" && strings.Contains(strings.ReplaceAll(p.Cmd, "\\", "/"), GamAppDir) {
					fmt.Println(p)
				}
			}
		}
	case "upgrade":
		shell_upgrade()
	case "version":
		fmt.Printf("%v", version)
	case "run":
		shell_run(argsWithoutProg[1])
	}
}
