package gam

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/maldan/go-cmhp"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/mod/semver"
)

func app_folder(prefix string) string {
	files, err := ioutil.ReadDir(GamAppDir)
	if err != nil {
		log.Fatal(err)
	}

	finalFiles := make([]fs.FileInfo, 0)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		if !strings.HasPrefix(file.Name(), prefix) {
			continue
		}

		finalFiles = append(finalFiles, file)
	}

	sort.Slice(finalFiles, func(i, j int) bool {
		a := strings.Replace(finalFiles[i].Name(), prefix+"-", "", 1)
		b := strings.Replace(finalFiles[j].Name(), prefix+"-", "", 1)
		return semver.Compare(a, b) > 0
	})

	if len(finalFiles) <= 0 {
		return ""
	}

	return finalFiles[0].Name()
}

func app_name(url string) string {
	return strings.ReplaceAll(url, "/", "-")
}

func app_name_without_version(name string) string {
	var re = regexp.MustCompile(`-v\d+\.\d+\.\d+`)
	return re.ReplaceAllString(name, `$1`)
}

func app_list() {
	files, err := ioutil.ReadDir(GamAppDir)
	if err != nil {
		log.Fatal(err)
	}

	appList := make([]Application, 0)

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		appList = append(appList, Application{
			Name: file.Name(),
			Path: fmt.Sprintf("%v/%v", GamAppDir, file.Name()),
		})
	}

	if Format == "json" {
		s, _ := json.Marshal(appList)
		fmt.Println(string(s))
	} else {
		for _, file := range appList {
			fmt.Println(file.Name)
		}
	}
}

func app_unzip(source, dest string) error {
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
	err = os.Chmod(dest+"/app", 0777)
	if err != nil {
		return err
	}
	return nil
}

func gam_upgrade(path string) {
	process_killDaemon()

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

func app_download(url string, appName string) {
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

	err := app_unzip(GamAppDir+"/"+appName+".zip", GamAppDir+"/"+appName)
	if err != nil {
		fmt.Println(err)
	}
}

func app_upgrade() {
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
				app_download(asset.DownloadUrl, "maldan-gam-"+release.TagName)
				gam_upgrade(GamAppDir + "/" + "maldan-gam-" + release.TagName)
				return
			}
		}
	}

	fmt.Println("New version not found")
}

func app_clean(url string) {
	// Convert name
	appName := app_name(url)
	folder := app_folder(appName)

	// Folder
	if folder == "" {
		ErrorMessage(fmt.Sprintf("App %v not found", url))
	}

	files, err := ioutil.ReadDir(GamAppDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		if !strings.Contains(file.Name(), appName) {
			continue
		}
		if file.Name() == folder {
			continue
		}
		fmt.Printf("Removing %v...\n", GamAppDir+"/"+file.Name())
		cmhp.DirDelete(GamAppDir + "/" + file.Name())
	}
}

func app_remove(url string) {
	// Convert name
	appName := app_name(url)
	folder := app_folder(appName)

	// Folder
	if folder == "" {
		ErrorMessage(fmt.Sprintf("App %v not found", url))
	}

	fmt.Printf("Removing %v...\n", folder)
	cmhp.DirDelete(folder)
}

func app_install(url string) {
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
	appName := app_name(url)

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
				app_download(asset.DownloadUrl, appName)
				return
			}
		}
	}

	fmt.Println("Asset not found")
}
