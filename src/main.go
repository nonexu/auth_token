package main

import (
	"./token"
	"fmt"
	"time"
)

func main() {
	tk, err := token.NewToken("GAME_ID", "auth")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(1, tk.String())

	newtoken, err := token.ParseToken(tk.String()+"")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(2, newtoken)
	time.Sleep(time.Second * 60)
}
