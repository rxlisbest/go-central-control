package main

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/op/go-logging"
	"net"
	"os"
	"time"
)

var response = map[string]interface{}{
	"code":    "0",
	"message": "",
}

func Msg(code int, message string) string {
	response["code"] = code
	response["message"] = message
	responseJson, _ := json.Marshal(response)
	return string(responseJson)
}

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}%{color:reset}`,
)

func initLogging() {
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(-1, "")
	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)
}

func main() {
	initLogging()

	workerConfig, err := config.NewConfig("json", "./conf/worker.json")
	if err != nil {
		log.Error(err)
	}
	workers, err := workerConfig.DIY("workers")
	if err != nil {
		log.Error(err)
	}
	for _, v := range workers.([]interface{}) {
		//create goroutine for each connect
		go process(v.(map[string]interface{}))
	}
	for {
		time.Sleep(1 * time.Second)
	}
}

func process(worker map[string]interface{}) {
	host := worker["host"]
	port := worker["port"]
	protocol := worker["protocol"]
	channel := worker["channel"]
	// simple tcp server
	// 1.listen ip+port

	listener, err := net.Listen(protocol.(string), host.(string)+":"+port.(string))
	if err != nil {
		log.Error(err)
		return
	}
	defer listener.Close()
	log.Infof("Listening:%s//%s:%s", protocol.(string), host.(string), port.(string))

	// 2.accept client request
	// 3.create goroutine for each request
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error(err)
			break
		}

		defer conn.Close()
		addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
		b, err := beanstalk.Dial("tcp", addr)
		tubeSet := beanstalk.NewTubeSet(b, channel.(string))
		if err != nil {
			log.Error(err)
			break
		}
		timeout, err := beego.AppConfig.Int("beanstalkdreservetimeout")

		for {
			id, body, err := tubeSet.Reserve(time.Duration(int(time.Duration(timeout) * time.Second)))
			str := ""
			if err != nil {
				str = Msg(500, "Heart")
			} else {
				str = string(body)
			}

			b.Delete(id)
			_, err = conn.Write([]byte(str))
			if err != nil {
				log.Error(err)
				conn.Close()
				break
			}
		}
	}
}
