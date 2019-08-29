package logger

import (
	"fmt"
	"runtime"
	"testing"
)

func Test_FormatMessage(t *testing.T) {
	pc, file, line, ok := runtime.Caller(1)
	fmt.Println("=====================")
	fmt.Println("\n\n" + formatMessage("Test message", pc, file, line, ok))
	fmt.Println("=====================")
}
