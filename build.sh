#!/bin/bash

GOARCH=amd64 GOOS=linux go build  -ldflags "-s -w" ./internal/app/gam
GOARCH=amd64 GOOS=windows go build  -ldflags "-s -w" ./internal/app/gam

zip -9 -r application-linux-amd64.zip gam
zip -9 -r application-windows-amd64.zip gam.exe