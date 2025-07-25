package superlative_test

import (
	"testing"
	"time"

	"github.com/Emyrk/strava/api/modelsdk/sdktype"
	"github.com/Emyrk/strava/api/superlative"
	"github.com/Emyrk/strava/database"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("NoActivities", func(t *testing.T) {
		activities := []database.HugelLeaderboardRow{}
		list := superlative.Parse(activities)
		require.Equal(t, superlative.List{}, list)
	})

	t.Run("SingleActivity", func(t *testing.T) {
		activities := []database.HugelLeaderboardRow{
			activity(1, stats{
				Stoppage:     1200,
				Watts:        100,
				Cadence:      200,
				Speed:        5,
				HeartRate:    160,
				Suffer:       200,
				Achievements: 200,
				Distance:     10000,
				Start:        time.Now().Add(time.Hour * -1),
				End:          time.Now(),
			}),
		}

		list := superlative.Parse(activities)
		require.Equal(t, list.EarliestStart, entry(1, activities[0].StartDate))
		require.Equal(t, list.MostStoppage, entry(1, int64(1200)))
		require.Equal(t, list.LeastStoppage, entry(1, int64(1200)))
		require.Equal(t, list.MostAverageWatts, entry(1, activities[0].AverageWatts))
		require.Equal(t, list.MostAverageCadence, entry(1, activities[0].AverageCadence))
		require.Equal(t, list.LeastAverageCadence, entry(1, activities[0].AverageCadence))
		require.Equal(t, list.MostAverageSpeed, entry(1, activities[0].AverageSpeed))
		require.Equal(t, list.LeastAverageSpeed, entry(1, activities[0].AverageSpeed))
		require.Equal(t, list.MostAverageHeartRate, entry(1, activities[0].AverageHeartrate))
		require.Equal(t, list.LeastAverageHeartRate, entry(1, activities[0].AverageHeartrate))
		require.Equal(t, list.MostSuffer, entry(1, int(activities[0].SufferScore)))
		require.Equal(t, list.MostAchievements, entry(1, int(activities[0].AchievementCount)))
		require.Equal(t, list.ShortestRide, entry(1, activities[0].Distance))
		require.Equal(t, list.LongestRide, entry(1, activities[0].Distance))
	})

	t.Run("TwoActivities", func(t *testing.T) {
		activities := []database.HugelLeaderboardRow{
			activity(1, stats{
				Stoppage:     1200,
				Watts:        100,
				Cadence:      200,
				Speed:        5,
				HeartRate:    160,
				Suffer:       200,
				Achievements: 200,
				Distance:     10000,
				Start:        time.Now().Add(time.Hour * -1),
				End:          time.Now(),
			}),

			activity(2, stats{
				Stoppage:     1000,
				Watts:        200,
				Cadence:      150,
				Speed:        10,
				HeartRate:    120,
				Suffer:       50,
				Achievements: 100,
				Distance:     20000,
				Start:        time.Now().Add(time.Hour * -2),
				End:          time.Now().Add(time.Hour),
			}),
		}

		list := superlative.Parse(activities)
		require.Equal(t, list.EarliestStart, entry(2, activities[1].StartDate))
		require.Equal(t, list.MostStoppage, entry(1, int64(1200)))
		require.Equal(t, list.LeastStoppage, entry(2, int64(1000)))
		require.Equal(t, list.MostAverageWatts, entry(2, activities[1].AverageWatts))
		require.Equal(t, list.MostAverageCadence, entry(1, activities[0].AverageCadence))
		require.Equal(t, list.LeastAverageCadence, entry(2, activities[1].AverageCadence))
		require.Equal(t, list.MostAverageSpeed, entry(2, activities[1].AverageSpeed))
		require.Equal(t, list.LeastAverageSpeed, entry(1, activities[0].AverageSpeed))
		require.Equal(t, list.MostAverageHeartRate, entry(1, activities[0].AverageHeartrate))
		require.Equal(t, list.LeastAverageHeartRate, entry(2, activities[1].AverageHeartrate))
		require.Equal(t, list.MostSuffer, entry(1, int(activities[0].SufferScore)))
		require.Equal(t, list.MostAchievements, entry(1, int(activities[0].AchievementCount)))
		require.Equal(t, list.ShortestRide, entry(1, activities[0].Distance))
		require.Equal(t, list.LongestRide, entry(2, activities[1].Distance))
	})

}

type stats struct {
	Stoppage     float64
	Watts        float64
	Cadence      float64
	Speed        float64
	HeartRate    float64
	Suffer       int32
	Achievements int32
	Distance     float64
	Start        time.Time
	End          time.Time
}

func activity(id int64, data stats) database.HugelLeaderboardRow {
	elapsedSeconds := data.End.Sub(data.Start).Seconds()
	movingSeconds := elapsedSeconds - data.Stoppage
	return database.HugelLeaderboardRow{
		ActivityID:         id,
		TotalTimeSeconds:   0,
		Efforts:            nil,
		Name:               "",
		DeviceWatts:        true,
		Distance:           data.Distance,
		MovingTime:         movingSeconds,
		ElapsedTime:        elapsedSeconds,
		TotalElevationGain: 0,
		StartDate:          database.Timestamptz(data.Start),
		AchievementCount:   data.Achievements,
		AverageHeartrate:   data.HeartRate,
		AverageSpeed:       data.Speed,
		SufferScore:        data.Suffer,
		AverageWatts:       data.Watts,
		AverageCadence:     data.Cadence,
		Firstname:          "",
		Lastname:           "",
		Username:           "",
		ProfilePicLink:     "",
		Sex:                "",
		HugelCount:         0,
	}
}

func entry[T comparable](id int64, v T) superlative.Entry[T] {
	return superlative.Entry[T]{Activity: sdktype.StringInt(id), Value: v}
}
