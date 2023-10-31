package hugeldate

import (
	"log"
	"time"
)

var CentralTimeZone *time.Location

func init() {
	var err error
	CentralTimeZone, err = time.LoadLocation("US/Central")
	if err != nil {
		log.Printf("error loading central timezone: %v", err)
		CentralTimeZone = time.Local
	}
}

// 3 day window. Nov 10, 11, and 12
var StartHugel = time.Date(2023, 11, 10, 0, 0, 0, 0, CentralTimeZone)
var EndHugel = StartHugel.Add(time.Hour * 24 * 3)
