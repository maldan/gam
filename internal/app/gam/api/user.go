package api

import "fmt"

func Fuck(p map[string]interface{}) {
	fmt.Println("S")
}

func Sasageo() map[string]func(map[string]interface{}) {
	return map[string]func(map[string]interface{}){
		"s": Fuck,
	}
}
