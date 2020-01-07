package routers

import (
	"central-control/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/signals", &controllers.SignalController{})
}
