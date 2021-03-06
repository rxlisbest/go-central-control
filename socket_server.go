package main

import (
	"go-central-control/protocols/tcp"
	"go-central-control/protocols/websocket"
	"go-central-control/utils"
	"github.com/astaxie/beego/config"
	"time"
)

func main() {
	utils.InitLogging()

	workerConfig, err := config.NewConfig("json", "./conf/worker.json")
	if err != nil {
		utils.Log.Error(err)
	}
	receivers, err := workerConfig.DIY("receivers")
	if err != nil {
		utils.Log.Error(err)
	}
	for _, v := range receivers.([]interface{}) {
		//create goroutine for each connect
		switch v.(map[string]interface{})["protocol"].(string) {
		case "tcp":
			go tcp.Receiver(v.(map[string]interface{}))
		case "ws":
			go websocket.Receiver(v.(map[string]interface{}))
		default:
			utils.Log.Error("Protocol not supported")
		}
	}
	senders, err := workerConfig.DIY("senders")
	if err != nil {
		utils.Log.Error(err)
	}
	for _, v := range senders.([]interface{}) {
		//create goroutine for each connect
		switch v.(map[string]interface{})["protocol"].(string) {
		case "tcp":
			go tcp.Sender(v.(map[string]interface{}))
		case "ws":
			go websocket.Sender(v.(map[string]interface{}))
		default:
			utils.Log.Error("Protocol not supported")
		}
	}
	for {
		time.Sleep(time.Second)
	}
}
