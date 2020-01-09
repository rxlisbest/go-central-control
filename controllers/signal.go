package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"github.com/beanstalkd/go-beanstalk"
	"time"
	//"fmt"
)

type SignalController struct {
	beego.Controller
}

var data = map[string]interface{}{
	"channel": 0,
	"body": map[string]interface{}{

	},
}

var response = map[string]interface{}{
	"code":    "0",
	"message": "",
}

func (s *SignalController) Error(code int, message string) {
	s.Ctx.Output.Status = code
	response["code"] = code
	response["message"] = message
	s.Data["json"] = response // json对象
	s.ServeJSON()
}

func (s *SignalController) Post() {
	request := s.Ctx.Input.RequestBody
	err := json.Unmarshal(request, &data)
	if err != nil {
		s.Error(400, err.Error())
		return
	}
	channel := data["channel"]
	if channel != "0" {
		switch channel.(type) {
		case string:
		default:
			s.Error(400, "The format of the field `channel` is incorrect")
			return
		}
	} else {
		s.Error(400, "The field `channel` is required")
		return
	}
	addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		s.Error(400, err.Error())
		return
	}

	tube := &beanstalk.Tube{Conn: conn, Name: channel.(string)}
	_, err = tube.Put([]byte(request), 1, 0, 120*time.Second)
	if err != nil {
		s.Error(500, err.Error())
		return
	}
	s.Data["json"] = data // json对象
	s.ServeJSON()
	return
}

func (s *SignalController) Get() {
	channel := s.GetString("channel")
	if channel == "" {
		s.Error(400, "The params `channel` is required")
		return
	}
	init := s.GetString("init")
	addr := beego.AppConfig.String("beanstalkdaddr") + ":" + beego.AppConfig.String("beanstalkdport")
	conn, err := beanstalk.Dial("tcp", addr)
	if err != nil {
		s.Error(500, err.Error())
		return
	}

	bm, err := cache.NewCache("file", `{}`)
	if err != nil {
		s.Error(500, err.Error())
		return
	}

	clientId := time.Now().Nanosecond()
	bm.Put("client_id", clientId, 0*time.Second)

	tubeSet := beanstalk.NewTubeSet(conn, channel)

	// clear queue
	if init == "1" {
		for {
			id, _, err := tubeSet.Reserve(0 * time.Second)
			if err != nil {
				break
			}
			conn.Delete(id)
		}
	}

	timeout, err := beego.AppConfig.Int("beanstalkdreservetimeout")
	if err != nil {
		s.Error(500, err.Error())
		return
	}

	id, body, err := tubeSet.Reserve(time.Duration(int(time.Duration(timeout) * time.Second)))
	if err != nil {
		s.Error(408, err.Error())
		return
	}

	currentClientId := bm.Get("client_id")
	if (currentClientId != clientId) {
		s.Error(409, "Request expired")
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		s.Error(500, err.Error())
		return
	}
	conn.Delete(id)
	s.Data["json"] = data // json对象
	s.ServeJSON()
	return
	//for {
	//	select {
	//	default:
	//		id, body, _ := tubeSet.Reserve(5 * time.Second)
	//		conn.Delete(id)
	//		fmt.Print(1)
	//		fmt.Print(string(body[:]))
	//	}
	//}
}
