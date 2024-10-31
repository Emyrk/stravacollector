package hugeldate

import (
	"log"
	"time"
)

var CentralTimeZone *time.Location
var Year2023 Dates
var Year2024 Dates

type Dates struct {
	Start time.Time
	End   time.Time
}

func init() {
	var err error
	CentralTimeZone, err = time.LoadLocation("US/Central")
	if err != nil {
		log.Printf("error loading central timezone: %v", err)
		CentralTimeZone = time.Local
	}

	start2023 := time.Date(2023, 11, 10, 0, 0, 0, 0, CentralTimeZone)
	Year2023 = Dates{
		Start: start2023,
		End:   start2023.Add(time.Hour * 24 * 3),
	}

	start2024 := time.Date(2024, 11, 8, 0, 0, 0, 0, CentralTimeZone)
	Year2024 = Dates{
		Start: start2024,
		End:   start2024.Add(time.Hour * 24 * 3),
	}
}
