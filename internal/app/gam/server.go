package gam

import (
	"mime/multipart"

	"github.com/maldan/gam/internal/app/gam/api"
	"github.com/maldan/go-restserver"
)

type SasageoArgs struct {
	Context     *restserver.RestServerContext
	AccessToken string
	Files       map[string][]*multipart.FileHeader

	Name string
	Age  int64
	Fuck bool
}

func server_start(addr string) {
	restserver.Start(addr, map[string]interface{}{
		"user":        new(api.UserApi),
		"application": new(api.ApplicationApi),
	})
}
