package controllers

import (
	"github.com/astaxie/beego"
	// "regexp"
	"fmt"
	"gopkg.in/mgo.v2"
)

var (
	MONGO_SESSION   *mgo.Session
	DATA_COLLECTION *mgo.Collection
	SENSOR_DB       *mgo.Database
)

func init() {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(fmt.Sprintf("connect to mongodb error: %s", err))
	} else {
		beego.Info("connect to mongodb OK")
	}
	MONGO_SESSION = session
	SENSOR_DB = MONGO_SESSION.DB("wifisensor")

}

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplNames = "index.tpl"
}
func (m *MainController) AddSensorData() {
	gateway := m.GetString("gateway")
	recordsRaw := m.GetString("records")
	// bodyRaw := string(m.Ctx.Input.RequestBody)
	//gateway=00:0A:C2:E9:FD:8C&records={{44:D4:E0:D3:BA:5F,,-84},{44:D4:E0:D3:BA:5F,,-83}}
	if len(gateway) <= 0 {
		beego.Trace("NO body")
	} else {
		beego.Trace(fmt.Sprintf("gateway: %s", gateway))
		beego.Trace(fmt.Sprintf("records: %s", recordsRaw))
		// beego.Trace((bodyRaw))

	}
	records := ParseRawDataToRecords(recordsRaw)
	for _, record := range records {
		SENSOR_DB.C(FormatMac(gateway)).Insert(record)
		beego.Trace(record.String())
	}

	m.ServeJson()
}
