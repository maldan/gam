package api

import "fmt"

type UserApi int

type SasageoArgs struct {
	Names []interface{}
}

// Fuck user
func (u UserApi) PostFuck(args SasageoArgs) []string {
	fmt.Println(args.Names[0])

	return []string{"x", "b"}
}

// Get user list
func (u UserApi) GetList(args SasageoArgs) []string {
	fmt.Println(args.Names)

	return []string{"x", "b"}
}
