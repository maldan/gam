package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/schollz/progressbar/v3"
)

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
		os.MkdirAll(path.Dir(name), os.ModeDir)
		create, err := os.Create(name)
		if err != nil {
			return err
		}
		defer create.Close()
		create.ReadFrom(open)
	}
	return nil
}

func shell_download(url string, appName string) {
	// Open file
	f, _ := os.OpenFile(GamAppDir+"/"+appName+".zip", os.O_CREATE|os.O_WRONLY, 0644)

	// Check exists
	/*if !os.IsNotExist(err) {
		shell_unzip(GamAppDir+"/"+appName+".zip", GamAppDir+"/application")
		return
	}*/
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

func shell_install(url string) {
	// Get release list
	resp, err := http.Get("https://api.github.com/repos/" + url + "/releases")
	if err != nil {
		ErrorMessage("Can't get release list")
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ErrorMessage("App not found")
	}

	// Convert name
	appName := application_convertName(url)

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
			if asset.Name == "application.zip" {
				// if "application-"+CurrentPlatform+".zip" == asset.Name {
				// fmt.Println("Found")
				appName += "-" + release.TagName
				fmt.Println(appName)
				shell_download(asset.DownloadUrl, appName)
				return
			}
		}
	}
	// fmt.Println(string(body))

	// Create a buffer to write our archive to.
	/*buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Register a custom Deflate compressor.
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})

	fType := reflect.TypeOf(users)
	fmt.Println(fType)*/
}
