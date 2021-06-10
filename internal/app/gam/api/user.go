package api

import "fmt"

type UserApi int

type SasageoArgs struct {
	Names []interface{}
}

func (u UserApi) PostFuck(args SasageoArgs) []string {
	fmt.Println(args.Names[0])

	return []string{"x", "b"}
}

func (u UserApi) GetList(args SasageoArgs) []string {
	fmt.Println(args.Names)

	return []string{"x", "b"}
}
