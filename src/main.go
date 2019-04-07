package main

import (
	"./token"
	"fmt"
)

func main() {
	tk, err := token.NewToken("client_Id", "user_name", []string{"auth", "login"}, "idfa")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("token str:", tk.String())

	newtoken, err := token.ParseToken(tk.String())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("new token str:", newtoken.String())

	if newtoken.CheckScopes([]string{"auth"}) {
		fmt.Println("token has permission")
	}

	if !newtoken.CheckScopes([]string{"pay"}) {
		fmt.Println("token has no permission")
	}
}
