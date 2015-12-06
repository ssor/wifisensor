package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SensorRecordList map[string]*SensorRecord

func (s SensorRecordList) Add(r *SensorRecord) {
	id := r.GetShortMAC()
	if _, ok := s[id]; !ok {
		s[id] = r
	}
}

type SensorRecord struct {
	Time        string
	UnixSeconds int64
	MAC         string
	SSID        string
	Signal      int
}

func (s *SensorRecord) String() string {
	return fmt.Sprintf("Time: %s  MAC: %s Signal: %d  SSID: %s", s.Time, s.MAC, s.Signal, s.SSID)
}
func (s *SensorRecord) GetShortMAC() string {
	return FormatMac(s.MAC)
}
func (s *SensorRecord) Equal(r *SensorRecord) bool {
	return s.MAC == r.MAC && s.Signal == r.Signal
}
func NewSensorRecord(mac, ssid, signal string, now time.Time) *SensorRecord {
	iSignal, err := strconv.Atoi(signal)
	if err != nil {
		return nil
	}
	return &SensorRecord{
		MAC:         mac,
		SSID:        ssid,
		Signal:      iSignal,
		Time:        now.Format(time.RFC3339),
		UnixSeconds: now.Unix(),
	}
}

// B4:0B:44:01:03:3E => MACB40B4401033E
func FormatMac(mac string) string {
	return "MAC" + strings.Replace(mac, ":", "", -1)
}
func ParseRawDataToRecords(data string) SensorRecordList {
	// data = "{{00:08:22:B8:D3:FB,,-89},{B4:0B:44:01:03:3E,,-38},{B4:0B:44:01:03:3E,,-38},{B4:0B:44:01:03:3E,,-37},{B4:0B:44:01:03:3E,,-38},{B4:0B:44:01:03:3E,,-38},{B4:0B:44:01:03:3E,,-38},{B4:0B:44:01:03:3E,,-36},{B4:0B:44:01:03:3E,,-37},{B4:0B:44:01:03:3E,,-37},{B4:0B:44:01:03:3E,,-37},{B4:0B:44:01:03:3E,,-46},{B4:0B:44:01:03:3E,,-58},{FC:64:BA:53:53:66,nadu,-88},{FC:64:BA:53:53:66,TP-LINK_NXW,-89},{FC:64:BA:53:53:66,PC-201309071840_Network,-89},{FC:64:BA:53:53:66,菊花朵朵开,-90}}"
	re := regexp.MustCompile(`\{[-0-9,:A-Z\w\p{Han}]+\}`)
	recordsRaw := re.FindAllString(data, -1)
	if len(recordsRaw) <= 0 {
		return nil
	}
	l := make(SensorRecordList)
	now := time.Now()
	for _, recordRaw := range recordsRaw {
		// beego.Trace(recordRaw)
		items := strings.Split(recordRaw[1:len(recordRaw)-1], ",")
		if len(items) < 2 {
			beego.Trace(data)
			beego.Trace(items)
			continue
		}
		record := NewSensorRecord(items[0], items[1], items[2], now)
		// beego.Trace(record)
		l.Add(record)
	}
	return l
}
