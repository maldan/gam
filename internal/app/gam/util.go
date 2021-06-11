package gam

import "strings"

func convertAppName(url string) string {
	return strings.ReplaceAll(url, "/", "-")
}
