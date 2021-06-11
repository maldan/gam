package main

import (
	_ "embed"
	"encoding/json"

	"github.com/maldan/gam/internal/app/gam"
)

//go:embed package.json
var packageJson string

var XConfig PackageJson

type PackageJson struct {
	Version string `json:"version"`
}

func init() {
	json.Unmarshal([]byte(packageJson), &XConfig)
}

func main() {
	gam.Start(XConfig.Version)
}
