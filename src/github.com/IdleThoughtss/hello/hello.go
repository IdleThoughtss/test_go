package main

import (
	"github.com/IdleThoughtss/login"
	"fmt"
	//"regexp"
)

func main() {
	err := login.GetQrcode()
	if err != nil {
		fmt.Print("request failed!")
		return
	}

}
