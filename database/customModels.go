package database

import (
	"database/sql/driver"
	"strconv"

	"github.com/lib/pq"
)

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

func (a *Floats) Value() (driver.Value, error) {
	return pq.Array(a).Value()
}
