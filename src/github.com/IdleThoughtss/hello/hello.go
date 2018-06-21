package main

import (
	"github.com/IdleThoughtss/wechat"
	"os"
	"encoding/json"
	"fmt"
)

func main(){
	wx := wechat.NewInstance()
	wx.Start()
	fh,err := os.Create(`aaa`)
	if err != nil {
		fmt.Println(err)
	}
	c,_ := json.Marshal(wx.ContactList)
	if err != nil {
		fmt.Println(err)
	}
	n,_ := fh.Write(c)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(n)
	wx.HandleMessage(handler)


}

func handler(wx *wechat.Server,message wechat.Msg)  {
		if wechat.MSG_TEXT == message.Type {
			fmt.Println(message.From.NickName + `:` + message.Content)
		}
}

