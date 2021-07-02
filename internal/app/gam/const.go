package gam

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const CurrentPlatform = runtime.GOOS + "-" + runtime.GOARCH
const DetachedProcess = 0x00000008

var Format string
var GamDir string
var GamAppDir string
var GamDataDir string
var Config GamConfig

type Release struct {
	Url     string  `json:"url"`
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type ReleaseList struct {
	ReleaseList []Release `json:"users"`
}

type Asset struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	DownloadUrl string `json:"browser_download_url"`
}

type GamConfig struct {
	Version           string
	GithubAccessToken string   `json:"GITHUB_ACCESS_TOKEN"`
	KeepAlive         []string `json:"KEEP_ALIVE"`
}

type Process struct {
	Pid  int64             `json:"pid"`
	Name string            `json:"name"`
	Cmd  string            `json:"cmd"`
	Args map[string]string `json:"args"`
}

type Application struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

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
	// Handle flags
	for _, v := range os.Args {
		if v == "--format=json" {
			Format = "json"
		}
	}

	// Get home dir and set app dir
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	GamAppDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam-app"
	GamDataDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam-data"

	// Set gam dir
	switch runtime.GOOS {
	case "windows":
		GamDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam"
	case "linux":
		GamDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam"
	default:
		panic("Unsupported platform")
	}

	// Create folders
	os.MkdirAll(GamDir, 0755)
	os.MkdirAll(GamAppDir, 0755)
	os.MkdirAll(GamDataDir, 0755)

	// Prepare
	loadConfig()
}
