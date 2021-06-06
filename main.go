package main

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/schollz/progressbar/v3"
)

type ReleaseList struct {
	ReleaseList []Release `json:"users"`
}

type Release struct {
	Url string `json:"url"`
}

func installApplication(x string) {
	resp, err := http.Get("https://api.github.com/repos/maldan/vsu-password-manager/releases")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	//fmt.Println(body)
	// bodyString := string(body)
	// fmt.Println(bodyString) s

	var users []Release
	json.Unmarshal(body, &users)

	fmt.Println(users[0].Url)

	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Register a custom Deflate compressor.
	w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestCompression)
	})
}

func main() {

	bar := progressbar.Default(100)
	for i := 0; i < 25; i++ {
		bar.Add(4)
		time.Sleep(10 * time.Millisecond)
	}

	// Get args
	argsWithoutProg := os.Args[1:]

	// Do action
	switch argsWithoutProg[0] {
	case "install":
		installApplication(argsWithoutProg[1])
	case "start":
		fmt.Println("FUCK")
	case "stop":
		fmt.Println("FUCK")
	case "run":
		fmt.Println("FUCK")
	default:
		fmt.Println("Gas")
	}

	dateCmd := exec.Command("date")
	dateOut, err := dateCmd.Output()
	if err != nil {
		panic(err)
	}
	fmt.Println("> date")
	fmt.Println(string(dateOut))

	// sss

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
