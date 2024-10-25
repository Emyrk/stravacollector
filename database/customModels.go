package database

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
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

type Floats []float64

func (a *Floats) Scan(src interface{}) error {
	var output []string
	err := pq.Array(&output).Scan(src)
	if err != nil {
		return err
	}

	vals := make([]float64, 0, len(output))
	for _, v := range output {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return err
		}
		vals = append(vals, f)
	}
	*a = vals
	return nil
}

func (a Floats) Value() (driver.Value, error) {
	return pq.Array(a).Value()
}
