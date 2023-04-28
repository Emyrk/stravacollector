package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Emyrk/strava/database"
	"github.com/vgarvardt/gue/v5"
	"golang.org/x/oauth2"

	"github.com/Emyrk/strava/strava"
)

type fetchActivityJobArgs struct {
	Source     database.ActivityDetailSource `json:"source"`
	ActivityID int64                         `json:"activity_id"`
	AthleteID  int64                         `json:"athlete_id"`
}

func (m *Manager) EnqueueFetchActivity(ctx context.Context, source database.ActivityDetailSource, athleteID int64, activityID int64) error {
	data, err := json.Marshal(fetchActivityJobArgs{
		ActivityID: activityID,
		AthleteID:  athleteID,
		Source:     source,
	})
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	return m.Client.Enqueue(ctx, &gue.Job{
		Type:  fetchActivityJob,
		Queue: stravaFetchQueue,
		Args:  data,
	})
}

func (m *Manager) fetchActivity(ctx context.Context, j *gue.Job) error {
	err := m.jobStravaCheck(j, 1)
	if err != nil {
		return err
	}

	logger := jobLogFields(m.Logger, j)

	var args fetchActivityJobArgs
	err = json.Unmarshal(j.Args, &args)
	if err != nil {
		logger.Error().Err(err).Msg("json unmarshal, job abandoned")
		return nil
	}
	if args.Source == "" {
		args.Source = database.ActivityDetailSourceUnknown
	}

	// Only track athletes we have in our database
	athlete, err := m.DB.GetAthleteLogin(ctx, args.AthleteID)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Error().Err(err).Msg("athlete not found, job abandoned")
		return nil
	}

	if err != nil {
		return err
	}
	logger = logger.With().Int64("athlete_id", athlete.AthleteID).Logger()

	cli := strava.NewOAuthClient(m.OAuthCfg.Client(ctx, &oauth2.Token{
		AccessToken:  athlete.OauthAccessToken,
		TokenType:    athlete.OauthTokenType,
		RefreshToken: athlete.OauthRefreshToken,
		Expiry:       athlete.OauthExpiry,
	}))

	activity, err := cli.GetActivity(ctx, args.ActivityID, true)
	if err != nil {
		return err
	}

	logger.Info().Int64("activity_id", activity.ID).Msg("activity fetched")

	// Parse the activity, save all efforts.
	err = m.DB.InTx(func(store database.Store) error {
		_, err := store.UpsertMap(ctx, database.UpsertMapParams{
			ID:              activity.Map.ID,
			Polyline:        activity.Map.Polyline,
			SummaryPolyline: activity.Map.SummaryPolyline,
		})
		if err != nil {
			return fmt.Errorf("upsert map: %w", err)
		}

		_, err = store.UpsertActivitySummary(ctx, database.UpsertActivitySummaryParams{
			ID:                 activity.ID,
			AthleteID:          activity.Athlete.ID,
			UploadID:           activity.UploadID,
			ExternalID:         activity.ExternalID,
			Name:               activity.Name,
			Distance:           activity.Distance,
			MovingTime:         activity.MovingTime,
			ElapsedTime:        activity.ElapsedTime,
			TotalElevationGain: activity.TotalElevationGain,
			ActivityType:       activity.Type,
			SportType:          activity.SportType,
			WorkoutType:        activity.WorkoutType,
			StartDate:          activity.StartDate,
			StartDateLocal:     activity.StartDateLocal,
			Timezone:           activity.Timezone,
			UtcOffset:          activity.UtcOffset,
			AchievementCount:   activity.AchievementCount,
			KudosCount:         activity.KudosCount,
			CommentCount:       activity.CommentCount,
			AthleteCount:       activity.AthleteCount,
			PhotoCount:         activity.PhotoCount,
			MapID:              activity.Map.ID,
			Trainer:            activity.Trainer,
			Commute:            activity.Commute,
			Manual:             activity.Manual,
			Private:            activity.Private,
			Flagged:            activity.Flagged,
			GearID:             activity.GearID,
			AverageSpeed:       activity.AverageSpeed,
			MaxSpeed:           activity.MaxSpeed,
			DeviceWatts:        activity.DeviceWatts,
			HasHeartrate:       activity.HasHeartrate,
			PrCount:            activity.PrCount,
			TotalPhotoCount:    activity.TotalPhotoCount,
			AverageHeartrate:   activity.AverageHeartrate,
			MaxHeartrate:       activity.MaxHeartrate,
		})
		if err != nil {
			return fmt.Errorf("upsert activity summary: %w", err)
		}

		_, err = store.UpsertActivityDetail(ctx, database.UpsertActivityDetailParams{
			ID:                       activity.ID,
			AthleteID:                activity.Athlete.ID,
			StartLatlng:              activity.StartLatlng,
			EndLatlng:                activity.EndLatlng,
			MapID:                    activity.Map.ID,
			FromAcceptedTag:          activity.FromAcceptedTag,
			AverageCadence:           activity.AverageCadence,
			AverageTemp:              activity.AverageTemp,
			AverageWatts:             activity.AverageWatts,
			WeightedAverageWatts:     activity.WeightedAverageWatts,
			Kilojoules:               activity.Kilojoules,
			MaxWatts:                 activity.MaxWatts,
			ElevHigh:                 activity.ElevHigh,
			ElevLow:                  activity.ElevLow,
			SufferScore:              int32(activity.SufferScore),
			EmbedToken:               activity.EmbedToken,
			SegmentLeaderboardOptOut: activity.SegmentLeaderboardOptOut,
			LeaderboardOptOut:        activity.LeaderboardOptOut,
			Calories:                 activity.Calories,
			Source:                   args.Source,
			//
			PremiumFetch:      athlete.Summit,
			NumSegmentEfforts: int32(len(activity.SegmentEfforts)),
		})
		if err != nil {
			return fmt.Errorf("upsert activity details: %w", err)
		}

		// Insert efforts.
		for i, effort := range activity.SegmentEfforts {
			_, err := store.UpsertSegmentEffort(ctx, database.UpsertSegmentEffortParams{
				ID:             effort.ID,
				AthleteID:      effort.Athlete.ID,
				SegmentID:      effort.Segment.ID,
				Name:           effort.Name,
				ElapsedTime:    effort.ElapsedTime,
				MovingTime:     effort.MovingTime,
				StartDate:      effort.StartDate,
				StartDateLocal: effort.StartDateLocal,
				Distance:       effort.Distance,
				StartIndex:     effort.StartIndex,
				EndIndex:       effort.EndIndex,
				DeviceWatts:    effort.DeviceWatts,
				AverageWatts:   effort.AverageWatts,
				KomRank: sql.NullInt32{
					Int32: effort.KomRank,
					Valid: effort.KomRank != 0,
				},
				PrRank: sql.NullInt32{
					Int32: effort.PrRank,
					Valid: effort.PrRank != 0,
				},
				ActivitiesID: activity.ID,
			})
			if err != nil {
				return fmt.Errorf("upsert segment effort index=%d, id=%d: %w", i, effort.ID, err)
			}
		}

		return nil
	}, nil)
	if err != nil {
		return fmt.Errorf("in tx: %w", err)
	}

	logger.Info().Int64("activity_id", activity.ID).Msg("activity inserted!")

	return nil
}
