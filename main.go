package main

import (
	"goDingRobot/db"
	"goDingRobot/ding"
	"goDingRobot/wx"
	"log"
	"net/http"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	log.Println(uri)
	if strings.HasPrefix(uri, ding.ServeDingRobotPrefix) {
		ding.HandleDingRobot(w, r)
	} else if strings.HasPrefix(uri, wx.ServerWxPrefix) {
		wx.HandleWxRobot(w, r)
	}
}

func main() {
	myGlobalConfig := handleConfig()
	url := myGlobalConfig.DataSource.Url
	ding.NewDingRobot(myGlobalConfig.Ding.AppKey, myGlobalConfig.Ding.AppSecret)
	wx.NewWx(myGlobalConfig.Wx.AppId, myGlobalConfig.Wx.AppSecret, myGlobalConfig.Wx.Token)

	db.InitSqlite(url)
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8000", nil)
}
