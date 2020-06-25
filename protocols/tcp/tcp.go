package tcp

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/beanstalkd/go-beanstalk"
	"go-central-control/utils"
	"net"
	"time"
)

var data = map[string]interface{}{
	"channel": 0,
	"body": map[string]interface{}{

	},
}

func Out(worker map[string]interface{}, mq func(msg string)) {
	host := worker["host"]
	port := worker["port"]
	protocol := worker["protocol"]
	channel := worker["channel"]
	// simple tcp server
	// 1.listen ip+port

	listener, err := net.Listen(protocol.(string), host.(string)+":"+port.(string))
	if err != nil {
		utils.Log.Error(err)
		return
	}
	defer listener.Close()
	utils.Log.Infof("Listening:%s//%s:%s", protocol.(string), host.(string), port.(string))

	// 2.accept client request
	// 3.create goroutine for each request
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.Log.Error(err)
			break
		}

		defer conn.Close()
		go func() {
			addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
			b, err := beanstalk.Dial("tcp", addr)
			tubeSet := beanstalk.NewTubeSet(b, channel.(string))
			if err != nil {
				utils.Log.Error(err)
				return
			}
			timeout, err := beego.AppConfig.Int("beanstalkdreservetimeout")

			for {
				id, body, err := tubeSet.Reserve(time.Duration(int(time.Duration(timeout) * time.Second)))
				str := ""
				if err != nil {
					str = utils.Msg(500, "Heart")
				} else {
					str = string(body)
				}

				b.Delete(id)
				_, err = conn.Write([]byte(str))
				if err != nil {
					utils.Log.Error(err)
					conn.Close()
					break
				}
			}
		}()

	}
	return
}

func Sender(worker map[string]interface{}) {
	host := worker["host"]
	port := worker["port"]
	protocol := worker["protocol"]
	channel := worker["channel"]
	// simple tcp server
	// 1.listen ip+port

	listener, err := net.Listen(protocol.(string), host.(string)+":"+port.(string))
	if err != nil {
		utils.Log.Error(err)
		return
	}
	defer listener.Close()
	utils.Log.Infof("Listening:%s//%s:%s", protocol.(string), host.(string), port.(string))

	// 2.accept client request
	// 3.create goroutine for each request
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.Log.Error(err)
			break
		}

		defer conn.Close()
		go func() {
			addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
			b, err := beanstalk.Dial("tcp", addr)
			tubeSet := beanstalk.NewTubeSet(b, channel.(string))
			if err != nil {
				utils.Log.Error(err)
				return
			}
			timeout, err := beego.AppConfig.Int("beanstalkdreservetimeout")

			for {
				id, body, err := tubeSet.Reserve(time.Duration(int(time.Duration(timeout) * time.Second)))
				str := ""
				if err != nil {
					str = utils.Msg(500, "Heart")
				} else {
					str = string(body)
				}

				b.Delete(id)
				_, err = conn.Write([]byte(str))
				if err != nil {
					utils.Log.Error(err)
					conn.Close()
					break
				}
			}
		}()

	}
	return
}

func Receiver(worker map[string]interface{}) {
	host := worker["host"]
	port := worker["port"]
	protocol := worker["protocol"]
	// simple tcp server
	// 1.listen ip+port
	listener, err := net.Listen(protocol.(string), host.(string)+":"+port.(string))
	if err != nil {
		utils.Log.Error(err)
		return
	}
	defer listener.Close()
	utils.Log.Infof("Listening:%s//%s:%s", protocol.(string), host.(string), port.(string))

	addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
	bConn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		utils.Log.Error(err)
		return
	}
	// 2.accept client request
	// 3.create goroutine for each request
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.Log.Error(err)
			continue
		}

		defer conn.Close()
		go func() {
			for {
				var buf [128]byte
				n, err := conn.Read(buf[:])
				if err != nil {
					utils.Log.Error(err)
					conn.Close()
					break
				}

				recvStr := string(buf[:n])
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
					break
				}

				str := utils.Msg(200, "Success")
				_, err = conn.Write([]byte(str))
				if err != nil {
					utils.Log.Error(err)
					conn.Close()
					continue
				}
			}
		}()
	}
}
