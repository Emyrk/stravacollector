package leaderboards

import (
	"context"
	"time"

	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/database/gencache"
)

type Leaderboards struct {
	db database.Store

	SuperHugelBoardCache    *gencache.LazyCache[[]database.SuperHugelLeaderboardRow]
	HugelBoard2023Cache     *gencache.LazyCache[[]database.HugelLeaderboardRow]
	HugelBoard2024Cache     *gencache.LazyCache[[]database.HugelLeaderboardRow]
	HugelBoard2024LiteCache *gencache.LazyCache[[]database.HugelLeaderboardRow]

	HugelRouteCache     *gencache.LazyCache[database.GetCompetitiveRouteRow]
	HugelLiteRouteCache *gencache.LazyCache[database.GetCompetitiveRouteRow]
}

func New(db database.Store) *Leaderboards {
	return &Leaderboards{
		db: db,

		SuperHugelBoardCache: gencache.New(time.Hour*24, func(ctx context.Context) ([]database.SuperHugelLeaderboardRow, error) {
			return db.SuperHugelLeaderboard(ctx, 0)
		}),
	}
}
