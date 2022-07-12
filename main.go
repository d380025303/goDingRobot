package main

import (
	"encoding/json"
	"fmt"
	"goDingRobot/db"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type GlobalConfig struct {
	Ding       Ding       `json:"ding"`
	DataSource DataSource `json:"datasource"`
}

type Ding struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
	AgentId   string `json:"agentId"`
}

type DataSource struct {
	url string `json:"url"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	fmt.Println(uri)
	if strings.HasPrefix(uri, dingRobotPrefix) {
		handleDingRobot(w, r)
	}
}

var MyGlobalConfig *GlobalConfig

func main() {
	handleConfig()
	url := MyGlobalConfig.DataSource.url
	db.InitSqlite(url)
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8000", nil)
}

func handleConfig() {
	MyGlobalConfig = &GlobalConfig{}
	file, _ := os.Open("application.json")
	defer file.Close()
	bytes, err2 := ioutil.ReadAll(file)
	if err2 != nil {
		panic("can not read config file named application.json")
	}
	err2 = json.Unmarshal(bytes, MyGlobalConfig)
	if err2 != nil {
		panic("can not read config file named application.json")
	}
}
