package superlative

import (
	"time"

	"github.com/Emyrk/strava/api/modelsdk/sdktype"
	"github.com/Emyrk/strava/database"
)

type List struct {
	EarliestStart Entry[time.Time] `json:"earliest_start"`
	LatestEnd     Entry[time.Time] `json:"latest_end"`
	// MostStoppage value is in seconds
	MostStoppage Entry[int64] `json:"most_stoppage"`
	// LeastStoppage value is in seconds
	LeastStoppage         Entry[int64]   `json:"least_stoppage"`
	MostAverageWatts      Entry[float64] `json:"most_avg_watts"`
	MostAverageCadence    Entry[float64] `json:"most_avg_cadence"`
	LeastAverageCadence   Entry[float64] `json:"least_avg_cadence"`
	MostAverageSpeed      Entry[float64] `json:"most_avg_speed"`
	LeastAverageSpeed     Entry[float64] `json:"least_avg_speed"`
	MostAverageHeartRate  Entry[float64] `json:"most_avg_hr"`
	LeastAverageHeartRate Entry[float64] `json:"least_avg_hr"`
	MostSuffer            Entry[int]     `json:"most_suffer"`
	MostAchievements      Entry[int]     `json:"most_achievements"`
	LongestRide           Entry[float64] `json:"longest_ride"`
	ShortestRide          Entry[float64] `json:"shortest_ride"`

	// Most elevation = mountain climber
	// Best of $Segment
}

type Entry[T comparable] struct {
	Activity sdktype.StringInt `json:"activity_id"`
	Value    T                 `json:"value"`
}

func entry[T comparable](id int64, v T) Entry[T] {
	return Entry[T]{Activity: sdktype.StringInt(id), Value: v}
}

func compare[T int | float64 | int64](entry Entry[T], lowest Entry[T], highest Entry[T]) (Entry[T], Entry[T]) {
	if entry.Value == 0 {
		// No change if the value is omitted
		return lowest, highest
	}

	if highest.Activity == 0 || entry.Value > highest.Value {
		highest = entry
	}

	if lowest.Activity == 0 || entry.Value < lowest.Value {
		lowest = entry
	}

	return lowest, highest

}

func Parse(activities []database.HugelLeaderboardRow) List {
	var list List

	for _, activity := range activities {
		if list.EarliestStart.Value.IsZero() || list.EarliestStart.Value.After(activity.StartDate.Time) {
			list.EarliestStart = entry(activity.ActivityID, activity.StartDate.Time)
		}

		endDate := activity.StartDate.Time.Add(time.Duration(activity.ElapsedTime) * time.Second)
		if list.LatestEnd.Value.IsZero() || list.LatestEnd.Value.Before(endDate) {
			list.LatestEnd = entry(activity.ActivityID, endDate)
		}

		stoppage := int64(activity.ElapsedTime - activity.MovingTime)
		list.LeastStoppage, list.MostStoppage = compare(entry(activity.ActivityID, stoppage), list.LeastStoppage, list.MostStoppage)
		list.LeastAverageCadence, list.MostAverageCadence = compare(entry(activity.ActivityID, activity.AverageCadence), list.LeastAverageCadence, list.MostAverageCadence)
		list.LeastAverageSpeed, list.MostAverageSpeed = compare(entry(activity.ActivityID, activity.AverageSpeed), list.LeastAverageSpeed, list.MostAverageSpeed)
		list.LeastAverageHeartRate, list.MostAverageHeartRate = compare(entry(activity.ActivityID, activity.AverageHeartrate), list.LeastAverageHeartRate, list.MostAverageHeartRate)
		_, list.MostSuffer = compare(entry(activity.ActivityID, int(activity.SufferScore)), entry(0, 0), list.MostSuffer)
		_, list.MostAchievements = compare(entry(activity.ActivityID, int(activity.AchievementCount)), entry(0, 0), list.MostAchievements)
		list.ShortestRide, list.LongestRide = compare(entry(activity.ActivityID, activity.Distance), list.ShortestRide, list.LongestRide)

		if activity.DeviceWatts {
			_, list.MostAverageWatts = compare(entry(activity.ActivityID, activity.AverageWatts), entry(0, float64(0)), list.MostAverageWatts)
		}
	}
	return list
}
