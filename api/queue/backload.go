package queue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vgarvardt/gue/v5"

	"github.com/Emyrk/strava/strava/stravalimit"

	"github.com/Emyrk/strava/database"

	"github.com/Emyrk/strava/strava"
)

const backloadWait = time.Second * 30

func (m *Manager) BackLoadAthleteRoutine(ctx context.Context) {
	logger := m.Logger.With().Str("job", "backload_athlete_data").Logger()
	lastRateLimitLog := time.Time{}
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("Back loading athletes ended")
			return
		default:
		}
		now := time.Now()

		iBuf, dBuf := int64(150), int64(500)

		if stravalimit.NextDailyReset(now) < time.Hour*3 {
			iBuf, dBuf = 80, 300
		}
		if stravalimit.NextDailyReset(now) < time.Hour*1 {
			iBuf, dBuf = 50, 150
		}
		if stravalimit.NextDailyReset(now) < time.Minute*20 {
			iBuf, dBuf = 50, 100
		}

		if ok, limitLogger := stravalimit.CanLogger(1, iBuf, dBuf, logger); !ok {
			if time.Since(lastRateLimitLog) > time.Minute*60 { // Debounce this
				// Do not nuke our api rate limits
				limitLogger.Error().
					Str("job", "backload_athlete_data").
					Msg("hitting strava rate limit, job will try again later")
				lastRateLimitLog = time.Now()
			}

			time.Sleep(backloadWait)
			continue
		}

		// Fetch an athlete that needs some loading.
		athlete := m.athleteToLoad(ctx)
		if athlete == nil {
			// No athletes to load, wait a bit.
			time.Sleep(time.Minute * 15)
			continue
		}

		start := time.Now()
		err := m.backloadAthlete(ctx, *athlete)
		m.backloadHistogram.WithLabelValues(strconv.FormatBool(err == nil)).Observe(time.Since(start).Seconds())
		if err != nil {
			sleepFor := time.Second
			next := time.Now().Add(time.Hour)
			if se := strava.IsAPIError(err); se != nil {
				if se.Response.StatusCode == http.StatusTooManyRequests {
					next = time.Now().Add(time.Minute * 15)
					err = fmt.Errorf("wait at least 15min: %w", err)
					sleepFor = time.Minute
				}

				if se.Response.StatusCode == http.StatusUnauthorized || se.Response.StatusCode == http.StatusForbidden {
					// This person needs to be fixed....
					// We should delete them?
					// TODO: Handle these people.
					next = time.Now().Add(time.Hour * 48)
					err = fmt.Errorf("unauthorized: %w", err)
				}
			} else {
				// Hitting disk IO rate limit on supabase
				sleepFor = time.Second * 15
			}
			// This could be bad
			_, dbErr := m.DB.UpsertAthleteLoad(ctx, database.UpsertAthleteLoadParams{
				AthleteID:                  athlete.AthleteLoad.AthleteID,
				LastBackloadActivityStart:  athlete.AthleteLoad.LastBackloadActivityStart,
				LastLoadAttempt:            time.Now(),
				LastLoadIncomplete:         false,
				LastLoadError:              err.Error(),
				ActivitesLoadedLastAttempt: 0,
				EarliestActivityID:         athlete.AthleteLoad.EarliestActivityID,
				EarliestActivity:           athlete.AthleteLoad.EarliestActivity,
				EarliestActivityDone:       athlete.AthleteLoad.EarliestActivityDone,
				NextLoadNotBefore:          next.UTC(),
			})
			logger.Error().
				Int64("athlete_id", athlete.AthleteLogin.AthleteID).
				AnErr("db_error", dbErr).
				Time("next_load_not_before", next).
				Err(err).
				Dur("sleep_for", sleepFor).
				Msg("backload athlete failed")
			time.Sleep(sleepFor)
			continue
		}
	}
}

func (m *Manager) athleteToLoad(ctx context.Context) *database.GetAthleteNeedsLoadRow {
	// Fetch an that needs some loading.
	athletes, err := m.DB.GetAthleteNeedsLoad(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		m.Logger.Error().Err(err).Msg("get athlete to load")
		return nil
	}

	for _, athlete := range athletes {
		athlete := athlete

		// If the athlete is incomplete, always return
		if athlete.AthleteLoad.LastLoadIncomplete {
			return &athlete
		}

		if !athlete.AthleteLoad.EarliestActivityDone {
			return &athlete
		}

		// If it has been over 24 hours, return it
		if time.Since(athlete.AthleteLoad.LastLoadAttempt) > time.Hour*24 {
			return &athlete
		}
	}

	return nil
}

