package main

import (
	"github.com/IdleThoughtss/wechat"
	"fmt"
)

func main(){
	server := wechat.Server{}
	server.SetHandler(handler)
	server.Start()

}

func handler(wx *wechat.Server,message wechat.Msg)  {
		//message.Reply(`123`
		fmt.Println(message)
}

