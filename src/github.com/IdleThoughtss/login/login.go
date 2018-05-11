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
)

var config struct{
	uuid string
	redirectUrl string
}

var BaseRequest = make(map[string]interface{})
var pass_ticket  string
const(
	LoginUri = "https://login.weixin.qq.com"
	BaseUrl = "https://wx.qq.com"
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



func GetQrcode() (err error) {
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
			redirectUrl = string(urlR.FindSubmatch(body)[1]) + "&fun=new"
			config.redirectUrl = redirectUrl
			fmt.Println(redirectUrl)
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
	response,err := http.Get(url)
	defer response.Body.Close()
	if err !=nil {
		fmt.Println("request redirectUri failed")
	}
	content,_ :=  ioutil.ReadAll(response.Body)
	xmlprefix :=[]byte(`<?xml version="1.0" encoding="UTF-8"?>`)
	xmlful := append(xmlprefix,content...)
	fmt.Println(string(xmlful))
	var result  CookieInfo
	err = xml.Unmarshal(xmlful,&result)
	if err !=nil {
		fmt.Println(err)
	}
	rand.Seed(time.Now().UnixNano())
	deviceNum := rand.Intn(999999999) + rand.Intn(1) * 1000000000
	deviceId := fmt.Sprintf("%s%d","e",deviceNum)
	baseRequest := make(map[string]interface{})
	baseRequest["Sid"] = result.Wxsid
	baseRequest["Uin"] = result.Wxuin
	baseRequest["Skey"] = result.Pass_ticket
	baseRequest["DeviceID"] = deviceId
	BaseRequest["BaseRequest"] = baseRequest
	pass_ticket = result.Pass_ticket
	Init()
}

func Init(){
	timeStamp := time.Now().UnixNano() /1000
	uri :=fmt.Sprintf("%s/cgi-bin/mmwebwx-bin/webwxinit?r=%d&pass_ticket=%s",BaseUrl,timeStamp,pass_ticket)
	fmt.Println(uri)
	fmt.Println(BaseRequest)
	body,err :=httpPost(uri,BaseRequest)
	if err !=nil{
		fmt.Println("init failed")
		return
	}
	//fmt.Println(string(body))
	synsStatus()
}

func synsStatus(){

}
func httpGet(uri string,param map[string]string)(body []byte, err error){

	if err != nil {
		fmt.Println("marshal json data failed")
		return
	}
	response,err := http.Get(uri)
	if err !=nil {
		fmt.Println("http get failed")
	}
	defer response.Body.Close()
	body,err = ioutil.ReadAll(response.Body)
	return
}

func httpPost(uri string,data map[string]interface{}) (body []byte,err error) {
	jsonData,err := json.Marshal(data)
	fmt.Println(string(jsonData))
	if err != nil {
		fmt.Println("marshal json data failed")
		return
	}
	response,err := http.Post(uri,"application/json",bytes.NewReader(jsonData))
	if err !=nil {
		fmt.Println("post failed")
		return
	}
	defer response.Body.Close()
	body,err = ioutil.ReadAll(response.Body)
	return

}