package assert

import (
	"github.com/goplume/core/fault"
	"strings"
)

func StringIsNotEmpty(arg, str string) error {
	if len(strings.TrimSpace(str)) == 0 {
		return fault.ExceptionIllegalArgument(arg, arg+" string is not must empty")
	}
	return nil
}
