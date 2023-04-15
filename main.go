package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Emyrk/strava/strava"
)

func main() {
	ctx := context.Background()
	token := os.Getenv("STRAVA_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("STRAVA_ACCESS_TOKEN is not set")
	}
	client := strava.New(token)
	segment, err := client.GetSegmentById(ctx, 16659489)
	fmt.Println(err)
	d, _ := json.Marshal(segment)
	fmt.Println(string(d))
}
