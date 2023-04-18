package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Emyrk/strava/cli"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := cli.RootCmd().ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
