package main

import (
	"fmt"
	"mime/multipart"

	"github.com/maldan/gam/internal/app/gam/api"
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

func post_run(args SasageoArgs) string {
	fmt.Println("files", args.Files)
	return args.Name
}

/*func get_sasageo(args SasageoArgs) string {
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
}*/

func server_start(addr string) {
	controller := map[string]interface{}{
		"user":        new(api.UserApi),
		"application": new(api.ApplicationApi),
	}
	rest_server.Start(addr, controller)
}