// backloadAthlete tries to make progress backloading activities for some athlete.
func (m *Manager) backloadAthlete(ctx context.Context, athlete database.GetAthleteNeedsLoadRow) error {
	logger := m.Logger.With().Int64("athlete_id", athlete.AthleteLogin.AthleteID).Logger()

	// Make progress on the athlete
	logger = logger.With().Int64("athlete_id", athlete.AthleteLogin.AthleteID).Logger()

	cli := strava.NewOAuthClient(m.OAuthCfg.Client(ctx, athlete.AthleteLogin.OAuthToken()))

	params := strava.GetActivitiesParams{
		Page:    0,
		PerPage: 50,
	}
	backloadingHistory := false
	athleteLoad := athlete.AthleteLoad
	if !athleteLoad.EarliestActivityDone {
		params.Before = athleteLoad.EarliestActivity.Add(time.Second * -1)
		backloadingHistory = true
	} else {
		params.After = athleteLoad.LastBackloadActivityStart
	}

	activities, err := cli.GetActivities(ctx, params)
	if err != nil {
		return fmt.Errorf("get activities: %w", err)
	}

	logger.Debug().
		Int("activities", len(activities)).
		Time("last_backload", athleteLoad.LastBackloadActivityStart).
		Int64("last_backload_unix", athleteLoad.LastBackloadActivityStart.Unix()).
		Time("param_before", params.Before).
		Time("param_after", params.After).
		Int("param_page", params.Page).
		Int("param_per_page", params.PerPage).
		Msg("backloading athlete")

	// No activities means we are done.
	if len(activities) == 0 {
		_, err := m.DB.UpsertAthleteLoad(ctx, database.UpsertAthleteLoadParams{
			AthleteID: athleteLoad.AthleteID,
			// This did not change
			LastBackloadActivityStart:  athleteLoad.LastBackloadActivityStart,
			LastLoadAttempt:            time.Now(),
			LastLoadIncomplete:         false,
			LastLoadError:              "",
			ActivitesLoadedLastAttempt: 0,
			EarliestActivity:           athleteLoad.EarliestActivity,
			EarliestActivityID:         athleteLoad.EarliestActivityID,
			EarliestActivityDone:       true,
			NextLoadNotBefore:          time.Now().Add(time.Minute * 15),
		})
		if err != nil {
			return fmt.Errorf("update athlete load: %w", err)
		}
		return nil
	}

	err = m.DB.InTx(func(store database.Store) error {
		for _, act := range activities {
			_, err := store.UpsertMapData(ctx, database.UpsertMapDataParams{
				ID:              act.Map.ID,
				SummaryPolyline: act.Map.SummaryPolyline,
			})
			if err != nil {
				return fmt.Errorf("upsert map summary (%d): %w", act.ID, err)
			}

			_, err = store.UpsertActivitySummary(ctx, database.UpsertActivitySummaryParams{
				ID:                 act.ID,
				AthleteID:          act.Athlete.ID,
				UploadID:           act.UploadID,
				ExternalID:         act.ExternalID,
				Name:               act.Name,
				Distance:           act.Distance,
				MovingTime:         act.MovingTime,
				ElapsedTime:        act.ElapsedTime,
				TotalElevationGain: act.TotalElevationGain,
				ActivityType:       act.Type,
				SportType:          act.SportType,
				WorkoutType:        act.WorkoutType,
				StartDate:          act.StartDate,
				StartDateLocal:     act.StartDateLocal,
				Timezone:           act.Timezone,
				UtcOffset:          act.UtcOffset,
				AchievementCount:   act.AchievementCount,
				KudosCount:         act.KudosCount,
				CommentCount:       act.CommentCount,
				AthleteCount:       act.AthleteCount,
				PhotoCount:         act.PhotoCount,
				MapID:              act.Map.ID,
				Trainer:            act.Trainer,
				Commute:            act.Commute,
				Manual:             act.Manual,
				Private:            act.Private,
				Flagged:            act.Flagged,
				GearID:             act.GearID,
				AverageSpeed:       act.AverageSpeed,
				MaxSpeed:           act.MaxSpeed,
				DeviceWatts:        act.DeviceWatts,
				HasHeartrate:       act.HasHeartrate,
				PrCount:            act.PrCount,
				TotalPhotoCount:    act.TotalPhotoCount,
				AverageHeartrate:   act.AverageHeartrate,
				MaxHeartrate:       act.MaxHeartrate,
			})
			if err != nil {
				return fmt.Errorf("upsert activity summary (%d): %w", act.ID, err)
			}

			// Backload bike rides for more deets
			if isBikeRide(act.Type) || isBikeRide(act.SportType) {
				err = m.EnqueueFetchActivity(ctx, database.ActivityDetailSourceBackload, athleteLoad.AthleteID, act.ID, canBeHugel(act) || canBeHugelLite(act), onHugelDate(act), activityGuePriority(act), func(j *gue.Job) {
					// Delay by 5minutes.
					// We do this because sometimes strava loads 0 segments for a ride, and it takes some time
					// for segments to be populated.
					j.RunAt = time.Now().Add(time.Minute * 5)
				})
				if err != nil {
					return fmt.Errorf("enqueue fetch activity: %w", err)
				}
			}
		}
		first := activities[len(activities)-1]
		lastActStart := activities[0].StartDate
		if athleteLoad.LastBackloadActivityStart.After(lastActStart) {
			lastActStart = athleteLoad.LastBackloadActivityStart
		}

		params := database.UpsertAthleteLoadParams{
			AthleteID:                  athleteLoad.AthleteID,
			LastBackloadActivityStart:  lastActStart,
			LastLoadAttempt:            time.Now(),
			LastLoadIncomplete:         true,
			LastLoadError:              "",
			ActivitesLoadedLastAttempt: int32(len(activities)),
			EarliestActivityID:         athleteLoad.EarliestActivityID,
			EarliestActivity:           athleteLoad.EarliestActivity,
			EarliestActivityDone:       athleteLoad.EarliestActivityDone,
			// When we are not done, do not prevent loading more.
			NextLoadNotBefore: time.Now(),
		}
		if backloadingHistory {
			params.EarliestActivity = first.StartDate
			params.EarliestActivityID = first.ID
			params.EarliestActivityDone = false
		}
		_, err := store.UpsertAthleteLoad(ctx, params)
		if err != nil {
			return fmt.Errorf("update athlete load after loading: %w", err)
		}
		return nil
	}, nil)
	if err != nil {
		return fmt.Errorf("in tx: %w", err)
	}
	m.backloadActivitiesLoaded.Add(float64(len(activities)))

	return nil
}

