package main

import (
	"fmt"
	"goDingRobot/db"
	"goDingRobot/ding"
	"goDingRobot/wx"
	"net/http"
	"strings"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	fmt.Println(uri)
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

	db.InitSqlite(url)
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8000", nil)
}
