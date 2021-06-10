package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const CurrentPlatform = runtime.GOOS + "-" + runtime.GOARCH

var GamDir string
var GamAppDir string

type ReleaseList struct {
	ReleaseList []Release `json:"users"`
}

type Release struct {
	Url     string  `json:"url"`
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	DownloadUrl string `json:"browser_download_url"`
}

func ErrorMessage(message string) {
	panic(message)
}

func main() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	GamDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam"
	GamAppDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam-app"
	os.MkdirAll(GamDir, os.ModeDir)
	os.MkdirAll(GamAppDir, os.ModeDir)
	fmt.Println(GamAppDir)

	// Get args
	argsWithoutProg := os.Args[1:]

	// Do action
	switch argsWithoutProg[0] {
	case "install":
		shell_install(argsWithoutProg[1])
	case "service":
		switch argsWithoutProg[1] {
		case "daemon":
			server_start("127.0.0.1:14393")
		default:
		}
	case "process":
		switch argsWithoutProg[1] {
		case "list":
			fmt.Println("SAS")
		default:
		}
	case "run":
		fmt.Println("FUCK")
	default:
		fmt.Println("Gas")
	}

	/*dateCmd := exec.Command("date")
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> date")
	fmt.Println(string(dateOut))*/
}
