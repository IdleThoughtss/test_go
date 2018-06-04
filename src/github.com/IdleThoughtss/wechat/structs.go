package wechat

import "encoding/xml"

type CookieInfo struct {
	XmlName xml.Name `xml:"error"`
	Ret int `xml:"ret"`
	Message  string  `xml:"message"`
	Skey   string  `xml:"skey"`
	Wxsid  string `xml:"wxsid"`
	Wxuin  string `xml:"wxuin"`
	Pass_ticket  string `xml:"pass_ticket"`
	Isgrayscale  string `xml:"isgrayscale"`
}

type ContactResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	MemberCount  int          `json:"MemberCount"`
	MemberList   []Contact    `json:"MemberList"`
	Seq          int          `json:"Seq"`
}

type BaseResponse struct {
	Ret int
	ErrMsg string `json:"ErrMsg"`
}

type InitResponse struct {
	BaseResponse BaseResponse `json:"BaseResponse"`
	Count int `json:"count"`
	ContactList []Contact `json:"ContactList"`
	ChatSet string `json:"ContactList"`
	ClickReportInterval int `json:"ClickReportInterval"`
	ClientVersion int `json:"ClientVersion"`
	GrayScale int `json:"GrayScale"`
	InviteStartCount int `json:"InviteStartCount"`
	MPSubscribeMsgCount int `json:"MPSubscribeMsgCount"`
	MPSubscribeMsgList []MPSubscribeMsg `json:"MPSubscribeMsgList"`
	SKey string `json:"SKey"`
	SyncKey SyncKey `json:"SyncKey"`
	SystemTime int `json:"SystemTime"`
	User User `json:"User"`
}

type SyncKey struct {
	Count int `json:"Count"`
	List []SyncKeyItem `json:"List"`
}

type  SyncKeyItem struct {
	Key  int `json:"Key"`
	Val int `json:"Val"`
}

type MPSubscribeMsg struct {
	MPArticleCount int `json:"MPArticleCount"`
	MPArticleList []MPArticle `json:"MPArticleList"`
	NickName string `json:"NickName"`
	Time int `json:"Time"`
	UserName string `json:"UserName"`
}

type MPArticle struct {
	Title  string `json:"Title"`
	Digest string `json:"Digest"`
	Cover  string `json:"Cover"`
	URL    string `json:"Url"`
}

type Member struct {
	AttrStatus int `json:"VerifyFlag"`
	DisplayName string `json:"VerifyFlag"`
	KeyWord string `json:"KeyWord"`
	MemberStatus int `json:"MemberStatus"`
	NickName string `json:"NickName"`
	PYInitial string `json:"PYInitial"`
	PYQuanPin string `json:"PYQuanPin"`
	RemarkPYInitial string `json:"RemarkPYInitial"`
	RemarkPYQuanPin string `json:"RemarkPYQuanPin"`
	Uin int `json:"Uin"`
	UserName string `json:"UserName"`

}

type Message struct {
	BaseResponse           BaseResponse  `json:"BaseResponse"`
	AddMsgCount            int           `json:"AddMsgCount"`
	AddMsgList             []AddMsg      `json:"AddMsgList"`
	ModContactCount        int           `json:"ModContactCount"`
	ModContactList         []Contact     `json:"ModContactList"`
	DelContactCount        int           `json:"DelContactCount"`
	DelContactList         []interface{} `json:"DelContactList"`
	ModChatRoomMemberCount int           `json:"ModChatRoomMemberCount"`
	ModChatRoomMemberList  []interface{} `json:"ModChatRoomMemberList"`
	Profile                Profile       `json:"Profile"`
	ContinueFlag           int           `json:"ContinueFlag"`
	SyncKey                SyncKey       `json:"SyncKey"`
	SKey                   string        `json:"SKey"`
	SyncCheckKey           SyncKey       `json:"SyncCheckKey"`
}

type Profile struct {
	BitFlag  int `json:"BitFlag"`
	UserName struct {
		Buff string `json:"Buff"`
	} `json:"UserName"`
	NickName struct {
		Buff string `json:"Buff"`
	} `json:"NickName"`
	BindUin   int `json:"BindUin"`
	BindEmail struct {
		Buff string `json:"Buff"`
	} `json:"BindEmail"`
	BindMobile struct {
		Buff string `json:"Buff"`
	} `json:"BindMobile"`
	Status            int    `json:"Status"`
	Sex               int    `json:"Sex"`
	PersonalCard      int    `json:"PersonalCard"`
	Alias             string `json:"Alias"`
	HeadImgUpdateFlag int    `json:"HeadImgUpdateFlag"`
	HeadImgURL        string `json:"HeadImgUrl"`
	Signature         string `json:"Signature"`
}

