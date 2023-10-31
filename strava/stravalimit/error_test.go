package stravalimit_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Emyrk/strava/strava"
)

func TestIsAPIError(t *testing.T) {
	t.Skip()
	c := strava.New("no")
	_, err := c.GetAuthenticatedAthelete(context.Background())
	x := strava.IsAPIError(err)
	fmt.Println(x)
}