func canBeHugel(summary strava.ActivitySummary) bool {
	return database.DistanceToMiles(summary.Distance) > 80 &&
		database.DistanceToFeet(summary.TotalElevationGain) > 8000
}

func canBeHugelLite(summary strava.ActivitySummary) bool {
	return database.DistanceToMiles(summary.Distance) > 35 &&
		database.DistanceToFeet(summary.TotalElevationGain) > 3500
}

func onHugelDate(summary strava.ActivitySummary) bool {
	if summary.StartDate.Year() == 2024 && summary.StartDate.Month() == time.November {
		if summary.StartDate.Day() > 7 && summary.StartDate.Day() < 12 {
			return true
		}

	}
	return false
}

func activityGuePriority(summary strava.ActivitySummary) gue.JobPriority {
	priority := gue.JobPriorityDefault
	if time.Since(summary.StartDate) < (time.Hour * 24 * 7) {
		// Add priority for recent activities
		priority -= 2000
	}

	if database.DistanceToMiles(summary.Distance) > 80 {
		// Priority for long activites
		priority -= 500
	}

	if database.DistanceToMiles(summary.Distance) > 40 {
		// Longer rides should just sync before shorter rides.
		// Even if they are non-hugels. We get more value from them.
		priority -= 1
	} else {
		// Reduce priority for short rides.
		priority += 200
	}

	if database.DistanceToFeet(summary.TotalElevationGain) > 7000 {
		// Priority for high vert
		priority -= 500
	} else {
		priority += 100
	}

	return priority
}

// isBikeRide covers the weird stuff like "VirtualRide", "EBikeRide", "MountainBikeRide"
func isBikeRide(act string) bool {
	act = strings.ToLower(act)
	if strings.Contains(act, "bike") {
		return true
	}
	if strings.Contains(act, "ride") {
		return true
	}
	return false
}
