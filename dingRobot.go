package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goDingRobot/db"
	"net/http"
	"strings"
)

var dingRobotPrefix string = "/dingRobot"

func HmacSha256(key, data string) []byte {
	keys := []byte(key)
	h := hmac.New(sha256.New, keys)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func checkSign(r *http.Request) bool {
	header := r.Header
	timestamp := header.Get("timestamp")
	stringToSign := timestamp + "\n" + MyGlobalConfig.Ding.AppKey
	b := HmacSha256(MyGlobalConfig.Ding.AppSecret, stringToSign)
	mySign := base64.StdEncoding.EncodeToString(b)
	sign := header.Get("sign")
	return mySign == sign
}

func handleDingRobot(w http.ResponseWriter, r *http.Request) {
	//uri := r.RequestURI
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
			saveArr := []string{"-s", "~s", "-save", "~save", "s", "save"}
			for _, v := range saveArr {
				if strings.HasPrefix(content, v) {
					// 保存
					msgReceive.Text.Content = strings.TrimSpace(strings.TrimPrefix(content, v))
					handleReturnMsg(w, saveMessage(msgReceive), msgReceive)
					return
				}
			}
			delArr := []string{"-d", "~d", "delete"}
			for _, v := range delArr {
				if strings.HasPrefix(content, v) {
					// 删除
					msgReceive.Text.Content = strings.TrimSpace(strings.TrimPrefix(content, v))
					handleReturnMsg(w, delMessage(msgReceive), msgReceive)
					return
				}
			}
			queryDetailArr := []string{"-qd", "~qd"}
			for _, v := range queryDetailArr {
				if strings.HasPrefix(content, v) {
					msgReceive.Text.Content = strings.TrimSpace(strings.TrimPrefix(content, v))
					handleReturnMsg(w, queryDetailMessage(msgReceive), msgReceive)
					return
				}
			}
			queryArr := []string{"-q", "~q", "-query", "~query"}
			for _, v := range queryArr {
				if strings.HasPrefix(content, v) {
					msgReceive.Text.Content = strings.TrimSpace(strings.TrimPrefix(content, v))
					handleReturnMsg(w, queryMessage(msgReceive), msgReceive)
					return
				}
			}
			msgReceive.Text.Content = strings.TrimSpace(content)
			handleReturnMsg(w, queryMessage(msgReceive), msgReceive)
		} else {
			handleReturnMsg(w, "只支持Text！", msgReceive)
		}
	} else {
		fmt.Fprintf(w, "暂不支持群聊哦！")
	}
}

func handleReturnMsg(w http.ResponseWriter, content string, receive *MsgReceive) {
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

func queryDetailMessage(receive *MsgReceive) string {
	myDb := db.GetDb()
	rows, _ := myDb.Query("SELECT ID, CONTENT FROM MESSAGE WHERE CONTENT LIKE $1 AND FROM_USER_ID = $2", "%"+receive.Text.Content+"%", receive.SenderStaffId)
	responseContent := ""
	haveContent := false
	for rows.Next() {
		if responseContent != "" {
			responseContent += "\r\n"
		}
		haveContent = true
		var id int64
		var content string
		rows.Scan(&id, &content)
		responseContent += fmt.Sprintf("id: %v, content: %v", id, content)
	}
	if !haveContent {
		return "没有查到任何内容哦~"
	} else {
		return responseContent
	}
}

func queryMessage(receive *MsgReceive) string {
	myDb := db.GetDb()
	rows, _ := myDb.Query("SELECT CONTENT FROM MESSAGE WHERE CONTENT LIKE $1 AND FROM_USER_ID = $2", "%"+receive.Text.Content+"%", receive.SenderStaffId)
	responseContent := ""
	haveContent := false
	for rows.Next() {
		if responseContent != "" {
			responseContent += "\r\n"
		}
		haveContent = true
		var content string
		rows.Scan(&content)
		responseContent += content
	}
	if !haveContent {
		return "没有查到任何内容哦~"
	} else {
		return responseContent
	}
}

func delMessage(receive *MsgReceive) string {
	db.InsertOrUpdateSqlByStmt("DELETE FROM MESSAGE WHERE ID = $1", receive.Text.Content)
	return "删除成功~"
}

func saveMessage(msg *MsgReceive) string {
	db.InsertOrUpdateSqlByStmt("INSERT INTO MESSAGE (CONTENT, THIRD_MESSAGE_ID, FROM_USER_ID) VALUES ($1, $2, $3)", msg.Text.Content, msg.MsgId, msg.SenderStaffId)
	return "数据已保存~"
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
