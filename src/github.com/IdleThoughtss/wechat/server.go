package wechat

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"os"
	"io"
	"runtime"
	"os/exec"
	"time"
	"regexp"
	"encoding/xml"
	"math/rand"
	"encoding/json"
	"bytes"
	"strconv"
	"net/http/cookiejar"
)

type Server struct{
	client http.Client
	uuid string
	redirectUrl string
	baseRequest map[string]string
	passTicket string
	syncKey SyncKey
	ContactList  map[string]Contact
	user User
	handler HandleFunc
}

type HandleFunc func (wx *Server,message Msg)

func (wx *Server) Start() (err error) {
	wx.init()
	url := LOGIN_HOST + "/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=" + timeNowStr(13)
	response,err :=wx.httpGet(url)
	if err != nil {
		return
	}
	uuid := string(response[50:62])
	wx.uuid = uuid
	wx.showQrcode()
	return
}

func (wx *Server) init()  {
	cookieJar ,_ :=  cookiejar.New(nil)
	wx.client = http.Client{
		Jar: cookieJar,
		CheckRedirect:nil,
	}
	wx.baseRequest = make(map[string]string)
	wx.ContactList = make(map[string]Contact)
}

func (wx *Server) showQrcode() {
	imgUrl := LOGIN_HOST + "/qrcode/" + wx.uuid
	response , err := http.Get(imgUrl)
	if err != nil{
		fmt.Println("open Qrcode failed")
		return
	}
	defer response.Body.Close()
	file ,err := os.Create("../tmp/qrcode.jpg")
	if err != nil {
		fmt.Println("create image file failed")
		return
	}
	defer file.Close()
	_,err = io.Copy(file,response.Body)
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
	wx.listenScan()
}

func (wx *Server) listenScan(){
	timeNow := timeNowStr(13)
	tip := "1"
	var response *http.Response
	var err error
	var code , url, redirectUrl string
	codeR,_ := regexp.Compile("window.code=([0-9]*);")
	urlR,_ := regexp.Compile(`window.redirect_uri="(.*)";`)
	for {
		url = LOGIN_HOST + "/cgi-bin/mmwebwx-bin/login?loginicon=true&uuid="+ wx.uuid +"&tip="+ tip +"&_=" + timeNow
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
			wx.redirectUrl = redirectUrl
			os.Remove("../tmp/qrcode.jpg")
			wx.login()
			return
		case "408":
		default:
			time.Sleep(time.Second *1)
		}
	}
	response.Body.Close()
}

func (wx *Server) login() {
	url := wx.redirectUrl
	response,err := wx.httpGet(url)

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
	fmt.Println(wx.baseRequest)
	wx.baseRequest["Sid"] = result.Wxsid
	wx.baseRequest["Uin"] = result.Wxuin
	wx.baseRequest["Skey"] = result.Skey
	wx.baseRequest["DeviceID"] = deviceId
	wx.passTicket = result.Pass_ticket
	wx.initInfo()
}

func (wx *Server)  initInfo(){
	timeStamp := time.Now().UnixNano() /1000
	uri :=fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/webwxinit?r=%d&pass_ticket=%s",BaseUrl,timeStamp,wx.passTicket)
	BaseRequest := make(map[string]interface{})
	BaseRequest["BaseRequest"] = wx.baseRequest
	body,err :=wx.httpPost(uri,BaseRequest)
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
	wx.pushContact(res.ContactList)
	wx.syncKey = res.SyncKey
	wx.user = res.User
	wx.getContact()
	wx.syncStatus()
}

func (wx *Server) getContact()  {
	uri := `https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxgetcontact`
	query := make(map[string]string)
	query["skey"] = wx.baseRequest["Skey"]
	query["lang"] = "zh_CN"
	query["r"] = timeNowStr(13)
	query["seq"] = "0"
	uri += formatQueryString(query)
	response,err :=wx.httpGet(uri)
	if err != nil {
		fmt.Println("getContact failed !")
	}
	var contactResponse ContactResponse
	if err := json.Unmarshal(response,&contactResponse); err !=nil{
		fmt.Println("Unmarshal ContactResponse failed!")
	}
	if contactResponse.MemberCount >0 {
		wx.pushContact(contactResponse.MemberList)
	}

}

func  (wx *Server) pushContact(list []Contact){
	for _,item := range list{
		if _,ok := wx.ContactList[item.UserName];!ok{
			wx.ContactList[item.UserName] = item
		}
	}
}

