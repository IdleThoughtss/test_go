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
	fh,_ := os.Create(`tmp.aaa`)
	c,_ := json.Marshal(wx.ContactList)
	n,_ := fh.Write(c)
	fmt.Print(n)
	//wx.HandleMessage(handler)


}

func handler(wx *wechat.Server,message wechat.Msg)  {
		if message.From.RemarkName == `静静宝贝` {
			message.Reply(`你别回来了！`)
		}
}

