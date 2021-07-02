package gam

import (
	"fmt"
	"strings"
	"time"
)

func daemon_start() {
	for {
		for _, v := range Config.KeepAlive {
			appName := strings.Split(v, " ")[0]
			/*pl := gamProcessList()
			fmt.Println(pl)
			for _, p := range pl {
				fmt.Println(p.Cmd, appName)
				if strings.Contains(p.Cmd, appName) {
					fmt.Println("Found")
					break
				}
			}*/
			fmt.Println(appName)
		}
		time.Sleep(time.Second * 2)
	}
}
