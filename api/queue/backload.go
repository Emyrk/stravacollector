package queue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Emyrk/strava/strava/stravalimit"

	"github.com/Emyrk/strava/database"
	"golang.org/x/oauth2"

	"github.com/Emyrk/strava/strava"
)

const backloadWait = time.Second * 30

func (m *Manager) BackLoadAthleteRoutine(ctx context.Context) {
	logger := m.Logger.With().Str("job", "backload_athlete_data").Logger()
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("Back loading athletes ended")
			return
		default:
		}

		if ok, limitLogger := stravalimit.CanLogger(1, 100, logger); !ok {
			// Do not nuke our api rate limits
			limitLogger.Error().
				Msg("hitting strava rate limit, job will try again later")
			time.Sleep(backloadWait)
			continue
		}

		// Fetch an that needs some loading.
		athlete := m.athleteToLoad(ctx)
		if athlete == nil {
			// No athletes to load, wait a bit.
			time.Sleep(backloadWait)
			continue
		}

		err := m.backloadAthlete(ctx, *athlete)
		if err != nil {
			// This could be bad
			_, dbErr := m.DB.UpsertAthleteLoad(ctx, database.UpsertAthleteLoadParams{
				AthleteID:                  athlete.AthleteID,
				LastBackloadActivityStart:  athlete.LastBackloadActivityStart,
				LastLoadAttempt:            time.Now(),
				LastLoadIncomplete:         false,
				LastLoadError:              err.Error(),
				ActivitesLoadedLastAttempt: 0,
			})
			logger.Error().
				AnErr("db_error", dbErr).
				Err(err).
				Msg("backload athlete failed")
			continue
		}
	}

}

func (m *Manager) athleteToLoad(ctx context.Context) *database.GetAthleteNeedsLoadRow {
	// Fetch an that needs some loading.
	athlete, err := m.DB.GetAthleteNeedsLoad(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		m.Logger.Error().Err(err).Msg("get athlete to load")
		return nil
	}

	// If the athlete is incomplete, always return
	if athlete.LastLoadIncomplete {
		return &athlete
	}

	// If it has been over 24 hours, return it
	if time.Since(athlete.LastLoadAttempt) > time.Hour*24 {
		return &athlete
	}

	return nil
}

// backloadAthlete tries to make progress backloading activities for some athlete.
func (m *Manager) backloadAthlete(ctx context.Context, athlete database.GetAthleteNeedsLoadRow) error {
	logger := m.Logger.With().Int64("athlete_id", athlete.AthleteID).Logger()

	// Make progress on the athlete
	logger = logger.With().Int64("athlete_id", athlete.AthleteID).Logger()

	cli := strava.NewOAuthClient(m.OAuthCfg.Client(ctx, &oauth2.Token{
		AccessToken:  athlete.OauthAccessToken,
		TokenType:    athlete.OauthTokenType,
		RefreshToken: athlete.OauthRefreshToken,
		Expiry:       athlete.OauthExpiry,
	}))

	activities, err := cli.GetActivities(ctx, strava.GetActivitiesParams{
		// Load everything after the last backload activity we have saved.
		After:   athlete.LastBackloadActivityStart,
		Page:    0,
		PerPage: 50,
	})
	if err != nil {
		return fmt.Errorf("get activities: %w", err)
	}

	logger.Debug().
		Int("activities", len(activities)).
		Time("last_backload", athlete.LastBackloadActivityStart).
		Int64("last_backload_unix", athlete.LastBackloadActivityStart.Unix()).
		Msg("backloading athlete")

	// No activities means we are done.
	if len(activities) == 0 {
		_, err := m.DB.UpsertAthleteLoad(ctx, database.UpsertAthleteLoadParams{
			AthleteID: athlete.AthleteID,
			// This did not change
			LastBackloadActivityStart:  athlete.LastBackloadActivityStart,
			LastLoadAttempt:            time.Now(),
			LastLoadIncomplete:         false,
			ActivitesLoadedLastAttempt: 0,
		})
		if err != nil {
			return fmt.Errorf("update athlete load: %w", err)
		}
		return nil
	}

	err = m.DB.InTx(func(store database.Store) error {
		for _, act := range activities {
			_, err := m.DB.UpsertMapSummary(ctx, database.UpsertMapSummaryParams{
				ID:              act.Map.ID,
				SummaryPolyline: act.Map.SummaryPolyline,
			})
			if err != nil {
				return fmt.Errorf("upsert map summary (%d): %w", act.ID, err)
			}

			_, err = m.DB.UpsertActivitySummary(ctx, database.UpsertActivitySummaryParams{
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

			err = m.EnqueueFetchActivity(ctx, athlete.AthleteID, act.ID)
			if err != nil {
				return fmt.Errorf("enqueue fetch activity: %w", err)
			}
		}
		lastAct := activities[len(activities)-1]

		_, err := store.UpsertAthleteLoad(ctx, database.UpsertAthleteLoadParams{
			AthleteID:                  athlete.AthleteID,
			LastBackloadActivityStart:  lastAct.StartDate,
			LastLoadAttempt:            time.Now(),
			LastLoadIncomplete:         true,
			ActivitesLoadedLastAttempt: int32(len(activities)),
		})
		if err != nil {
			return fmt.Errorf("update athlete load after loading: %w", err)
		}
		return nil
	}, nil)
	if err != nil {
		return fmt.Errorf("in tx: %w", err)
	}

	return nil
}
