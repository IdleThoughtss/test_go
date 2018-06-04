package wechat

import "time"

type Msg struct {
	From Contact
	To Contact
	Content string
	ContentType int
	CreateTime int
	Wx *Server
}

type SendMsg struct {
	ClientMsgID  int64
	Content      string
	FromUserName string
	LocalID      int64
	ToUserName   string
	Type         int
}

func (msg *Msg)Reply(content string)  {
	id := time.Now().UnixNano()
	 sendMsg := SendMsg{
	 	id,
	 	`content`,
	 	msg.Wx.user.UserName,
	 	id,
	 	msg.From.UserName,
	 	1,
	 }
	msg.Wx.send(sendMsg)
}

func (msg *Msg)init(wx *Server,msgItem AddMsg)  {
	msg.From,_ = wx.getUserInfo(msgItem.FromUserName)
	msg.To,_ = wx.getUserInfo(msgItem.ToUserName)
	msg.ContentType = msgItem.MsgType
	msg.CreateTime = msgItem.CreateTime
	msg.Wx = wx
	switch msgItem.MsgType {
	case 1:
		msg.Content = msgItem.Content
	case 34: // voice
		msg.Content = BaseUrl +`/cgi-bin/mmwebwx-bin/webwxgetvoice?msgid=`+ msgItem.MsgID
	case 3: // img
	case 47://表情包
		msg.Content = BaseUrl + `/cgi-bin/mmwebwx-bin/webwxgetmsgimg?&MsgID=` + msgItem.MsgID
	case 43: // video
		msg.Content = BaseUrl + `/cgi-bin/mmwebwx-bin/webwxgetvideo?msgid=` + msgItem.MsgID
	case 49://公众号 转账 乱七八糟
	case 51:
		//open dialog
	case 10000 : // redBag
		msg.Content = msgItem.Content
	}
}

// todo 转账
/*
<msg>
	<br/>
	<appmsg appid="" sdkver="">
		<br/>
		<title><![CDATA[å¾®ä¿¡è½¬è´¦]]></title>
		<br/>
		<des><![CDATA[æ”¶åˆ°è½¬è´¦1.00å…ƒã€‚å¦‚éœ€æ”¶é’±ï¼Œè¯·ç‚¹æ­¤å‡çº§è‡³æœ€æ–°ç‰ˆæœ¬]]></des>
		<br/>
		<action>
		</action>
		<br/>
		<type>2000</type>
		<br/>
		<content><![CDATA[]]></content>
		<br/>
		<url><![CDATA[https://support.weixin.qq.com/cgi-bin/mmsupport-bin/readtemplate?t=page/common_page__upgrade&amp;text=text001&amp;btn_text=btn_text_0]]>
		 </url>
		 <br/>
		<thumburl>
		<![CDATA[https://support.weixin.qq.com/cgi-bin/mmsupport-bin/readtemplate?t=page/common_page__upgrade&amp;text=text001&amp;btn_text=btn_text_0]]></thumburl>
		 <br/>
		<lowurl></lowurl>
		<br/>
		<extinfo>
		<br/>
		</extinfo>
		<br/>
		<wcpayinfo>
		<br/>
		<paysubtype>1</paysubtype>
		<br/>
		<feedesc>
		<![CDATA[ï¿¥1.00]]>
		</feedesc><br/>
		<transcationid>
		<![CDATA[100005020118060400065311429252186418]]>
		</transcationid>
		<br/>
		<transferid>
		<![CDATA[1000050201201806041000689109104]]>
		</transferid>
		<br/>
		<invalidtime>
		<![CDATA[1528210851]]>
		</invalidtime>
		<br/>
		<begintransfertime>
		<![CDATA[1528119051]]>
		</begintransfertime>
		<br/>
		<effectivedate>
		<![CDATA[1]]></effectivedate>
		<br/>
		<pay_memo><![CDATA[æµ‹è¯•]]>
		</pay_memo>
		<br/><br/>
		</wcpayinfo>
		<br/>
		</appmsg>
		<br/>
		</msg>
*/

