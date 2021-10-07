package core

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/maldan/go-cmhp/cmhp_file"
)

type Process struct {
	Pid  int64             `json:"pid"`
	Cmd  string            `json:"cmd"`
	Args map[string]string `json:"args"`
}

type Application struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
}

type RepoApplication struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Rating      int       `json:"rating"`
	LastUpdate  time.Time `json:"lastUpdate"`
}

type Config struct {
	DefaultHost string

	SPACES_KEY      string `json:"SPACES_KEY"`
	SPACES_SECRET   string `json:"SPACES_SECRET"`
	SPACES_ENDPOINT string `json:"SPACES_ENDPOINT"`
	SPACES_BUCKET   string `json:"SPACES_BUCKET"`
}

// Global path
var GamDir string
var GamAppDir string
var GamDataDir string
var GamConfig Config

// Global const
const CurrentPlatform = runtime.GOOS + "-" + runtime.GOARCH

func init() {
	// Get home dir and set app dir
	dirname, err := os.UserHomeDir()
	if err != nil {
		Exit(err.Error())
	}
	GamAppDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam-app"
	GamDataDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam-data"
	GamDir = strings.ReplaceAll(dirname, "\\", "/") + "/.gam"

	// Create folders
	os.MkdirAll(GamDir, 0755)
	os.MkdirAll(GamAppDir, 0755)
	os.MkdirAll(GamDataDir, 0755)

	// Load config
	cmhp_file.ReadJSON(GamDataDir+"/gam/config.json", &GamConfig)
	if GamConfig.DefaultHost == "" {
		GamConfig.DefaultHost = "localhost"
	}
}

func Exit(msg string) {
	println(msg)
	os.Exit(1)
}
