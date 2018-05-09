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
)

var config struct{
	uuid string
	redirectUrl string
}

type cookieInfo struct {
	e errorMsg `xml:"error"`
}
type errorMsg struct {
	ret int `xml:"ret"`
	message  string  `xml:"message"`
	skey   string  `xml:"skey"`
	wxsid  string `xml:"wxsid"`
	wxuin  string `xml:"wxuin"`
	pass_ticket  string `xml:"pass_ticket"`
	isgrayscale  string `xml:"isgrayscale"`
}


func GetQrcode() (err error) {
	url := "https://login.wx.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=1525777185095"
	response,err :=http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body,_ := ioutil.ReadAll(response.Body)
	uuid := string(body[50:62])
	config.uuid = uuid
	ShowQrcode()
	return
}

func ShowQrcode() {
	imgUrl :=  "https://login.weixin.qq.com/qrcode/" + config.uuid
	response , err := http.Get(imgUrl)
	defer response.Body.Close()
	if err != nil{
		fmt.Println("open Qrcode failed")
		return
	}
	file ,err := os.Create("/tmp/qrcode.png")
	defer file.Close()
	if err != nil {
		fmt.Println("create image file failed")
		return
	}
	_,err = io.Copy(file,response.Body)
	if err != nil {
		fmt.Println("write Qrcode failed")
		return
	}
	cmd := exec.Command("open","/tmp/qrcode.png")
	err = cmd.Run()
	if err != nil {
		fmt.Println("open /tmp/qrcode.png  failed!")
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
	v := cookieInfo{}
	err = xml.Unmarshal(xmlful,&v)
	if err !=nil {
		fmt.Println(err)
	}
	fmt.Println(v)

}