func (wx *Server) syncStatus(){
	uri := "https://webpush.wx.qq.com/cgi-bin/mmwebwx-bin/synccheck"
	query := make(map[string]string)
	query["skey"] = wx.baseRequest["Skey"]
	query["sid"] = wx.baseRequest["Sid"]
	query["uin"] = wx.baseRequest["Uin"]
	IdStr := strconv.Itoa(rand.Int())
	query["deviceid"] = IdStr[2:17]
	firstTIme := time.Now().UnixNano() / 1000
	var timeNow string

	for {
		timeNow = timeNowStr(13)
		query["synckey"] = wx.syncKeyStr()
		query["_"] = strconv.FormatInt(firstTIme + 1 ,10)
		query["r"] = timeNow
		uri += formatQueryString(query)
		res,_ := wx.httpGet(uri)
		rule,_ :=regexp.Compile(`retcode:"([0-9]*)",selector:"([0-9]*)"`)
		match := rule.FindSubmatch(res)
		retcode := string( match[1])
		selector :=string( match[2])

		if retcode != "0"{
			fmt.Println("sync Failed")
			switch retcode {
			case "1101":
				fmt.Println("logout from mobile")
			}
			return
		}
		switch selector {
		case "2","4","6":
			wx.getMessage()
		case "0":
			time.Sleep(time.Second)

		}
		time.Sleep(time.Second)
	}

}


func (wx *Server) syncKeyStr () (keyStr string) {
	delimiter := ""
	for _,keyItem := range wx.syncKey.List {
		keyStr += fmt.Sprintf(`%s%d_%d`,delimiter,keyItem.Key,keyItem.Val)
		delimiter = "|"
	}
	return
}

func (wx *Server) getMessage()  {
	uri := "https://wx.qq.com/cgi-bin/mmwebwx-bin/webwxsync"
	query := make(map[string]string)
	query["skey"] = wx.baseRequest["Skey"]
	query["sid"] = wx.baseRequest["Sid"]
	query["pass_ticket"] = wx.passTicket
	query["lang"] = "zh_CN"
	uri += formatQueryString(query)
	requestBody := make(map[string]interface{})
	requestBody["BaseRequest"] = wx.baseRequest
	requestBody["SyncKey"] = wx.syncKey
	requestBody["rr"] = timeNowStr(10)
	response,err :=wx.httpPost(uri,requestBody)
	if err !=nil{
		fmt.Println("get message failed")
	}
	msg := Message{}
	if err := json.Unmarshal(response,&msg) ;err !=nil {
		fmt.Println("Unmarshal message json failed!")
	}
	wx.handle(msg)
}

func (wx *Server) handle(message Message)  {
	if message.BaseResponse.Ret !=0 {
		fmt.Println(`message error %s`,message.BaseResponse.ErrMsg)
	}
	wx.syncKey = message.SyncKey
	if message.ModContactCount > 0{
		wx.pushContact(message.ModContactList)
	}

	for _,msgItem := range message.AddMsgList {
			var msg Msg
			msg.init(wx,msgItem)
			wx.handler(wx,msg)
	}
}

func (wx *Server)send(msg SendMsg) (response []byte,err error) {
	uri := `/cgi-bin/mmwebwx-bin/webwxsendmsg?pass_ticket=` + wx.passTicket
	data := make(map[string]interface{})
	data[`BaseRequest`] = wx.baseRequest
	data[`Msg`] = msg
	response,err = wx.httpPost(uri,data)
	return
}

func (wx *Server) SetHandler (handler HandleFunc) {
	wx.handler = handler
}

func (wx *Server) getUserInfo (userName string)(contact Contact ,ok bool){
	contact,ok = wx.ContactList[userName]
	return
}

func (wx *Server) httpGet(uri string)(body []byte, err error){

	req,_ := http.NewRequest("GET",uri,nil)


	req.Header.Add("Referer", Referer)
	req.Header.Add("User-agent", UserAgent)
	response,err := wx.client.Do(req)
	if err !=nil {
		fmt.Println("http get failed")
		return
	}
	defer response.Body.Close()
	body,err = ioutil.ReadAll(response.Body)
	return
}

func  (wx *Server) httpPost(uri string,data map[string]interface{}) (body []byte,err error) {
	jsonData,err := json.Marshal(data)

	if err != nil {
		fmt.Println("marshal json data failed")
		return
	}
	req,_ := http.NewRequest("POST",uri,bytes.NewReader(jsonData))
	req.Header.Add("Referer", Referer)
	req.Header.Add("User-agent", UserAgent)

	response,err := wx.client.Do(req)
	if err !=nil {
		fmt.Println("post failed")
		return
	}
	defer response.Body.Close()
	body,err = ioutil.ReadAll(response.Body)
	return

}