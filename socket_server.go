package main

import (
	"central-control/utils"
	"central-control/protocols"
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
		go protocols.Tcp(v.(map[string]interface{}))
	}
	for {
		time.Sleep(time.Second)
	}
}
