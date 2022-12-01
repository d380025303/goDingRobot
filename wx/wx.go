package wx

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"goDingRobot/robot"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

var ServerWxPrefix = "/wx"
var appId string
var appSecret string
var token string

func NewWx(inAppId string, inAppSecret string, inToken string) {
	appId = inAppId
	appSecret = inAppSecret
	token = inToken
}

func checkConfig() bool {
	if appId == "" || appSecret == "" {
		log.Println(fmt.Sprintf("微信未配置, appId: %v, appSecret: %v", appId, appSecret))
		return false
	}
	return true
}

func handleCheck(signature, timestamp, nonce string) bool {
	arr := []string{token, timestamp, nonce}
	sort.Strings(arr)
	h := sha1.New()
	h.Write([]byte(strings.Join(arr, "")))
	bs := h.Sum(nil)
	if hex.EncodeToString(bs) == signature {
		return true
	} else {
		return false
	}
}

func HandleWxRobot(w http.ResponseWriter, r *http.Request) {
	//uri := r.RequestURI
	if !checkConfig() {
		fmt.Fprintf(w, "微信未配置！")
		return
	}
	_ = r.ParseForm()
	signature := r.Form.Get("signature")
	timestamp := r.Form.Get("timestamp")
	nonce := r.Form.Get("nonce")
	echostr := r.Form.Get("echostr")
	// 签名验证
	if handleCheck(signature, timestamp, nonce) {
		if echostr != "" {
			fmt.Fprintf(w, echostr)
		} else {
			handleMsg(w, r)
		}
		return
	} else {
		log.Println("验证失败...")
	}
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
