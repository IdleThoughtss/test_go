package wechat

type Msg struct {
	fromUser interface{}
	toUser interface{}
	content interface{}
}

func (msg *Msg)reply(wx Server)  {
	wx.send(msg)
}