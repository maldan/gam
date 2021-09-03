package app

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/maldan/gam/internal/app/gam/core"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_net"
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

// Download apps
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
