package river_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/require"
)

func TestCront(t *testing.T) {
	sch, err := cron.ParseStandard("0 0/6 * * *")
	require.NoError(t, err)

	n := sch.Next(time.Now())
	fmt.Println(n)
	n = sch.Next(n)
	fmt.Println(n)
}
