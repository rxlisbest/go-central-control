package tcp

import (
	"central-control/utils"
	"github.com/astaxie/beego"
	"github.com/beanstalkd/go-beanstalk"
	"net"
	"time"
)

func Output(worker map[string]interface{}) {
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
		addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
		b, err := beanstalk.Dial("tcp", addr)
		tubeSet := beanstalk.NewTubeSet(b, channel.(string))
		if err != nil {
			utils.Log.Error(err)
			break
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
	}
}

func Input(worker map[string]interface{}) {
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

	addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		utils.Log.Error(err)
		return
	}
	tube := &beanstalk.Tube{Conn: conn, Name: channel.(string)}
	// 2.accept client request
	// 3.create goroutine for each request
	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.Log.Error(err)
			break
		}

		defer conn.Close()

		_, err = tube.Put([]byte(request), 1, 0, 120*time.Second)

		addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
		b, err := beanstalk.Dial("tcp", addr)
		tubeSet := beanstalk.NewTubeSet(b, channel.(string))
		if err != nil {
			utils.Log.Error(err)
			break
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
	}
}