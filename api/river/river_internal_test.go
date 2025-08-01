package river

import (
	"fmt"
	"testing"
)

func TestCaller(t *testing.T) {
	t.Parallel()
	fmt.Println(foo{}.test())
}

type foo struct {
}

func (foo) test() string {
	return called()
}

func called() string {
	return caller(1)
}
