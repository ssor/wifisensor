package main

import (
	"github.com/astaxie/beego"
	_ "wifisensor/routers"
)

func main() {
	beego.Run()
}
