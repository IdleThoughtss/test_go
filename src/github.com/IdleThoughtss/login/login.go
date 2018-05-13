package login

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"os"
	"io"
	"os/exec"
	"time"
	"strconv"
	"regexp"
	"encoding/xml"
	"runtime"
	"math/rand"
	"encoding/json"
	"bytes"
	"net/url"
	"net/http/cookiejar"
	"sync"
)

var config struct{
	uuid string
	redirectUrl string
}

var baseRequest = make(map[string]string)
var pass_ticket  string
var syncKey SyncKey
var httpClient http.Client
var fileLocker = sync.Mutex{}
var contactList = make(map[string]Contact)
var user User

const(
	LoginUri = "https://login.weixin.qq.com"
	BaseUrl = "https://wx.qq.com"
	Referer = "https://wx.qq.com/?&lang=zh_CN"
	UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/5"
)

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

type User struct {

	Uin               int    `json:"Uin"`
	UserName          string `json:"UserName"`
	NickName          string `json:"NickName"`
	HeadImgURL        string `json:"HeadImgUrl"`
	RemarkName        string `json:"RemarkName"`
	PYInitial         string `json:"PYInitial"`
	PYQuanPin         string `json:"PYQuanPin"`
	RemarkPYInitial   string `json:"RemarkPYInitial"`
	RemarkPYQuanPin   string `json:"RemarkPYQuanPin"`
	HideInputBarFlag  int    `json:"HideInputBarFlag"`
	StarFriend        int    `json:"StarFriend"`
	Sex               int    `json:"Sex"`
	Signature         string `json:"Signature"`
	AppAccountFlag    int    `json:"AppAccountFlag"`
	VerifyFlag        int    `json:"VerifyFlag"`
	ContactFlag       int    `json:"ContactFlag"`
	WebWxPluginSwitch int    `json:"WebWxPluginSwitch"`
	HeadImgFlag       int    `json:"HeadImgFlag"`
	SnsFlag           int    `json:"SnsFlag"`
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

//type ContactList struct {
//	Alias string `json:"Alias"`
//	AppAccountFlag int `json:"AppAccountFlag"`
//	AttrStatus int `json:"AttrStatus"`
//	ChatRoomId int `json:"ChatRoomId"`
//	City string `json:"City"`
//	ContactFlag int `json:"ContactFlag"`
//	DisplayName string `json:"DisplayName"`
//	EncryChatRoomId string `json:"EncryChatRoomId"`
//	HeadImgUrl string `json:"HeadImgUrl"`
//	HideInputBarFlag int `json:"HideInputBarFlag"`
//	IsOwner int `json:"IsOwner"`
//	KeyWord string `json:"KeyWord"`
//	MemberCount int `json:"MemberCount"`
//	MemberList []Member `json:"MemberList"`
//	NickName string `json:"NickName"`
//	OwnerUin int `json:"OwnerUin"`
//	PYInitial string `json:"PYInitial"`
//	PYQuanPin string `json:"PYQuanPin"`
//	Province string `json:"Province"`
//	RemarkName string `json:"RemarkName"`
//	RemarkPYInitial string `json:"RemarkPYInitial"`
//	RemarkPYQuanPin string `json:"RemarkPYQuanPin"`
//	Sex  int `json:"Sex"`
//	Signature string `json:"Signature"`
//	SnsFlag int `json:"SnsFlag"`
//	StarFriend int `json:"StarFriend"`
//	Statues int `json:"Statues"`
//	Uin int `json:"Uin"`
//	UniFriend int `json:"UniFriend"`
//	UserName string `json:"UserName"`
//	VerifyFlag int `json:"VerifyFlag"`
//}

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

func GetQrcode() (err error) {
	initHttpClient()
	url := "https://login.wx.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=1525777185095"
	response,err :=http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	body,_ := ioutil.ReadAll(response.Body)
	uuid := string(body[50:62])
	config.uuid = uuid
	ShowQrcode()
	return
}

func ShowQrcode() {
	imgUrl :=  "https://login.weixin.qq.com/qrcode/" + config.uuid
	response , err := http.Get(imgUrl)
	if err != nil{
		fmt.Println("open Qrcode failed")
		return
	}
	defer response.Body.Close()
	file ,err := os.Create("../tmp/qrcode.jpg")
	fileLocker.Lock()
	if err != nil {
		fmt.Println("create image file failed")
		return
	}
	defer file.Close()
	_,err = io.Copy(file,response.Body)
	fileLocker.Unlock()
	if err != nil {
		fmt.Println("write Qrcode failed")
		return
	}
	system := runtime.GOOS
	var cmdName string
	if system == "linux" {
		// ubuntu
		cmdName = "eog"
	}else if system == "darwin" {
		// mac os
		cmdName = "open"
	}
	 cmd := exec.Command(cmdName,"../tmp/qrcode.jpg")
	err = cmd.Start()
	if err != nil {
		fmt.Println("open ../tmp/qrcode.jpg  failed!")
		return
	}
	fmt.Println("scan the qrcode please!")
	listenScan()
}

func listenScan(){
	timeNow := time.Now().UnixNano() / 1000000
	tip := "1"
	var response *http.Response
	var err error
	var code , url, redirectUrl string
	codeR,_ := regexp.Compile("window.code=([0-9]*);")
	urlR,_ := regexp.Compile(`window.redirect_uri="(.*)";`)
	for {
		url = "https://login.wx.qq.com/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid="+ config.uuid +"&tip="+ tip +"&_=" + strconv.FormatInt(timeNow,10)
		response,err = http.Get(url)
		if err != nil {
			fmt.Println("request login status failed")
			continue
		}
		body,_ :=ioutil.ReadAll(response.Body)
		code = string(codeR.FindSubmatch(body)[1])
		switch code {
		case "201":
			tip = "0"
			fmt.Println("please confirm")
		case "200":
			fmt.Println(`login success`)
			redirectUrl = string(urlR.FindSubmatch(body)[1]) + "&fun=new"
			config.redirectUrl = redirectUrl
			os.Remove("../tmp/qrcode.jpg")
			login()
			return
		case "408":
		default:
			time.Sleep(time.Second *1)
		}
	}
	response.Body.Close()
}

func login() {
	url := config.redirectUrl
	response,err := httpGet(url)

	if err !=nil {
		fmt.Println("request redirectUri failed")
	}
	xmlprefix :=[]byte(`<?xml version="1.0" encoding="UTF-8"?>`)
	xmlful := append(xmlprefix,response...)

	var result  CookieInfo
	err = xml.Unmarshal(xmlful,&result)
	if err !=nil {
		fmt.Println(err)
	}
	rand.Seed(time.Now().UnixNano())
	deviceNum := rand.Intn(999999999) + rand.Intn(1) * 1000000000
	deviceId := fmt.Sprintf("%s%d","e",deviceNum)
	baseRequest["Sid"] = result.Wxsid
	baseRequest["Uin"] = result.Wxuin
	baseRequest["Skey"] = result.Skey
	baseRequest["DeviceID"] = deviceId


	pass_ticket = result.Pass_ticket
	Init()
}

func Init(){
	timeStamp := time.Now().UnixNano() /1000
	uri :=fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/webwxinit?r=%d&pass_ticket=%s",BaseUrl,timeStamp,pass_ticket)
	BaseRequest := make(map[string]interface{})
	BaseRequest["BaseRequest"] = baseRequest
	body,err :=httpPost(uri,BaseRequest)
	if err !=nil{
		fmt.Println("get init info failed")
		return
	}
	var res InitResponse
	if err :=json.Unmarshal(body,&res);err != nil{
		fmt.Println("Unmarshal json failed")
		return
	}

	if res.BaseResponse.Ret !=0 {
		fmt.Println("init failed")
		return
	}
	pushContact(res.ContactList)
	syncKey = res.SyncKey
	user = res.User
	getContact()
	syncStatus()
}



func syncStatus(){
	uri := "https://webpush.wx.qq.com/cgi-bin/mmwebwx-bin/synccheck"
	query := make(map[string]string)
	query["skey"] = baseRequest["Skey"]
	query["sid"] = baseRequest["Sid"]
	query["uin"] = baseRequest["Uin"]
	firstTIme := time.Now().UnixNano() / 1000
	var timeNow string
	
	for {
		timeNow = timeNowStr(13)
		query["synckey"] = syncKeyStr()
		IdStr := strconv.Itoa(rand.Int())
		query["deviceid"] = IdStr[2:17]
		query["_"] = strconv.FormatInt(firstTIme + 1 ,10)
		query["r"] = timeNow
		uri += formatQueryString(query)
		res,_ := httpGet(uri)
		rule,_ :=regexp.Compile(`retcode:"([0-9]*)",selector:"([0-9]*)"`)
		match := rule.FindSubmatch(res)
		retcode := string( match[1])
		selector :=string( match[2])

		if retcode != "0"{
			fmt.Println("sync Failed")
			fmt.Println(retcode)
			return
		}
		switch selector {
		case "2","4","6":
			//fmt.Println(syncKeyStr())
			getMessage()
		case "0":
			time.Sleep(time.Second)

		}
		time.Sleep(time.Second)
	}

}

func getMessage()  {
	uri := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsync"
	query := make(map[string]string)
	query["skey"] = baseRequest["Skey"]
	query["sid"] = baseRequest["Sid"]
	query["pass_ticket"] = pass_ticket
	query["lang"] = "zh_CN"
	uri += formatQueryString(query)
	requestBody := make(map[string]interface{})
	requestBody["BaseRequest"] = baseRequest
	requestBody["SyncKey"] = syncKey
	requestBody["rr"] = timeNowStr(10)
	response,err :=httpPost(uri,requestBody)
	if err !=nil{
		fmt.Println("get message failed")
	}
	msg := Message{}
	if err := json.Unmarshal(response,&msg) ;err !=nil {
		fmt.Println("Unmarshal message json failed!")
	}
	handle(msg)
}

func handle(message Message)  {
	if message.BaseResponse.Ret !=0 {
		fmt.Println(`message error %s`,message.BaseResponse.ErrMsg)
	}
	syncKey = message.SyncKey
	if message.ModContactCount > 0{
		pushContact(message.ModContactList)
	}

	for _,msgItem := range message.AddMsgList {
		user,ok := getUserInfo(msgItem.FromUserName)
		if ok {
			fmt.Println(user.NickName + `:` + msgItem.Content)

		}
	}
}

func formatQueryString(query map[string]string) string {

		var q = url.Values{}
		for key,val := range query {
			q.Add(key,val)
		}
		return "?" + q.Encode()
}

func getContact()  {
	uri := `https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact`
	query := make(map[string]string)
	query["skey"] = baseRequest["Skey"]
	query["lang"] = "zh_CN"
	query["r"] = timeNowStr(13)
	query["seq"] = "0"
	uri += formatQueryString(query)
	response,err :=httpGet(uri)
	if err != nil {
		fmt.Println("getContact failed !")
	}
	var contactResponse ContactResponse
	if err := json.Unmarshal(response,&contactResponse); err !=nil{
		fmt.Println("Unmarshal ContactResponse failed!")
	}
	if contactResponse.MemberCount >0 {
		 pushContact(contactResponse.MemberList)
	}

}

func pushContact(list []Contact){
	for _,item := range list{
		if _,ok := contactList[item.UserName];!ok{
			contactList[item.UserName] = item
		}
	}
}

func getUserInfo (userName string)(contact Contact ,ok bool){
	contact,ok = contactList[userName]
	return
}

func timeNowStr (len int) (timeStr string) {
	if len == 10{
		timeStr = fmt.Sprintf(`%d`, time.Now().Unix())
	}
	if len ==13 {
		timeStr = fmt.Sprintf(`%d`, time.Now().UnixNano() / 1000)
	}
	return
}

func syncKeyStr () (keyStr string) {
	delimiter := ""
	for _,keyItem := range syncKey.List {
		keyStr += fmt.Sprintf(`%s%d_%d`,delimiter,keyItem.Key,keyItem.Val)
		delimiter = "|"
	}
	return
}

func initHttpClient()  {
	cookieJar ,_ :=  cookiejar.New(nil)
	client := http.Client{
		Jar: cookieJar,
		CheckRedirect:nil,
	}
	httpClient = client
}
func httpGet(uri string)(body []byte, err error){

	req,_ := http.NewRequest("GET",uri,nil)

	httpClient.Do(req)
	req.Header.Add("Referer", Referer)
	req.Header.Add("User-agent", UserAgent)
	response,err := httpClient.Do(req)
	if err !=nil {
		fmt.Println("http get failed")
		return
	}
	defer response.Body.Close()
	body,err = ioutil.ReadAll(response.Body)
	return
}

func httpPost(uri string,data map[string]interface{}) (body []byte,err error) {
	jsonData,err := json.Marshal(data)

	if err != nil {
		fmt.Println("marshal json data failed")
		return
	}
	req,_ := http.NewRequest("POST",uri,bytes.NewReader(jsonData))
	req.Header.Add("Referer", Referer)
	req.Header.Add("User-agent", UserAgent)

	response,err := httpClient.Do(req)
	if err !=nil {
		fmt.Println("post failed")
		return
	}
	defer response.Body.Close()
	body,err = ioutil.ReadAll(response.Body)
	return

}