/*
小程序
<?xml version="1.0"?><br/><msg><br/>	<appmsg appid="" sdkver="0"><br/>		<title>å…¨ä¸–ç•Œèªæ˜Žäººéƒ½åœ¨çŽ©çš„æ¸¸æˆï¼Œå¿«æ¥è¯•è¯•ï¼</title><br/>		<des /><br/>		<username /><br/>		<action>view</action><br/>		<type>33</type><br/>		<showtype>0</showtype><br/>		<content /><br/>		<url>https://mp.weixin.qq.com/mp/waerrpage?appid=wxbe77d50ac82c441d&amp;amp;type=upgrade&amp;amp;upgradetype=3#wechat_redirect</url><br/>		<lowurl /><br/>		<dataurl /><br/>		<lowdataurl /><br/>		<contentattr>0</contentattr><br/>		<streamvideo><br/>			<streamvideourl /><br/>			<streamvideototaltime>0</streamvideototaltime><br/>			<streamvideotitle /><br/>			<streamvideowording /><br/>			<streamvideoweburl /><br/>			<streamvideothumburl /><br/>			<streamvideoaduxinfo /><br/>			<streamvideopublishid /><br/>		</streamvideo><br/>		<canvasPageItem><br/>			<canvasPageXml><![CDATA[]]></canvasPageXml><br/>		</canvasPageItem><br/>		<appattach><br/>			<attachid /><br/>			<cdnthumburl>305c02010004553053020100020468d7e8b002033d11fe020433f516d202045b14f186042e6175706170706d73675f383032333539386663643365643964315f313532383039393230353039325f31353536380204010400030201000400</cdnthumburl><br/>			<cdnthumbmd5>1bbbd558ef85f10023dae4c5f5f4d89a</cdnthumbmd5><br/>			<cdnthumblength>146845</cdnthumblength><br/>			<cdnthumbheight>576</cdnthumbheight><br/>			<cdnthumbwidth>720</cdnthumbwidth><br/>			<cdnthumbaeskey>48c57264f846446b92c41f75f842ee55</cdnthumbaeskey><br/>			<aeskey>48c57264f846446b92c41f75f842ee55</aeskey><br/>			<encryver>1</encryver><br/>			<fileext /><br/>			<islargefilemsg>0</islargefilemsg><br/>		</appattach><br/>		<extinfo /><br/>		<androidsource>0</androidsource><br/>		<sourceusername></sourceusername><br/>		<sourcedisplayname>æœ€å¼ºå¼¹ä¸€å¼¹</sourcedisplayname><br/>		<commenturl /><br/>		<thumburl /><br/>		<mediatagname /><br/>		<messageaction><![CDATA[]]></messageaction><br/>		<messageext><![CDATA[]]></messageext><br/>		<emoticongift><br/>			<packageflag>0</packageflag><br/>			<packageid /><br/>		</emoticongift><br/>		<emoticonshared><br/>			<packageflag>0</packageflag><br/>			<packageid /><br/>		</emoticonshared><br/>		<designershared><br/>			<designeruin>0</designeruin><br/>			<designername>null</designername><br/>			<designerrediretcturl>null</designerrediretcturl><br/>		</designershared><br/>		<emotionpageshared><br/>			<tid>0</tid><br/>			<title>null</title><br/>			<desc>null</desc><br/>			<iconUrl>null</iconUrl><br/>			<secondUrl>null</secondUrl><br/>			<pageType>0</pageType><br/>		</emotionpageshared><br/>		<webviewshared><br/>			<shareUrlOriginal /><br/>			<shareUrlOpen /><br/>			<jsAppId /><br/>			<publisherId>wxapp_wxbe77d50ac82c441d</publisherId><br/>		</webviewshared><br/>		<template_id /><br/>		<md5>1bbbd558ef85f10023dae4c5f5f4d89a</md5><br/>		<weappinfo><br/>			<username></username><br/>			<appid>wxbe77d50ac82c441d</appid><br/>			<version>29</version><br/>			<type>2</type><br/>			<weappiconurl><![CDATA[http://mmbiz.qpic.cn/mmbiz_png/M3ZDbDAdgsYXjGG6jvRA243UVwKkJZxDNhXcX7FAH0aC4SnBdOgMicwqKf8GRUsOXrcHDiamTRS9jpAO6qgZ1oicQ/0?wx_fmt=png]]></weappiconurl><br/>			<shareId><![CDATA[1_wxbe77d50ac82c441d_23051855_1528120477_0]]></shareId><br/>			<appservicetype>4</appservicetype><br/>		</weappinfo><br/>		<statextstr /><br/>		<websearch><br/>			<rec_category>0</rec_category><br/>		</websearch><br/>	</appmsg><br/>	<fromusername></fromusername><br/>	<scene>0</scene><br/>	<appinfo><br/>		<version>1</version><br/>		<appname></appname><br/>	</appinfo><br/>	<commenturl></commenturl><br/></msg><br/>"


 */
type OtherMsg struct {

}

{{0 @68519bfab6585599a66c7e1619737d92 懒人闲思 /cgi-bin/mmwebwx-bin/webwxgeticon?seq=685124159&username=@68519bfab6585599a66c7e1619737d92&skey=@crypt_2984_6c70aebd4e0741fe520beaeaeb4f64b7 7 0 []  0 1 我想静静 0 0 LRXS lanrenxiansai   0 0 0 103527 江苏 南京  17 0  0 lld  0} {0 @4d1a2598c462d967cdb2a3d902713f024d8b364cecdd34df80dff6e1 静静 /cgi-bin/mmwebwx-bin/webwxgeticon?seq=700448348&username=@4d1a2598c462d967cdb2a3b31578a8d902713f024d8b364cecdddf80dff6e1&skey=@crypt_c0972984_6c70aebd4e0741fe520beaeaeb4f64b7 3 0 [] 静静宝贝 0 2 我听见了你的声音，也藏着颗不敢见的心。 0 0 JJ jingjing JJBB jingj上海 徐汇  17 0  0   0} 你什么时候睡 1 1528128379 0xc4200f0000}
