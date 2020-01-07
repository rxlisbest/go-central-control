package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/beanstalkd/go-beanstalk"
	"strconv"
	"time"

	//"github.com/beanstalkd/go-beanstalk"
	//"time"
)

type SignalController struct {
	beego.Controller
}

var data = map[string]interface{}{
	"channel": 0,
	"body": map[string]interface{}{

	},
}

var reponse = map[string]interface{}{
	"code":    0,
	"message": "",
}

func (c *SignalController) Post() {
	request := c.Ctx.Input.RequestBody
	err := json.Unmarshal(request, &data)
	fmt.Println(data)
	if err != nil {
		fmt.Println("json.Unmarshal is err:", err.Error())
		return
	}
	b, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		fmt.Println("json.Unmarshal is err:", err.Error())
	}
	channel, ok := data ["channel"]
	if ok {
		b.Tube.Name = strconv.FormatInt(int64(channel.(float64)), 10)
	}

	_, err = b.Put([]byte(request), 2, 0, 120*time.Second)
	if err != nil {
		fmt.Println("json.Unmarshal is err:", err.Error())
	}
	//_, _, _ = b.Reserve(5 * time.Second)
	c.Data["json"] = data // json对象
	c.ServeJSON()
	return

	//c.Data["Website"] = string(body[:])
	//c.Data["Email"] = "astaxie@gmail.com"
	//c.TplName = "index.tpl"
}
