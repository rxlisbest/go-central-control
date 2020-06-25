package websocket

import (
	"go-central-control/utils"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/beanstalkd/go-beanstalk"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var data = map[string]interface{}{
	"channel": 0,
	"body": map[string]interface{}{

	},
}

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
		for {
			id, body, err := tubeSet.Reserve(time.Duration(int(time.Duration(timeout) * time.Second)))
			str := ""
			if err != nil {
				str = utils.Msg(500, "Heart")
			} else {
				str = string(body)
			}

			b.Delete(id)

			err = conn.WriteMessage(1, []byte(str))
			if err != nil {
				utils.Log.Error(err)
				conn.Close()
				continue
			}
		}
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

	addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
	bConn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		utils.Log.Error(err)
		return
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			utils.Log.Error(err)
			return
		}
		defer conn.Close()
		for {
			mt, recvStr, err := conn.ReadMessage()
			if err != nil {
				utils.Log.Error(err)
				conn.Close()
				continue
			}
			err = json.Unmarshal([]byte(recvStr), &data)
			if err != nil {
				utils.Log.Error(err)
				continue
			}

			channel := data["channel"]
			if channel != "0" {
				switch channel.(type) {
				case string:
				default:
					utils.Log.Error(400, "The format of the field `channel` is incorrect")
					continue
				}
			} else {
				utils.Log.Error(400, "The field `channel` is required")
				continue
			}

			tube := &beanstalk.Tube{Conn: bConn, Name: channel.(string)}

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
	})
	utils.Log.Infof("Listening:%s//%s:%s", protocol.(string), host.(string), port.(string))
	err = http.ListenAndServe(host.(string)+":"+port.(string), mux)
	if err != nil {
		utils.Log.Error(err)
		return
	}
}
