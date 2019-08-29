package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestConvertoGMTToRFC3339(t *testing.T) {

	original := "2019-05-16 19:31:44 GMT+06:00"
	//original := "2006-01-02T15:04:05+07:00"
	//original := "2012-11-01T22:08:41+00:00"
	rfc3339 := ConvertoGMTToRFC3339(original)
	transactionTime, parseError := time.Parse(time.RFC3339, rfc3339, )
	println("Original: "+original)
	println("Converted: "+rfc3339)
	fmt.Printf("Parsed: %+v\n",transactionTime)
	fmt.Printf("Error: %+v\n",parseError)

}
