package robot

import (
	"fmt"
	"goDingRobot/db"
	"strings"
)

var (
	ACTION_SAVE         = "SAVE"
	ACTION_DEL          = "DELETE"
	ACTION_QUERY_DETAIL = "QUERY_DETAIL"
	ACTION_QUERY        = "ACTION_QUERY"
)

type RobotRespVo struct {
	Action string
	Msg    string
}

func saveMessage(content, msgId, fromUserId string) string {
	db.InsertOrUpdateSqlByStmt("INSERT INTO MESSAGE (CONTENT, THIRD_MESSAGE_ID, FROM_USER_ID) VALUES ($1, $2, $3)", content, msgId, fromUserId)
	return "数据已保存~"
}

func delMessage(content string) string {
	db.InsertOrUpdateSqlByStmt("DELETE FROM MESSAGE WHERE ID = $1", content)
	return "删除成功~"
}

func queryMessage(content string) string {
	myDb := db.GetDb()
	rows, _ := myDb.Query("SELECT CONTENT FROM MESSAGE WHERE CONTENT LIKE $1", "%"+content+"%")
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

func queryDetailMessage(content string) string {
	myDb := db.GetDb()
	rows, _ := myDb.Query("SELECT ID, CONTENT FROM MESSAGE WHERE CONTENT LIKE $1", "%"+content+"%")
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

func HandleContent(content, msgId, fromUserId string) *RobotRespVo {
	saveArr := []string{"-s", "~s", "-save", "~save", "s", "save"}
	for _, v := range saveArr {
		if strings.HasPrefix(content, v) {
			// 保存
			content = strings.TrimSpace(strings.TrimPrefix(content, v))
			respMsg := saveMessage(content, msgId, fromUserId)
			return &RobotRespVo{
				Action: ACTION_DEL,
				Msg:    respMsg,
			}
		}
	}
	delArr := []string{"-d", "~d", "delete"}
	for _, v := range delArr {
		if strings.HasPrefix(content, v) {
			// 删除
			content = strings.TrimSpace(strings.TrimPrefix(content, v))
			respMsg := delMessage(content)
			return &RobotRespVo{
				Action: ACTION_SAVE,
				Msg:    respMsg,
			}
		}
	}
	queryDetailArr := []string{"-qd", "~qd"}
	for _, v := range queryDetailArr {
		if strings.HasPrefix(content, v) {
			content = strings.TrimSpace(strings.TrimPrefix(content, v))
			respMsg := queryDetailMessage(content)
			return &RobotRespVo{
				Action: ACTION_QUERY_DETAIL,
				Msg:    respMsg,
			}
		}
	}
	queryArr := []string{"-q", "~q", "-query", "~query"}
	for _, v := range queryArr {
		if strings.HasPrefix(content, v) {
			content = strings.TrimSpace(strings.TrimPrefix(content, v))
			respMsg := queryMessage(content)
			return &RobotRespVo{
				Action: ACTION_QUERY,
				Msg:    respMsg,
			}
		}
	}
	content = strings.TrimSpace(content)
	respMsg := queryMessage(content)
	return &RobotRespVo{
		Action: ACTION_QUERY,
		Msg:    respMsg,
	}
}
