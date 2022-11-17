package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type GlobalConfig struct {
	Ding       Ding       `json:"ding"`
	DataSource DataSource `json:"datasource"`
	Wx         Wx         `json:"wx"`
}

type Wx struct {
	AppId     string `json:"appId"`
	AppSecret string `json:"appSecret"`
}

type Ding struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
	AgentId   string `json:"agentId"`
}

type DataSource struct {
	Url string `json:"url"`
}

func handleConfig() *GlobalConfig {
	var MyGlobalConfig *GlobalConfig

	log.Println("resolve config...")
	var applicationLocation string
	flag.StringVar(&applicationLocation, "c", "", "位置")

	var dockerInd string
	flag.StringVar(&dockerInd, "dockerInd", "N", "是否为docker启动")

	var url string
	var appKey string
	flag.StringVar(&appKey, "appKey", "", "钉钉Key")
	var appSecret string
	flag.StringVar(&appSecret, "appSecret", "", "钉钉Secret")

	var wxAppId string
	flag.StringVar(&wxAppId, "wxAppId", "", "公众号id")
	var wxAppSecret string
	flag.StringVar(&wxAppSecret, "wxAppSecret", "", "公众号Secret")

	flag.Parse()

	if "Y" == dockerInd {
		url = "/usr/data/ding.db"
		MyGlobalConfig = &GlobalConfig{
			Ding: Ding{
				AppKey:    appKey,
				AppSecret: appSecret,
			},
			DataSource: DataSource{
				Url: url,
			},
			Wx: Wx{
				AppId:     wxAppId,
				AppSecret: wxAppSecret,
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
wxAppId: %v
wxAppSecret: %v
url: %v
`, MyGlobalConfig.Ding.AppKey, MyGlobalConfig.Ding.AppSecret, MyGlobalConfig.Wx.AppId, MyGlobalConfig.Wx.AppSecret, MyGlobalConfig.DataSource.Url))
	return MyGlobalConfig
}
