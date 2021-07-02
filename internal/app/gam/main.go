package gam

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	"github.com/k0kubun/pp"
)

func Start(version string) {
	// Set version
	Config.Version = "v" + version

	if len(os.Args) <= 1 {
		fmt.Println("No params")
		return
	}

	// Get args
	argsWithoutProg := os.Args[1:]

	// Do action
	switch argsWithoutProg[0] {
	case "init":
		if runtime.GOOS == "windows" {
			source, _ := os.Open(os.Args[0])
			destination, err := os.Create(GamDir + "/gam.exe")
			if err != nil {
				panic(err)
			}
			io.Copy(destination, source)
		}
		if runtime.GOOS == "linux" {
			source, _ := os.Open(os.Args[0])
			destination, err := os.Create(GamDir + "/gam")
			if err != nil {
				panic(err)
			}
			io.Copy(destination, source)
			os.Chmod(GamDir+"/gam", 0755)

			// Create script
			d1 := []byte("#!/bin/bash\n~/.gam/gam \"$@\"")
			ioutil.WriteFile("/usr/local/bin/gam", d1, 0755)
		}
	case "install":
		app_install(argsWithoutProg[1])
	case "update":
		app_install(argsWithoutProg[1])
		app_clean(argsWithoutProg[1])
	case "service":
		switch argsWithoutProg[1] {
		case "daemon":
			daemon_start()
		case "kill":
			process_killDaemon()
		}
	case "app":
		switch argsWithoutProg[1] {
		case "list":
			app_list()
		case "clean":
			app_clean(argsWithoutProg[2])
		}
	case "process":
		switch argsWithoutProg[1] {
		case "list":
			pl := process_gamList()

			if Format == "json" {
				r, _ := json.Marshal(pl)
				fmt.Println(string(r))
			} else {
				for _, p := range pl {
					pp.Println(p)
				}
			}

		case "kill":
			if argsWithoutProg[2] == "all" {
				process_killAllGamList()
			} else {
				pid, _ := strconv.Atoi(argsWithoutProg[2])
				process_kill(pid)
			}
		}
	case "upgrade":
		app_upgrade()
	case "version":
		fmt.Printf("%v", version)
	case "run":
		app_run(argsWithoutProg[1], argsWithoutProg[2:])
	case "help":
		commands := map[string]string{
			"init":              "Install and init gam",
			"install $url":      "Download and install app frop gihtub",
			"run $url":          "Run app",
			"process list":      "Show running gam applications",
			"process kill $pid": "Kill process by pid",
			"process kill all":  "Kill all running gam applications",
			"version":           "Show current gam version",
			"upgrade":           "Install and download new gam version",
		}

		for k, v := range commands {
			color.Green(k)
			fmt.Println(" - " + v)
		}
	default:
		fmt.Println("Unknown command")
	}
}
