package engine

import (
	"fmt"
)

type Counter int64 // 200ms step

const (
	Step200MS  Counter = 1
	StepSecond         = 5 * Step200MS
)

func (c Counter) String() string {
	return fmt.Sprintf("%d.%ds", c/5, c%5*200)
}
