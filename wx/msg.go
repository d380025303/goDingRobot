package wx

import "encoding/xml"

type ReceiveTextMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        string   `xml:"MsgId"`
	MsgDataId    string   `xml:"MsgDataId"`
	Idx          string   `xml:"Idx"`
}
