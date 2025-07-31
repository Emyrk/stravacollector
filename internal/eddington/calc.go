package eddington

import (
	"github.com/Emyrk/strava/database"
)

func FromActivities(acts []database.EddingtonActivitiesRow) Sums {
	edds := Sums{}
	for _, act := range acts {
		edds.Add(int(database.DistanceToMiles(act.Distance)))
	}
	return edds
}

// Sums is a slice of integers representing the total number of entries with at
// least a value equal to the index.
type Sums []int32

func (e Sums) Current() int32 {
	for need, have := range e {
		need = need + 1 // 1-indexed
		if int32(need) > have {
			return int32(need) - 1
		}
	}

	if len(e) == 0 {
		return 0
	}

	return e[len(e)-1]
}

func (e *Sums) Add(value int) {
	if value < 0 {
		return
	}
	if *e == nil {
		*e = make(Sums, 0, value)
	}
	if value > len(*e) {
		*e = append(*e, make(Sums, value-len(*e))...)
	}
	for i := 0; i < value; i++ {
		(*e)[i]++
	}
}
