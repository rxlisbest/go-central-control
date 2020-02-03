package main

import (
	"central-control/protocols/tcp"
	"central-control/utils"
	"github.com/astaxie/beego/config"
	"time"
)

func main() {
	utils.InitLogging()

	workerConfig, err := config.NewConfig("json", "./conf/worker.json")
	if err != nil {
		utils.Log.Error(err)
	}
	workers, err := workerConfig.DIY("workers")
	if err != nil {
		utils.Log.Error(err)
	}
	for _, v := range workers.([]interface{}) {
		//create goroutine for each connect
		switch v.(map[string]interface{})["protocol"].(string) {
		case "tcp":
			go tcp.Input(v.(map[string]interface{}))
		default:
			utils.Log.Error("Protocol not supported")
		}
	}
	for {
		time.Sleep(time.Second)
	}
}