type Contact struct {
	Uin              int      `json:"Uin"`
	UserName         string   `json:"UserName"`
	NickName         string   `json:"NickName"`
	HeadImgURL       string   `json:"HeadImgUrl"`
	ContactFlag      int      `json:"ContactFlag"`
	MemberCount      int      `json:"MemberCount"`
	MemberList       []Member `json:"MemberList"`
	RemarkName       string   `json:"RemarkName"`
	HideInputBarFlag int      `json:"HideInputBarFlag"`
	Sex              int      `json:"Sex"`
	Signature        string   `json:"Signature"`
	VerifyFlag       int      `json:"VerifyFlag"`
	OwnerUin         int      `json:"OwnerUin"`
	PYInitial        string   `json:"PYInitial"`
	PYQuanPin        string   `json:"PYQuanPin"`
	RemarkPYInitial  string   `json:"RemarkPYInitial"`
	RemarkPYQuanPin  string   `json:"RemarkPYQuanPin"`
	StarFriend       int      `json:"StarFriend"`
	AppAccountFlag   int      `json:"AppAccountFlag"`
	Statues          int      `json:"Statues"`
	AttrStatus       int      `json:"AttrStatus"`
	Province         string   `json:"Province"`
	City             string   `json:"City"`
	Alias            string   `json:"Alias"`
	SnsFlag          int      `json:"SnsFlag"`
	UniFriend        int      `json:"UniFriend"`
	DisplayName      string   `json:"DisplayName"`
	ChatRoomID       int      `json:"ChatRoomId"`
	KeyWord          string   `json:"KeyWord"`
	EncryChatRoomID  string   `json:"EncryChatRoomId"`
	IsOwner          int      `json:"IsOwner"`
}

type AddMsg struct {
	MsgID                string        `json:"MsgId"`
	FromUserName         string        `json:"FromUserName"`
	ToUserName           string        `json:"ToUserName"`
	MsgType              int           `json:"MsgType"`
	Content              string        `json:"Content"`
	Status               int           `json:"Status"`
	ImgStatus            int           `json:"ImgStatus"`
	CreateTime           int           `json:"CreateTime"`
	VoiceLength          int           `json:"VoiceLength"`
	PlayLength           int           `json:"PlayLength"`
	FileName             string        `json:"FileName"`
	FileSize             string        `json:"FileSize"`
	MediaID              string        `json:"MediaId"`
	URL                  string        `json:"Url"`
	AppMsgType           int           `json:"AppMsgType"`
	StatusNotifyCode     int           `json:"StatusNotifyCode"`
	StatusNotifyUserName string        `json:"StatusNotifyUserName"`
	RecommendInfo        RecommendInfo `json:"RecommendInfo"`
	ForwardFlag          int           `json:"ForwardFlag"`
	AppInfo              AppInfo       `json:"AppInfo"`
	HasProductID         int           `json:"HasProductId"`
	Ticket               string        `json:"Ticket"`
	ImgHeight            int           `json:"ImgHeight"`
	ImgWidth             int           `json:"ImgWidth"`
	SubMsgType           int           `json:"SubMsgType"`
	NewMsgID             int64         `json:"NewMsgId"`
	OriContent           string        `json:"OriContent"`
	EncryFileName        string        `json:"EncryFileName"`
}

type AppInfo struct {
	AppID string `json:"AppID"`
	Type  int    `json:"Type"`
}

type RecommendInfo struct {
	UserName   string `json:"UserName"`
	NickName   string `json:"NickName"`
	QQNum      int    `json:"QQNum"`
	Province   string `json:"Province"`
	City       string `json:"City"`
	Content    string `json:"Content"`
	Signature  string `json:"Signature"`
	Alias      string `json:"Alias"`
	Scene      int    `json:"Scene"`
	VerifyFlag int    `json:"VerifyFlag"`
	AttrStatus int    `json:"AttrStatus"`
	Sex        int    `json:"Sex"`
	Ticket     string `json:"Ticket"`
	OpCode     int    `json:"OpCode"`
}

