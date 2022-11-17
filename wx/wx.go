package wx

import (
	"encoding/xml"
	"fmt"
	"goDingRobot/robot"
	"log"
	"net/http"
	"time"
)

var ServerWxPrefix = "/wx"
var appId string
var appSecret string

func NewWx(inAppId string, inAppSecret string) {
	appId = inAppId
	appSecret = inAppSecret
}

func checkConfig() bool {
	if appId == "" || appSecret == "" {
		return false
	}
	return true
}

func HandleWxRobot(w http.ResponseWriter, r *http.Request) {
	//uri := r.RequestURI
	if checkConfig() {
		fmt.Fprintf(w, "微信未配置！")
		return
	}
	handleMsg(w, r)
}

func handleMsg(w http.ResponseWriter, r *http.Request) {
	bodyByte := make([]byte, r.ContentLength)
	r.Body.Read(bodyByte)

	log.Println(string(bodyByte))

	receiveTextMsg := &ReceiveTextMsg{}
	_ = xml.Unmarshal(bodyByte, receiveTextMsg)

	log.Println(receiveTextMsg)

	msgType := receiveTextMsg.MsgType
	if "text" == msgType {
		content := receiveTextMsg.Content
		fromUserId := receiveTextMsg.FromUserName
		msgId := receiveTextMsg.MsgId
		respVo := robot.HandleContent(content, msgId, fromUserId)
		handleReturnTextMsg(w, respVo.Msg, receiveTextMsg)
	}
}

func handleReturnTextMsg(w http.ResponseWriter, content string, receiveTextMsg *ReceiveTextMsg) {
	unix := time.Now().Unix()
	var respReceiveTextMsg = &ReceiveTextMsg{
		ToUserName:   receiveTextMsg.FromUserName,
		FromUserName: receiveTextMsg.ToUserName,
		CreateTime:   unix,
		MsgType:      "text",
		Content:      content,
	}

	marshal, _ := xml.Marshal(respReceiveTextMsg)
	marshalStr := string(marshal)
	log.Println(marshalStr)
	fmt.Fprintf(w, marshalStr)
}
