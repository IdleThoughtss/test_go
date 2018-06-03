package main

import (
	"github.com/IdleThoughtss/wechat"
)

func main(){
	server := wechat.Server{}
	server.SetHandler(handler)
	server.Start()

}

func handler(wx *wechat.Server,message wechat.AddMsg)  {

}

