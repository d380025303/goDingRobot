package ding

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goDingRobot/robot"
	"net/http"
)

var ServeDingRobotPrefix = "/dingRobot"
var appKey string
var appSecret string
var token = ""

func NewDingRobot(inAppKey string, inAppSecret string) {
	appKey = inAppKey
	appKey = inAppSecret
}

func HmacSha256(key, data string) []byte {
	keys := []byte(key)
	h := hmac.New(sha256.New, keys)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func checkSign(r *http.Request) bool {
	header := r.Header
	timestamp := header.Get("timestamp")
	stringToSign := timestamp + "\n" + appKey
	b := HmacSha256(appSecret, stringToSign)
	mySign := base64.StdEncoding.EncodeToString(b)
	sign := header.Get("sign")
	return mySign == sign
}

func checkConfig() bool {
	if appKey == "" || appSecret == "" {
		return false
	}
	return true
}

func HandleDingRobot(w http.ResponseWriter, r *http.Request) {
	//uri := r.RequestURI
	if !checkConfig() {
		fmt.Fprintf(w, "钉钉未配置")
		return
	}
	if checkSign(r) {
		fmt.Fprintf(w, "签名校验失败！")
		return
	}
	handleMsg(w, r)
}

func handleMsg(w http.ResponseWriter, r *http.Request) {
	bodyByte := make([]byte, r.ContentLength)
	r.Body.Read(bodyByte)

	msgReceive := &MsgReceive{}
	_ = json.Unmarshal(bodyByte, msgReceive)

	if "1" == msgReceive.ConversationType {
		msgType := msgReceive.MsgType
		if "text" == msgType {
			text := msgReceive.Text
			content := text.Content
			respVo := robot.HandleContent(content, msgReceive.MsgId, msgReceive.SenderStaffId)
			handleReturnTextMsg(w, respVo.Msg, msgReceive)
		} else {
			handleReturnTextMsg(w, "只支持Text！", msgReceive)
		}
	} else {
		fmt.Fprintf(w, "暂不支持群聊哦！")
	}
}

func handleReturnTextMsg(w http.ResponseWriter, content string, receive *MsgReceive) {
	v := fmt.Sprintf(`{
	"msgtype": "text",
	"text": {
		"content": "%v"
	},
	"at": {
		"atUserIds": ["%v"]
	}
}`, content, receive.SenderStaffId)
	fmt.Println(v)
	fmt.Fprintf(w, v)
}

type MsgReceive struct {
	ConversationId            string `json:"conversationId"`
	AtUsers                   []User `json:"atUsers"`
	ChatBotCorpId             string `json:"chatbotCorpId"`
	ChatBotUserId             string `json:"chatbotUserId"`
	MsgId                     string `json:"msgId"`
	SenderNick                string `json:"senderNick"`
	IsAdmin                   bool   `json:"isAdmin"`
	SenderStaffId             string `json:"senderStaffId"`
	SessionWebhookExpiredTime int64  `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64  `json:"createAt"`
	SenderCorpId              string `json:"SenderCorpId"`
	ConversationType          string `json:"conversationType"`
	SenderId                  string `json:"senderId"`
	ConversationTitle         string `json:"conversationTitle"`
	IsInAtList                bool   `json:"isInAtList"`
	SessionWebhook            string `json:"sessionWebhook"`
	Text                      Text   `json:"text"`
	MsgType                   string `json:"msgtype"`
}

type User struct {
	DingTalkId string `json:"dingtalkId"`
	StaffId    string `json:"staffId"`
}

type Text struct {
	Content string `json:"content"`
}
