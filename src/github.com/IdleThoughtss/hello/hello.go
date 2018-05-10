package main

import (
	"github.com/IdleThoughtss/login"
	"fmt"
	//"regexp"
	"encoding/xml"
	"encoding/json"
)



type Err struct {
	XmlName xml.Name `xml:"error"`
	Ret int `xml:"ret"`
	Message  string  `xml:"message"`
	Skey   string  `xml:"skey"`
	Wxsid  string `xml:"wxsid"`
	Wxuin  string `xml:",wxuin"`
	Pass_ticket  string `xml:"pass_ticket"`
	Isgrayscale  string `xml:"isgrayscale"`
}

func main() {
	//readjson()

	//return
	err := login.GetQrcode()
	if err != nil {
		fmt.Print("request failed!")
		return
	}

}


func readjson()  {

jsonData := `{"BaseRequest":{"Uin":0,"Sid":0},"Count":1,"List":[{"Type":1,"Text":"/cgi-bin/mmwebwx-bin/login, Second Request Start, uuid: 454d958c7f6243"}]}`
var v map[string]interface{}
json.Unmarshal([]byte(jsonData),&v)
fmt.Println(v)
}
