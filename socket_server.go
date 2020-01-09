package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/config"
	_ "github.com/astaxie/beego/config"
	"github.com/beanstalkd/go-beanstalk"
	"net"
	"time"
)

var response = map[string]interface{}{
	"code":    "0",
	"message": "",
}

func Error(code int, message string) string {
	response["code"] = code
	response["message"] = message
	responseJson, _ := json.Marshal(response)
	return string(responseJson)
}

func main() {
	workerConfig, err := config.NewConfig("json", "./conf/worker.json")
	if err != nil {
		fmt.Printf(Error(500, err.Error()))
	}
	workers, err := workerConfig.DIY("workers")
	if err != nil {
		fmt.Printf(Error(500, err.Error()))
	}
	for _, v := range workers.([]interface{}) {
		host := v.(map[string]interface{})["host"]
		port := v.(map[string]interface{})["port"]
		protocol := v.(map[string]interface{})["protocol"]
		channel := v.(map[string]interface{})["channel"]
		// simple tcp server
		// 1.listen ip+port
		listener, err := net.Listen(protocol.(string), host.(string)+":"+port.(string))
		if err != nil {
			fmt.Printf(Error(500, err.Error()))
			return
		}

		// 2.accept client request
		// 3.create goroutine for each request

		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf(Error(500, err.Error()))
		}

		//create goroutine for each connect
		go process(conn, channel.(string))
	}
	time.Sleep(10 * time.Second)
}

func process(conn net.Conn, channel string) {
	defer conn.Close()
	addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
	b, err := beanstalk.Dial("tcp", addr)
	tubeSet := beanstalk.NewTubeSet(b, channel)
	if err != nil {
		return
	}
	timeout, err := beego.AppConfig.Int("beanstalkdreservetimeout")

	fmt.Printf("success")
	for {
		id, body, err := tubeSet.Reserve(time.Duration(int(time.Duration(timeout) * time.Second)))
		if err != nil {
		}

		b.Delete(id)
		//var buf [128]byte
		//n, err := conn.Read(buf[:])
		//if err != nil {
		//	fmt.Printf("read from connect failed, err: %v\n", err)
		//	break
		//}
		str := body
		//fmt.Printf("receive from client, data: %v\n", str)
		conn.Write(str)
	}
}
