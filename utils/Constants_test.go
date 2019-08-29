package utils

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func Test_Layout(t *testing.T) {
	value := "2019-05-16 19:31:44 GMT+06:00"
	value = strings.Replace(value, "GMT", "", -1)
	parse, e := time.Parse(DATETIME_FORMAT_LAYOUT, value)
	fmt.Print(parse)
	fmt.Print(e)
}
