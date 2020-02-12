package test

import (
	"central-control/utils"
	"github.com/astaxie/beego"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	// 跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Sender(worker map[string]interface{}) {
	host := worker["host"]
	port := worker["port"]
	protocol := worker["protocol"]
	channel := worker["channel"]

	addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
	b, err := beanstalk.Dial("tcp", addr)
	tubeSet := beanstalk.NewTubeSet(b, channel.(string))
	if err != nil {
		utils.Log.Error(err)
		return
	}
	timeout, err := beego.AppConfig.Int("beanstalkdreservetimeout")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			utils.Log.Error(err)
			return
		}
		defer conn.Close()
		go func() {
			for {
				mt, _, err := conn.ReadMessage()
				if err != nil {
					utils.Log.Error(err)
					conn.Close()
					continue
				}

				id, body, err := tubeSet.Reserve(time.Duration(int(time.Duration(timeout) * time.Second)))
				str := ""
				if err != nil {
					str = utils.Msg(500, "Heart")
				} else {
					str = string(body)
				}

				b.Delete(id)

				err = conn.WriteMessage(mt, []byte(str))
				if err != nil {
					utils.Log.Error(err)
					conn.Close()
					continue
				}
			}
		}()
	})
	utils.Log.Infof("Listening:%s//%s:%s", protocol.(string), host.(string), port.(string))
	err = http.ListenAndServe(host.(string)+":"+port.(string), mux)
	if err != nil {
		utils.Log.Error(err)
		return
	}
}

func Receiver(worker map[string]interface{}) {
	host := worker["host"]
	port := worker["port"]
	protocol := worker["protocol"]
	channel := worker["channel"]

	addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		utils.Log.Error(err)
		return
	}
	tube := &beanstalk.Tube{Conn: conn, Name: channel.(string)}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			utils.Log.Error(err)
			return
		}
		defer conn.Close()
		go func() {
			for {
				mt, recvStr, err := conn.ReadMessage()
				if err != nil {
					utils.Log.Error(err)
					conn.Close()
					continue
				}
				_, err = tube.Put([]byte(recvStr), 1, 0, 120*time.Second)
				if err != nil {
					utils.Log.Error(err)
					conn.Close()
					continue
				}

				str := utils.Msg(200, "Success")
				err = conn.WriteMessage(mt, []byte(str))
				if err != nil {
					utils.Log.Error(err)
					conn.Close()
					continue
				}
			}
		}()
	})
	utils.Log.Infof("Listening:%s//%s:%s", protocol.(string), host.(string), port.(string))
	err = http.ListenAndServe(host.(string)+":"+port.(string), nil)
	if err != nil {
		utils.Log.Error(err)
		return
	}
}
