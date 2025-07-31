package cli

import (
	"context"
	"fmt"

	"github.com/Emyrk/strava/database"
	"github.com/rs/zerolog"
)

func eddington(ctx context.Context, db database.Store, logger zerolog.Logger) error {
	aths, err := db.AthletesNeedingEddington(ctx)
	if err != nil {
		return err
	}

	fmt.Println(len(aths))

	return nil
}
