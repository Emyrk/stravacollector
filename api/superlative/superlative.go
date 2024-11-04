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

	//Most elevation = mountain climber
	//Best of $Segment
}

type Entry[T comparable] struct {
	Activity sdktype.StringInt `json:"activity_id"`
	Value    T                 `json:"value"`
}

func entry[T comparable](id int64, v T) Entry[T] {
	return Entry[T]{Activity: sdktype.StringInt(id), Value: v}
}

func Parse(activities []database.HugelLeaderboardRow) List {
	var list List

	for _, activity := range activities {
		if list.EarliestStart.Value.IsZero() || list.EarliestStart.Value.After(activity.StartDate) {
			list.EarliestStart = entry(activity.ActivityID, activity.StartDate)
		}

		endDate := activity.StartDate.Add(time.Duration(activity.ElapsedTime) * time.Second)
		if list.LatestEnd.Value.IsZero() || list.LatestEnd.Value.Before(endDate) {
			list.LatestEnd = entry(activity.ActivityID, endDate)
		}

		stoppage := int64(activity.ElapsedTime - activity.MovingTime)

		if list.MostStoppage.Value == 0 || list.MostStoppage.Value < stoppage {
			list.MostStoppage = entry(activity.ActivityID, stoppage)
		}

		if list.LeastStoppage.Value == 0 || list.LeastStoppage.Value > stoppage {
			list.LeastStoppage = entry(activity.ActivityID, stoppage)
		}

		if list.MostAverageWatts.Value == 0 || list.MostAverageWatts.Value < activity.AverageWatts {
			list.MostAverageWatts = entry(activity.ActivityID, activity.AverageWatts)
		}

		if list.MostAverageCadence.Value == 0 || list.MostAverageCadence.Value < activity.AverageCadence {
			list.MostAverageCadence = entry(activity.ActivityID, activity.AverageCadence)
		}

		if list.LeastAverageCadence.Value == 0 || list.LeastAverageCadence.Value > activity.AverageCadence {
			list.LeastAverageCadence = entry(activity.ActivityID, activity.AverageCadence)
		}

		if list.MostAverageSpeed.Value == 0 || list.MostAverageSpeed.Value < activity.AverageSpeed {
			list.MostAverageSpeed = entry(activity.ActivityID, activity.AverageSpeed)
		}

		if list.LeastAverageSpeed.Value == 0 || list.LeastAverageSpeed.Value > activity.AverageSpeed {
			list.LeastAverageSpeed = entry(activity.ActivityID, activity.AverageSpeed)
		}

		if list.MostAverageHeartRate.Value == 0 || list.MostAverageHeartRate.Value < activity.AverageHeartrate {
			list.MostAverageHeartRate = entry(activity.ActivityID, activity.AverageHeartrate)
		}

		if list.LeastAverageHeartRate.Value == 0 || list.LeastAverageHeartRate.Value > activity.AverageHeartrate {
			list.LeastAverageHeartRate = entry(activity.ActivityID, activity.AverageHeartrate)
		}

		if list.MostSuffer.Value == 0 || list.MostSuffer.Value < int(activity.SufferScore) {
			list.MostSuffer = entry(activity.ActivityID, int(activity.SufferScore))
		}

		if list.MostAchievements.Value == 0 || list.MostAchievements.Value < int(activity.AchievementCount) {
			list.MostAchievements = entry(activity.ActivityID, int(activity.AchievementCount))
		}

		if list.LongestRide.Value == 0 || list.LongestRide.Value < activity.Distance {
			list.LongestRide = entry(activity.ActivityID, activity.Distance)
		}

		if list.ShortestRide.Value == 0 || list.ShortestRide.Value > activity.Distance {
			list.ShortestRide = entry(activity.ActivityID, activity.Distance)
		}
	}
	return list
}
