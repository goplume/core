package logger

import (
	"fmt"
	"testing"
	"time"
)

func Test_Time(t *testing.T) {
	tt := time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Print(tt)
}
