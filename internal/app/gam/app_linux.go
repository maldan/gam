package gam

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/phayes/freeport"
)

func app_run(url string, args []string) {
	// Port
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}

	// Convert name
	appName := app_name(url)
	folder := app_folder(appName)

	// Folder
	if folder == "" {
		ErrorMessage(fmt.Sprintf("App %v not found", url))
	}

	// Port
	portFound := false
	for _, v := range args {
		if strings.Contains(v, "--port=") {
			portFound = true
			break
		}
	}

	// Port
	if !portFound {
		args = append(args, fmt.Sprintf("--port=%d", port))
	}

	// Add data dir
	args = append(args, fmt.Sprintf("--dataDir=%v/%v", GamDataDir, app_name_without_version(appName)))
	args = append(args, fmt.Sprintf("--appId=%v", app_name_without_version(appName)))
	argsFinal := append([]string{GamAppDir + "/" + folder + "/app"}, args...)

	sysproc := &syscall.SysProcAttr{
		Noctty: true,
	}

	// Set process params
	logs, _ := os.OpenFile("/var/log/"+appName+"_info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	errors, _ := os.OpenFile("/var/log/"+appName+"_error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)

	attr := os.ProcAttr{
		Dir: GamAppDir + "/" + folder,
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			logs,
			errors,
		},
		Sys: sysproc,
	}

	// Run process
	process, err := os.StartProcess(GamAppDir+"/"+folder+"/app", argsFinal, &attr)
	if err == nil {
		err = process.Release()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("pid:%v, port:%v\n", process.Pid, port)
		}
	} else {
		fmt.Println(err.Error())
	}
}
