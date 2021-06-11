package main

import "runtime"

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