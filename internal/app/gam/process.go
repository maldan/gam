package gam

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Process struct {
	Pid  int64
	Name string
	Cmd  string
}

func processList() []Process {
	cmd := exec.Command("wmic", "process", "get", "processid,caption,commandline", "/format:CSV")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()

	xx := fmt.Sprintf("%s", out)
	lines := strings.Split(xx, "\n")

	processList := make([]Process, 1)

	for _, v := range lines {
		line := strings.Split(v, ",")
		if len(line) < 4 {
			continue
		}

		pid, err := strconv.ParseInt(strings.ReplaceAll(line[3], "\r", ""), 10, 0)
		if err != nil {
			continue
		}

		processList = append(processList, Process{
			Pid:  pid,
			Name: line[1],
			Cmd:  line[2],
		})
	}
	return processList
}

func killProcess(pid int) {
	cmd := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Start()
	if err == nil {
		fmt.Println("Daemon stopped") //
	}
}

func killDaemon() {
	pl := processList()
	for _, p := range pl {
		if p.Name == "gam.exe" && strings.Contains(p.Cmd, "service daemon") {
			killProcess(int(p.Pid))
		}
	}
}
