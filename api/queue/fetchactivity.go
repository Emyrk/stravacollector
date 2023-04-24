package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vgarvardt/gue/v5"
	"golang.org/x/oauth2"

	"github.com/Emyrk/strava/strava"
)

type fetchActivityJobArgs struct {
	ActivityID int64 `json:"activity_id"`
	AthleteID  int64 `json:"athlete_id"`
}

func (m *Manager) EnqueueFetchActivity(ctx context.Context, athleteID int64, activityID int64) error {
	data, err := json.Marshal(fetchActivityJobArgs{
		ActivityID: activityID,
		AthleteID:  athleteID,
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
	err := m.stravaCheck(j, 1)
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

	// Only track athletes we have in our database
	athlete, err := m.DB.GetAthleteLogin(ctx, args.AthleteID)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Error().Err(err).Msg("athlete not found, job abandoned")
		return nil
	}

	if err != nil {
		logger.Error().Err(err).Msg("job failed to get athlete from DB, will retry")
		return err
	}

	cli := strava.NewOAuthClient(m.OAuthCfg.Client(ctx, &oauth2.Token{
		AccessToken:  athlete.OauthAccessToken,
		TokenType:    athlete.OauthTokenType,
		RefreshToken: athlete.OauthRefreshToken,
		Expiry:       athlete.OauthExpiry,
	}))

	activity, err := cli.GetActivity(ctx, args.ActivityID, true)
	if err != nil {
		logger.Error().Err(err).Msg("job failed to fetch activity, will retry")
		return err
	}

	m.Logger.Info().Interface("activity", activity).Msg("activity fetched")

	return nil
}
