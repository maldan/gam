package core

import (
	"os"
	"runtime"
	"strings"
)

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

// Global path
var GamDir string
var GamAppDir string
var GamDataDir string

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
}

func Exit(msg string) {
	println(msg)
	os.Exit(1)
}
