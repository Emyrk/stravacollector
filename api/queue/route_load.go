package queue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Emyrk/strava/database"

	"github.com/Emyrk/strava/strava"

	"github.com/Emyrk/strava/strava/stravalimit"
)

const segmentWait = time.Minute * 10

func (m *Manager) BackLoadRouteSegments(ctx context.Context) {
	logger := m.Logger.With().Str("job", "backload_segment_data").Logger()
	first := true
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("Back loading segments ended")
			return
		default:
		}
		// First one does not sleep
		if !first {
			time.Sleep(segmentWait)
		} else {
			first = false
		}

		// This is "Steven Masley"
		ath, err := m.DB.GetAthleteLogin(ctx, 2661162)
		if err != nil {
			logger.Error().Err(err).Msg("Steven Masley is required to load segments")
			continue
		}

		iBuf, dBuf := int64(100), int64(500)
		if stravalimit.NextDailyReset(time.Now()) < time.Hour*3 {
			iBuf, dBuf = 50, 500
		}

		routes, err := m.DB.AllCompetitiveRoutes(ctx)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch competitive routes")
			continue
		}

		neededSegments := make(map[int64]int, 0)
		for _, route := range routes {
			for _, seg := range route.Segments {
				neededSegments[seg] += 1
			}
		}

		loaded, err := m.DB.LoadedSegments(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("failed to fetch loaded segments")
			continue
		}

		for _, seg := range loaded {
			delete(neededSegments, seg.ID)
		}

		if len(neededSegments) == 0 {
			continue
		}

		logger.Debug().Int("needed", len(neededSegments)).Msg("need to load segments")
		if ok, limitLogger := stravalimit.CanLogger(int64(len(neededSegments)), iBuf, dBuf, logger); !ok {
			// Do not nuke our api rate limits
			limitLogger.Error().
				Str("job", "backload_segment_data").
				Msg("hitting strava rate limit, job will try again later")
			continue
		}

		cli := strava.NewOAuthClient(m.OAuthCfg.Client(ctx, ath.OAuthToken()))
		for segmentID := range neededSegments {
			segment, err := cli.GetSegment(ctx, segmentID)
			if err != nil {
				logger.Error().Err(err).Int64("segment", segmentID).Msg("failed to fetch segment")
				continue
			}

			err = m.DB.InTx(func(store database.Store) error {
				_, err := store.UpsertMapData(ctx, database.UpsertMapDataParams{
					ID:              segment.Map.ID,
					Polyline:        segment.Map.Polyline,
					SummaryPolyline: segment.Map.SummaryPolyline,
				})
				if err != nil {
					return fmt.Errorf("failed to upsert map data: %w", err)
				}

				_, err = store.UpsertSegment(ctx, database.UpsertSegmentParams{
					ID:                 segment.ID,
					Name:               segment.Name,
					ActivityType:       segment.ActivityType,
					Distance:           segment.Distance,
					AverageGrade:       segment.AverageGrade,
					MaximumGrade:       segment.MaximumGrade,
					ElevationHigh:      segment.ElevationHigh,
					ElevationLow:       segment.ElevationLow,
					StartLatlng:        segment.StartLatlng,
					EndLatlng:          segment.EndLatlng,
					ElevationProfile:   segment.ElevationProfile,
					ClimbCategory:      segment.ClimbCategory,
					City:               segment.City,
					State:              segment.State,
					Country:            segment.Country,
					Private:            segment.Private,
					Hazardous:          segment.Hazardous,
					CreatedAt:          database.Timestamp(segment.CreatedAt),
					UpdatedAt:          database.Timestamp(segment.UpdatedAt),
					TotalElevationGain: segment.TotalElevationGain,
					MapID:              segment.Map.ID,
					TotalEffortCount:   segment.EffortCount,
					TotalAthleteCount:  segment.AthleteCount,
					TotalStarCount:     segment.StarCount,
				})
				if err != nil {
					return fmt.Errorf("failed to upsert segment data: %w", err)
				}

				return nil
			}, nil)
			if err != nil {
				logger.Error().Err(err).Int64("segment", segmentID).Msg("failed to insert segment, waiting 1 hour")
				time.Sleep(time.Hour)
				continue
			}
		}
	}
}
