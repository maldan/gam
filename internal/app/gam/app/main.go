package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/maldan/gam/internal/app/gam/core"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_net"
	"github.com/maldan/go-cmhp/cmhp_process"
	"github.com/phayes/freeport"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/mod/semver"
)

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

// Install applications
func Install(input string) {
	url := GetAssetFromGithub(GetUrlFromInput(input))
	appName := GetNameFromUrl(url)

	// Check if already exists
	if cmhp_file.Exists(core.GamAppDir + "/" + appName) {
		core.Exit("Application already installed: " + appName)
	}

	// Download app
	zipPath := Download(url)

	// Unzip to app folder
	cmhp_process.Exec("unzip", zipPath, "-d", core.GamAppDir+"/"+appName)

	// Set app executable
	err := os.Chmod(core.GamAppDir+"/"+appName+"/app", 0777)
	if err != nil {
		core.Exit(err.Error())
	}
}

// Download app
func Download(url string) string {
	// Get app name
	appName := GetNameFromUrl(url)

	// Open file
	f, err := os.OpenFile(os.TempDir()+"/"+appName+".zip", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		core.Exit(err.Error())
	}
	defer f.Close()

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		core.Exit(err.Error())
	}

	// Do request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		core.Exit(err.Error())
	}
	defer resp.Body.Close()

	// Draw progress bar
	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)

	return os.TempDir() + "/" + appName + ".zip"
}

// Show application list
func List() {
	// App list
	appList := make([]core.Application, 0)

	// Files
	files, _ := cmhp_file.List(core.GamAppDir)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		// Add application to list
		appList = append(appList, core.Application{
			Name: file.Name(),
			Path: fmt.Sprintf("%v/%v", core.GamAppDir, file.Name()),
		})
	}

	// Print list
	for _, file := range appList {
		version := strings.Replace(file.Name, RemoveVersionFromName(file.Name), "", 1)
		version = strings.Replace(version, "-v", "", 1)

		author := strings.Split(file.Name, "-")[0]
		name := strings.Replace(RemoveVersionFromName(file.Name), author+"-gam-app-", "", 1)

		fmt.Println("name: " + name)
		fmt.Println("author: " + author)
		fmt.Println("version: " + version)
		fmt.Println("path: " + file.Path)
		fmt.Println()
	}
}

func SearchApp(input string) []core.Application {
	list := make([]core.Application, 0)

	// Files
	files, _ := cmhp_file.List(core.GamAppDir)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		// Skip
		if !strings.HasPrefix(file.Name(), input) {
			continue
		}
		list = append(list, core.Application{
			Name:    file.Name(),
			Path:    core.GamAppDir + "/" + file.Name(),
			Version: GetVersionInName(file.Name()),
		})
	}

	sort.SliceStable(list, func(i, j int) bool {
		return semver.Compare(list[i].Version, list[j].Version) == 0
	})

	return list
}

// Get asset
func GetAssetFromGithub(url string) string {
	// Get version
	tuple := strings.Split(url, "@")
	version := ""
	if len(tuple) == 2 {
		version = tuple[1]
	}

	// Request
	response := cmhp_net.Request(cmhp_net.HttpArgs{
		Url: tuple[0],
	})

	// Check error
	if response.StatusCode != 200 {
		core.Exit("Url " + url + " not found")
	}

	// Parse release list
	var releaseList []Release
	json.Unmarshal(response.Body, &releaseList)

	// No releases
	if len(releaseList) == 0 {
		core.Exit("There is no releases")
	}

	// Set default version
	if version == "" {
		version = releaseList[0].TagName
	}

	// Get releases
	for _, release := range releaseList {
		// Check version
		if release.TagName == version {
			// Find asset for current platform
			for _, asset := range release.Assets {
				if asset.Name == "application-"+core.CurrentPlatform+".zip" {
					return asset.DownloadUrl
				}
			}
		}
	}

	// Release not found
	core.Exit("No release for current platform: " + core.CurrentPlatform)
	return ""
}

// Create url from user input
func GetUrlFromInput(input string) string {
	// Parse name
	tuple := strings.Split(input, "/")

	// Default author is maldan
	if len(tuple) == 1 {
		v := strings.Split(tuple[0], "@")
		if len(v) == 2 {
			return "https://api.github.com/repos/maldan/gam-app-" + v[0] + "/releases@v" + v[1]
		}
		return "https://api.github.com/repos/maldan/gam-app-" + v[0] + "/releases"
	}

	// With author
	v := strings.Split(tuple[1], "@")
	if len(v) == 2 {
		return "https://api.github.com/repos/" + tuple[0] + "/gam-app-" + v[0] + "/releases@v" + v[1]
	}
	return "https://api.github.com/repos/" + tuple[0] + "/gam-app-" + tuple[1] + "/releases"
}

