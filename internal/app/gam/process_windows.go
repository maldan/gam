package gam

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func process_list() []Process {
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

		// Parse args
		args := make(map[string]string)
		t := strings.Split(line[2], " ")
		for _, v := range t {
			v2 := strings.Split(v, "=")
			if len(v2) < 2 {
				continue
			}
			args[strings.Replace(v2[0], "--", "", 1)] = v2[1]
		}

		processList = append(processList, Process{
			Pid:  pid,
			Name: line[1],
			Cmd:  line[2],
			Args: args,
		})
	}
	return processList
}

func process_gamList() []Process {
	out := make([]Process, 0)

	pl := process_list()
	for _, p := range pl {
		// Skip if go run process
		if p.Name == "go.exe" {
			continue
		}
		cmd := strings.ReplaceAll(p.Cmd, "\\", "/")
		if strings.Contains(cmd, "--appId") && strings.Contains(cmd, "--port") {
			out = append(out, p)
		}
	}

	return out
}

func process_killAllGamList() {
	pl := process_gamList()
	for _, p := range pl {
		process_kill(int(p.Pid))
	}
}

func process_kill(pid int) {
	cmd := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid))
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Start()
	if err == nil {
		fmt.Printf("Process %d stopped\n", pid)
	}
}

func process_killDaemon() {
	pl := process_list()
	for _, p := range pl {
		if p.Name == "gam.exe" && strings.Contains(p.Cmd, "service daemon") {
			process_kill(int(p.Pid))
		}
	}
}
