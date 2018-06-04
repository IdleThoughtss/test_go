package wechat

type Msg struct {
	FromUser *User
	ToUser *User
	Content string
	ContentType string
}

func (msg *Msg)Reply(wx *Server)  {
	wx.send(msg)
}

func (msg *Msg)From(wx *Server){

}