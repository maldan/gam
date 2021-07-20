package gam

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func process_list() []Process {
	// Process list
	cmd := exec.Command("ps", "-aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	processList := make([]Process, 0)

	// Lines
	lines := strings.Split(out.String(), "\n")
	for _, v := range lines {
		// Prepare lines
		line := strings.ReplaceAll(v, "  ", " ")
		for i := 0; i < 10; i += 1 {
			line = strings.ReplaceAll(line, "  ", " ")
		}

		// Split tuple
		tuple := strings.Split(line, " ")
		if len(tuple) < 4 {
			continue
		}

		// Parse int
		pid, err := strconv.ParseInt(tuple[1], 10, 0)
		if err != nil {
			continue
		}

		// Cmd
		cmd := strings.Join(tuple[10:], " ")

		// Parse args
		args := make(map[string]string)
		t := strings.Split(cmd, " ")
		for _, v := range t {
			v2 := strings.Split(v, "=")
			if len(v2) < 2 {
				continue
			}
			args[strings.Replace(v2[0], "--", "", 1)] = v2[1]
		}

		processList = append(processList, Process{
			Pid:  pid,
			Name: strings.Split(strings.Replace(cmd, GamAppDir, "", 1), " ")[0],
			Cmd:  cmd,
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
	cmd := exec.Command("kill", strconv.Itoa(pid))
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
