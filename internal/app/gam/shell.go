package gam

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/mod/semver"
)

func shell_run(url string) {
	// Convert name
	appName := convertAppName(url)
	dateCmd := exec.Command(GamAppDir+"/"+appName+"/app.exe", "--port=58123")
	dateCmd.Dir = GamAppDir + "/" + appName
	err := dateCmd.Start()
	if err != nil {
		panic(err)
	}
}

func shell_unzip(source, dest string) error {
	read, err := zip.OpenReader(source)
	if err != nil {
		fmt.Println("Can't open zip")
		return err
	}

	defer read.Close()
	for _, file := range read.File {
		if file.Mode().IsDir() {
			continue
		}
		open, err := file.Open()
		if err != nil {
			return err
		}
		name := path.Join(dest, file.Name)
		os.MkdirAll(path.Dir(name), 0755)
		create, err := os.Create(name)
		if err != nil {
			return err
		}
		defer create.Close()
		create.ReadFrom(open)
	}
	fmt.Println("Unpacked to " + dest)
	return nil
}

func gam_upgrade(path string) {
	killDaemon()

	d1 := []byte("timeout 2\ngam.exe init")
	err := ioutil.WriteFile(path+"/upgrade.cmd", d1, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command(path + "/upgrade.cmd")
	cmd.Dir = path
	err = cmd.Start()
	if err != nil {
		fmt.Println("Can't update gam")
	}
}

func shell_download(url string, appName string) {
	// Open file
	f, _ := os.OpenFile(GamAppDir+"/"+appName+".zip", os.O_CREATE|os.O_WRONLY, 0644)

	defer f.Close()

	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)

	shell_unzip(GamAppDir+"/"+appName+".zip", GamAppDir+"/"+appName)
}

func shell_upgrade() {
	// Get release list
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/maldan/gam/releases", nil)
	resp, err := client.Do(req)
	if err != nil {
		ErrorMessage("Can't get release list")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ErrorMessage("App not found")
	}

	// Parse json
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ErrorMessage("App not found")
	}
	var releaseList []Release
	json.Unmarshal(body, &releaseList)

	// Search release
	for _, release := range releaseList {
		if semver.Compare(release.TagName, Config.Version) <= 0 {
			continue
		}

		for _, asset := range release.Assets {
			if "application-"+CurrentPlatform+".zip" == asset.Name {
				fmt.Println("Found new version", release.TagName)
				shell_download(asset.DownloadUrl, "maldan-gam-"+release.TagName)
				gam_upgrade(GamAppDir + "/" + "maldan-gam-" + release.TagName)
				return
			}
		}
	}

	fmt.Println("New version not found")
}

func shell_install(url string) {
	// Get release list
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+url+"/releases", nil)
	if Config.GithubAccessToken != "" {
		req.Header.Set("Authorization", "token "+Config.GithubAccessToken)
	}

	resp, err := client.Do(req)
	if err != nil {
		ErrorMessage("Can't get release list")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ErrorMessage("App not found")
	}

	// Convert name
	appName := convertAppName(url)

	// Parse json
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ErrorMessage("App not found")
	}
	var releaseList []Release
	json.Unmarshal(body, &releaseList)

	// Search release
	for _, release := range releaseList {
		for _, asset := range release.Assets {
			if "application-"+CurrentPlatform+".zip" == asset.Name {
				fmt.Println("Found")
				appName += "-" + release.TagName
				fmt.Println(appName)
				shell_download(asset.DownloadUrl, appName)
				return
			}
		}
	}

	fmt.Println("Asset not found")
}
