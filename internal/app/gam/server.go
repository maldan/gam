package main

import (
	"fmt"
	"mime/multipart"

	"github.com/maldan/gam/internal/app/rest_server"
)

type SasageoArgs struct {
	Context     *rest_server.RestServerContext
	AccessToken string
	Files       map[string][]*multipart.FileHeader

	Name string
	Age  int64
	Fuck bool
}

func post_sasageo(args SasageoArgs) string {
	fmt.Println("files", args.Files)
	return args.Name
}

func get_sasageo(args SasageoArgs) string {
	return "SAs"
}

func get_gagaseo(args SasageoArgs) int {
	return 32
}

func get_konodioda(args SasageoArgs) map[string]string {
	return map[string]string{
		"a": "b",
		"c": "x",
	}
}

func start() {
	controller := map[string]interface{}{
		"post_sasageo":  post_sasageo,
		"get_sasageo":   get_sasageo,
		"get_gagaseo":   get_gagaseo,
		"get_konodioda": get_konodioda,
	}
	rest_server.Start("127.0.0.1:8080", controller)
}
