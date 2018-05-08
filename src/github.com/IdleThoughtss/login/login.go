package login

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"os"
	"io"
	"os/exec"
)

func GetQrcode()(imgUrl string,err error) {
	url := "https://login.wx.qq.com/jslogin?appid=wx782c26e4c19acffb&redirect_uri=https%3A%2F%2Fwx.qq.com%2Fcgi-bin%2Fmmwebwx-bin%2Fwebwxnewloginpage&fun=new&lang=zh_CN&_=1525777185095"
	response,err :=http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return
	}
	body,_ := ioutil.ReadAll(response.Body)
	imgUri := string(body[50:62])
	imgUrl = "https://login.weixin.qq.com/qrcode/" + imgUri
	fmt.Println(imgUrl)
	saveQrcode(imgUrl)
	return
}

func saveQrcode(imgUrl string) {
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


}