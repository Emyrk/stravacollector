package superlative

import (
	"time"

	"github.com/Emyrk/strava/api/modelsdk/sdktype"
	"github.com/Emyrk/strava/database"
)

type List struct {
	EarliestStart Entry[time.Time] `json:"early_bird"`
	LatestEnd     Entry[time.Time] `json:"night_owl"`
	// MostStoppage value is in seconds
	MostStoppage Entry[int64] `json:"most_stoppage"`
	// LeastStoppage value is in seconds
	LeastStoppage    Entry[int64]   `json:"least_stoppage"`
	MostWatts        Entry[float64] `json:"most_watts"`
	MostCadence      Entry[float64] `json:"most_cadence"`
	LeastCadence     Entry[float64] `json:"least_cadence"`
	MostSuffer       Entry[int]     `json:"most_suffer"`
	MostAchievements Entry[int]     `json:"most_achievements"`
	LongestRide      Entry[float64] `json:"longest_ride"`
	ShortestRide     Entry[float64] `json:"shortest_ride"`

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

		// TODO: Loop over efforts
		//if list.MostWatts.Value == 0 || list.MostWatts.Value < activity. {
		//	list.MostWatts = entry(activity.ActivityID, activity.AverageWatts)
		//}

		//if list.MostCadence.Value == 0 || list.MostCadence.Value < activity. {
		//	list.MostCadence = entry(activity.ActivityID, activity.AverageCadence)
		//}

		//if list.LeastCadence.Value == 0 || list.LeastCadence.Value > activity. {
		//	list.LeastCadence = entry(activity.ActivityID, activity.AverageCadence)
		//}

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
