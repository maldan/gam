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

	/*bar := progressbar.Default(100)
	for i := 0; i < 1; i++ {
		bar.Add(2)
		time.Sleep(20 * time.Millisecond)
	}*/

	// Get args
	argsWithoutProg := os.Args[1:]

	// Do action
	switch argsWithoutProg[0] {
	case "install":
		shell_install(argsWithoutProg[1])
	case "start":
		fmt.Println("FUCK")
	case "stop":
		fmt.Println("FUCK")
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

	/*myData := map[string]interface{}{
		"Name": "Tony",
		//"Age":  "fuck",
		//"Fuck": true,
		//"Sex":  23,
	}
	var result MyStruct
	transcode(myData, &result)
	fmt.Printf("%+v\n", result)*/

	// sss

	start()
}
