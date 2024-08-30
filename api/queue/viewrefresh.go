package queue

import (
	"context"
	"fmt"
	"time"
)

func (m *Manager) refreshViews(ctx context.Context) {
	logger := m.Logger.With().Str("job", "refresh_views").Logger()
	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("View refresh ended")
			return
		default:
		}

		time.Sleep(time.Minute * 5)

		start := time.Now()
		err := m.DB.RefreshHugelActivities(ctx)
		logger.Error().
			Err(err).
			Str("duration", fmt.Sprintf("%.3fs", time.Since(start).Seconds())).
			Msg("refresh views")
	}
}
