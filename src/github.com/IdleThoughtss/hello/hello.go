package main

import (
	"github.com/IdleThoughtss/login"
	"fmt"
	//"regexp"
	"encoding/xml"
	"log"
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
	//readxml()
	err := login.GetQrcode()
	if err != nil {
		fmt.Print("request failed!")
		return
	}

}

func readxml()  {
	xmlL := []byte(`<?xml version="1.0" encoding="UTF-8"?>
	<error>
<ret>1</ret>
<message></message>
<skey>@crypt_c0972984_47384cd60eb3f598dc1a1dcf21916cbc</skey>
<wxsid>3i1uYbu1OwxVLJ5r</wxsid>
<wxuin>23051855</wxuin>
<pass_ticket>u8VmBMbOq4n7W5JPQ08bTd7mLt%2FK1%2FAKIolXRwtJwUo%3D</pass_ticket>
<isgrayscale>1</isgrayscale>
</error>
`)
	var result Err
	xml.Unmarshal(xmlL,&result)
	log.Println(result)

}