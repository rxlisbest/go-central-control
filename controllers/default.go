package controllers

import (
	"github.com/astaxie/beego"
	"github.com/beanstalkd/go-beanstalk"
	"time"
)

type MainController struct {
	beego.Controller
}

var JsonOutput = map[string]interface{}{
	"code" : 0,
	"msg" : "",
	"data" : map[string]interface{}{

	},
}


func (c *MainController) Get() {
	//b, _ := beanstalk.Dial("tcp", "127.0.0.1:11300")
	//_, _ = b.Put([]byte("hello"), 1, 0, 120*time.Second)
	//_, _, _ = b.Reserve(5 * time.Second)
	//c.Data["json"] = JsonOutput                     // json对象
	//c.ServeJSON()
	//return

	c.Data["Website"] = "rxl.net"
	c.Data["Email"] = "astaxie@gmail.com"               // json对象
	c.ServeJSON()
	return
}

func (c *MainController) Post() {
	b, _ := beanstalk.Dial("tcp", "127.0.0.1:11300")
	_, _ = b.Put([]byte("hello"), 1, 0, 120*time.Second)
	_, _, _ = b.Reserve(5 * time.Second)
	c.Data["json"] = JsonOutput                     // json对象
	c.ServeJSON()
	return

	//c.Data["Website"] = string(body[:])
	//c.Data["Email"] = "astaxie@gmail.com"
	//c.TplName = "index.tpl"
}
