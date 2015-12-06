package routers

import (
	"github.com/astaxie/beego"
	"wifisensor/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/1.php", &controllers.MainController{}, "post:AddSensorData")
}
