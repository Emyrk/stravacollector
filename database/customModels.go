package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type HugelSegmentEfforts []HugelSegmentEffort

type HugelSegmentEffort struct {
	ActivityID   int64     `json:"activity_id"`
	EffortID     int64     `json:"effort_id"`
	StartDate    time.Time `json:"start_date"`
	SegmentID    int       `json:"segment_id"`
	ElapsedTime  int       `json:"elapsed_time"`
	MovingTime   int       `json:"moving_time"`
	DeviceWatts  bool      `json:"device_watts"`
	AverageWatts float64   `json:"average_watts"`
}

func (a *HugelSegmentEfforts) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return json.Unmarshal([]byte(v), &a)
	case []byte:
		return json.Unmarshal(v, &a)
	}
	return fmt.Errorf("unexpected type %T", src)
}

func (a *HugelSegmentEfforts) Value() (driver.Value, error) {
	return json.Marshal(a)
}
