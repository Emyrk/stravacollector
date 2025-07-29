package river

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/strava"
	"github.com/Emyrk/strava/strava/stravalimit"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

func (m *Manager) EnqueueReloadSegments(ctx context.Context, opts ...func(j *river.InsertOpts)) (bool, error) {
	iopts := &river.InsertOpts{}
	for _, opt := range opts {
		opt(iopts)
	}

	fi, err := m.cli.Insert(ctx, ReloadSegmentsArgs{}, iopts)

	skipped := false
	if fi != nil {
		skipped = fi.UniqueSkippedAsDuplicate
	}

	return !skipped, err
}

type ReloadSegmentsArgs struct {
}

func (ReloadSegmentsArgs) Kind() string { return "reload_segments" }
func (ReloadSegmentsArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:    riverDatabaseQueue,
		Priority: PriorityHighest,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: time.Minute * 5,
		},
	}
}

type ReloadSegmentsWorker struct {
	mgr *Manager
	river.WorkerDefaults[ReloadSegmentsArgs]
}

func (*ReloadSegmentsWorker) Middleware(job *rivertype.JobRow) []rivertype.WorkerMiddleware {
	return []rivertype.WorkerMiddleware{}
}

func (w *ReloadSegmentsWorker) Work(ctx context.Context, job *river.Job[ReloadSegmentsArgs]) error {
	logger := jobLogFields(w.mgr.logger, job)

	// This is "Steven Masley"
	ath, err := w.mgr.db.GetAthleteLogin(ctx, 2661162)
	if err != nil {
		logger.Error().Err(err).Msg("Steven Masley is required to load segments")
		return errors.New("Steven Masley is required to load segments")
	}

	iBuf, dBuf := int64(100), int64(500)
	if stravalimit.NextDailyReset(time.Now()) < time.Hour*3 {
		iBuf, dBuf = 50, 500
	}

	routes, err := w.mgr.db.AllCompetitiveRoutes(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch competitive routes")
		return err
	}

	neededSegments := make(map[int64]int, 0)
	for _, route := range routes {
		for _, seg := range route.Segments {
			neededSegments[seg] += 1
		}
	}

	loaded, err := w.mgr.db.LoadedSegments(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch loaded segments")
		return err
	}

	for _, seg := range loaded {
		delete(neededSegments, seg.ID)
	}

	if len(neededSegments) == 0 {
		_ = river.RecordOutput(ctx, "no segments to load")
		return nil
	}

	logger.Debug().Int("needed", len(neededSegments)).Msg("need to load segments")
	if ok, limitLogger := stravalimit.CanLogger(int64(len(neededSegments)), iBuf, dBuf, logger); !ok {
		// Do not nuke our api rate limits
		limitLogger.Error().
			Str("job", "backload_segment_data").
			Msg("hitting strava rate limit, job will try again later")
		return w.mgr.StravaSnooze(ctx)
	}

	cli := strava.NewOAuthClient(w.mgr.oauthCfg.Client(ctx, ath.OAuthToken()))
	for segmentID := range neededSegments {
		_ = river.RecordOutput(ctx, fmt.Sprintf("loading segment https://www.strava.com/segments/%d", segmentID))
		segment, err := cli.GetSegment(ctx, segmentID)
		if err != nil {
			return err
		}

		err = w.mgr.db.InTx(func(store database.Store) error {
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
			return fmt.Errorf("failed to upsert segment data: %w", err)
		}
	}
	_ = river.RecordOutput(ctx, fmt.Sprintf("%d segments loaded", len(neededSegments)))
	return nil
}
