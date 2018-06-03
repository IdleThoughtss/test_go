package wechat

import (
	"fmt"
	"time"
	"net/url"
)


func timeNowStr (len int) (timeStr string) {
	if len == 10{
		timeStr = fmt.Sprintf(`%d`, time.Now().Unix())
	}
	if len ==13 {
		timeStr = fmt.Sprintf(`%d`, time.Now().UnixNano() / 1000)
	}
	return
}

func formatQueryString(query map[string]string) string {

	var q = url.Values{}
	for key,val := range query {
		q.Add(key,val)
	}
	return "?" + q.Encode()
}

