package main

import (
	"github.com/IdleThoughtss/wechat"
	"fmt"
)

func main(){
	wx := wechat.NewInstance()
	wx.Start()
	wx.HandleMessage(handler)

}

func handler(wx *wechat.Server,message wechat.Msg)  {
	fmt.Println(message.From.NickName + `:` + message.Content)	
}
