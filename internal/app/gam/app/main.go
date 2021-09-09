package app

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/maldan/gam/internal/app/gam/core"
	"github.com/maldan/go-cmhp/cmhp_crypto"
	"github.com/maldan/go-cmhp/cmhp_file"
	"github.com/maldan/go-cmhp/cmhp_net"
	"github.com/maldan/go-cmhp/cmhp_process"
	"github.com/maldan/go-cmhp/cmhp_s3"
	"github.com/maldan/go-cmhp/cmhp_time"
	"github.com/phayes/freeport"
)

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

// Backup data
func Backup(input string) {
	// Get app name
	appName := RemoveVersionFromName(GetNameFromInput(input))

	// Check if data exists
	if !cmhp_file.Exists(core.GamDataDir + "/" + appName) {
		core.Exit("Data not found")
	}

	// Init config
	cmhp_s3.Start(core.GamDataDir + "/gam/config.json")

	// Zip app data
	zipFile := os.TempDir() + "/" + cmhp_crypto.UID(10) + ".zip"
	p, _ := cmhp_process.Create("zip", "-9", "-r", zipFile, ".", "-i", "*")
	p.Dir = core.GamDataDir + "/" + appName
	p.Run()
	defer cmhp_file.Delete(zipFile)

	// Get zip file
	zipData, err := cmhp_file.ReadBin(zipFile)
	if err != nil {
		core.Exit(err.Error())
	}

	// Upload
	zipHash, _ := cmhp_file.HashMd5(zipFile)
	url, err := cmhp_s3.Write(cmhp_s3.WriteArgs{
		Path:        "backup/gam-data/" + appName + "/" + cmhp_time.Format(time.Now(), "YYYY-MM-DD") + "_" + zipHash[0:8] + ".zip",
		InputData:   zipData,
		Visibility:  "public-read",
		ContentType: "application/zip",
	})
	if err != nil {
		core.Exit(err.Error())
	}
	fmt.Println("Uploaded to", url)
}

// Print backup list
func BackupList(input string) {
	// Init config
	cmhp_s3.Start(core.GamDataDir + "/gam/config.json")

	// Get app name
	appName := RemoveVersionFromName(GetNameFromInput(input))
	files := cmhp_s3.List("backup/gam-data/" + appName)
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].LastModified.Unix() < files[j].LastModified.Unix()
	})
	for _, file := range files {
		fileName := strings.Replace(file.Path, "backup/gam-data/"+appName+"/", "", 1)
		if fileName == "" {
			continue
		}
		fmt.Println("file:", fileName)
		fmt.Println("size:", humanize.Bytes(uint64(file.Size)))
		fmt.Println("lastModified:", cmhp_time.Format(file.LastModified, "YYYY-MM-DD HH:mm:ss"))
		fmt.Println()
	}
}

func ShowRepoList(input string) {
	// Load repo list
	list := make([]string, 0)
	cmhp_file.ReadJSON("repo.json", &list)

	// Load cache
	cache := make([]core.RepoApplication, 0)
	cmhp_file.ReadJSON(core.GamDataDir+"/gam/repo.cache.json", &cache)

	// https://raw.githubusercontent.com/maldan/gam-app-fileman/main/icon.svg

	// Go over list
	for _, v := range list {
		tuple := strings.Split(strings.Replace(v, "https://github.com/", "", 1), "/")
		author := tuple[0]
		appName := strings.Replace(tuple[1], "gam-app-", "", 1)

		// Search in cache
		isFound := false
		for _, vv := range cache {
			if vv.Name == appName && vv.Author == author {
				isFound = true
				break
			}
		}

		// Not found
		if !isFound {
			// Get info about repo
			jj := make(map[string]interface{})
			cmhp_net.Request(cmhp_net.HttpArgs{
				Url:        "https://api.github.com/repos/" + author + "/gam-app-" + appName,
				Method:     "GET",
				OutputJSON: &jj,
			})

			// Create app
			tt := strings.Replace(jj["updated_at"].(string), "Z", "", 1)
			tt = strings.Replace(tt, "T", " ", 1)
			updatedAt, _ := time.Parse("2006-01-02 15:04:05", tt)
			repoApp := core.RepoApplication{
				Name:        appName,
				Description: jj["description"].(string),
				Rating:      int(jj["stargazers_count"].(float64)),
				Author:      author,
				LastUpdate:  updatedAt,
			}
			cache = append(cache, repoApp)
		}
	}

	// Save file
	cmhp_file.WriteJSON(core.GamDataDir+"/gam/repo.cache.json", &cache)
	for _, v := range cache {
		fmt.Println("name:", v.Name)
		fmt.Println("description:", v.Description)
		fmt.Println("author:", v.Author)
		fmt.Println("rating:", v.Rating)
		fmt.Println("update:", cmhp_time.Format(v.LastUpdate, "YYYY-MM-DD HH:mm:ss"))
		fmt.Println()
	}
}

// Exec command
func Execute(input string, args []string) {
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

	// Run process
	argsFinal := []string{core.GamAppDir + "/" + appName + "/app", "--cmd=" + strings.Join(args, " ")}
	process, err := os.StartProcess(core.GamAppDir+"/"+appName+"/app", argsFinal, &os.ProcAttr{
		Dir: core.GamAppDir + "/" + appName,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
		Sys: &syscall.SysProcAttr{},
	})
	if err != nil {
		core.Exit(err.Error())
	}
	process.Release()
}
