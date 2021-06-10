package api

import (
	"os/exec"
)

type ApplicationApi int

type ApplicationRunArgs struct {
	Url string
}

func (a ApplicationApi) PostRun(args ApplicationRunArgs) string {
	dateCmd := exec.Command("gam", "run", "maldan/vsu-password-manager")
	_, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	// fmt.Println("> date")
	// fmt.Println(string(dateOut))

	return "ok"
}