// Get name from url
func GetNameFromUrl(url string) string {
	first := strings.Replace(url, "https://github.com/", "", 1)
	first = strings.Replace(first, "releases/download", "", 1)

	u := strings.Split(first, "/")
	return strings.Join(u[0:2], "-") + "-" + u[3]
}

// Create app name from user input
func GetNameFromInput(input string) string {
	input = strings.Replace(input, "@", "-v", 1)
	tuple := strings.Split(input, "/")
	if len(tuple) == 1 {
		return "maldan-gam-app-" + tuple[0]
	}
	return tuple[0] + "-gam-app-" + tuple[1]
}

func RemoveVersionFromName(name string) string {
	var re = regexp.MustCompile(`-v\d+\.\d+\.\d+`)
	return re.ReplaceAllString(name, `$1`)
}

func HasVersionInName(name string) bool {
	var re = regexp.MustCompile(`-v\d+\.\d+\.\d+`)
	return re.Match([]byte(name))
}

func GetVersionInName(name string) string {
	var re = regexp.MustCompile(`\d+\.\d+\.\d+`)
	return re.FindAllString(name, 1)[0]
}

// Remove app
func Delete(input string) {
	// Get app name
	appName := GetNameFromInput(input)

	// Check if already exists
	if !cmhp_file.Exists(core.GamAppDir + "/" + appName) {
		core.Exit("Application not found: " + appName)
	}

	// Delete app dir
	cmhp_file.DeleteDir(core.GamAppDir + "/" + appName)
	fmt.Println("Application deleted: " + appName)
}

// Run app
func Run(input string, args []string) {
	// Get app name
	appName := GetNameFromInput(input)

	// Version not specified
	if !HasVersionInName(appName) {
		list := SearchApp(appName)
		if len(list) == 0 {
			core.Exit("Application not found")
		}
		appName += "-v" + list[0].Version
	}

	// Check if already exists
	if !cmhp_file.Exists(core.GamAppDir + "/" + appName) {
		Install(input)
	}

	// Check port
	port := 0
	portFound := false
	for _, v := range args {
		if strings.Contains(v, "--port=") {
			portFound = true
			x := strings.ReplaceAll(v, "--port=", "")
			xx, _ := strconv.ParseInt(x, 10, 64)
			port = int(xx)
			break
		}
	}

	// Check houst
	hostFound := false
	for _, v := range args {
		if strings.Contains(v, "--host=") {
			hostFound = true
		}
	}

	// Set port
	if !portFound {
		// Port
		p, err := freeport.GetFreePort()
		if err != nil {
			core.Exit(err.Error())
		}
		args = append(args, fmt.Sprintf("--port=%d", p))
		port = p
	}

	// Set host
	if !hostFound {
		args = append(args, "--host="+core.GamConfig.DefaultHost)
	}

	// Add data dir
	args = append(args, fmt.Sprintf("--dataDir=%v/%v", core.GamDataDir, RemoveVersionFromName(appName)))
	args = append(args, fmt.Sprintf("--appId=%v", RemoveVersionFromName(appName)))
	argsFinal := append([]string{core.GamAppDir + "/" + appName + "/app"}, args...)

	// Set logs
	logs, _ := os.OpenFile("/var/log/"+appName+"_info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	errors, _ := os.OpenFile("/var/log/"+appName+"_error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)

	// Run process
	process, err := os.StartProcess(core.GamAppDir+"/"+appName+"/app", argsFinal, &os.ProcAttr{
		Dir: core.GamAppDir + "/" + appName,
		Env: os.Environ(),
		Files: []*os.File{
			nil,
			logs,
			errors,
		},
		Sys: &syscall.SysProcAttr{},
	})
	pid := process.Pid

	if err == nil {
		err = process.Release()
		if err != nil {
			core.Exit(err.Error())
		} else {
			fmt.Println("pid:", pid)
			fmt.Println("port:", port)
			fmt.Println("path:", core.GamAppDir+"/"+appName)
		}
	} else {
		core.Exit(err.Error())
	}
}
