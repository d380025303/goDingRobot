package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"goDingRobot/db"
	"io/ioutil"
	"log"
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
	Url string `json:"url"`
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
	url := MyGlobalConfig.DataSource.Url
	db.InitSqlite(url)
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8000", nil)
}

func handleConfig() {
	log.Println("resolve config...")
	var applicationLocation string
	flag.StringVar(&applicationLocation, "c", "", "位置")

	var url string
	var appKey string
	flag.StringVar(&appKey, "appKey", "", "钉钉Key")
	var appSecret string
	flag.StringVar(&appSecret, "appSecret", "", "钉钉Secret")
	flag.Parse()

	if appKey != "" && appSecret != "" {
		url = "/usr/data/ding.db"
		MyGlobalConfig = &GlobalConfig{
			Ding: Ding{
				AppKey:    appKey,
				AppSecret: appSecret,
			},
			DataSource: DataSource{
				Url: url,
			},
		}
	} else {
		MyGlobalConfig = &GlobalConfig{}

		var file *os.File
		defer file.Close()
		if applicationLocation != "" {
			file, _ = os.Open(applicationLocation)
		}
		if file == nil {
			file, _ = os.Open("application.json")
		}
		if file == nil {
			panic("can not read config file named application.json")
		}
		bytes, err2 := ioutil.ReadAll(file)
		if err2 != nil {
			panic("can not read config file named application.json")
		}
		err2 = json.Unmarshal(bytes, MyGlobalConfig)
		if err2 != nil {
			panic("can not read config file named application.json")
		}
	}
	log.Println(fmt.Sprintf(`appKey: %v
appSecret: %v
url: %v`, MyGlobalConfig.Ding.AppKey, MyGlobalConfig.Ding.AppSecret, MyGlobalConfig.DataSource.Url))
}
