package main

import (
	_ "embed"
	"encoding/json"

	"github.com/maldan/gam/internal/app/gam"
)

//go:embed package.json
var packageJson string

// Package json config
var Config struct {
	Version string `json:"version"`
}

func main() {
	json.Unmarshal([]byte(packageJson), &Config)
	gam.Start(Config.Version)
}
