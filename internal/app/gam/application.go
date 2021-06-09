package main

import "strings"

func application_convertName(url string) string {
	return strings.ReplaceAll(url, "/", "-")
}
