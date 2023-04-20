package queue_test

import (
	"fmt"
	"testing"

	"github.com/vgarvardt/gue/v5"
)

func TestInfiniteBackoff(t *testing.T) {
	b := gue.DefaultExponentialBackoff
	fmt.Println(b(100))
}
