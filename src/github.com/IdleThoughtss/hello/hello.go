package main

import (
	"github.com/IdleThoughtss/login"
	"fmt"
)

func main() {
	_ ,err := login.GetQrcode()
	if err != nil {
		fmt.Print("request failed!")
		return
	}
	fmt.Println("Scan the qrcode please!")

}
