package process

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/maldan/gam/internal/app/gam/core"
	"github.com/maldan/go-cmhp/cmhp_process"
)

func List() []core.Process {
	// Process list
	cmd := exec.Command("ps", "-aux")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	processList := make([]core.Process, 0)

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

		// Process list
		name := strings.Split(strings.Replace(cmd, core.GamAppDir+"/", "", 1), " ")[0]
		name = strings.Replace(name, "/app", "", 1)
		processList = append(processList, core.Process{
			Pid:  pid,
			Name: name,
			Cmd:  cmd,
			Args: args,
		})
	}

	return processList
}

func GamList() []core.Process {
	out := make([]core.Process, 0)

	pl := List()
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

// Kill process
func Kill(input string) {
	if input == "all" {
		pl := GamList()
		for _, p := range pl {
			cmhp_process.Exec("kill", fmt.Sprintf("%v", p.Pid))
		}
		return
	}
	b := cmhp_process.Exec("kill", input)
	core.Exit(b)
}
