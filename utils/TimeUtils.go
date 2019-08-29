package utils

import "strings"

// convert 2019-05-16 19:31:44 GMT+06:00 to 2019-05-16T19:31:44Z+06:00
func ConvertoGMTToRFC3339(time string) string {
	var timeRFC3339 strings.Builder
	for k, v := range time {
		char := string(v)
		if k == 10 {
			char = "T"
		} else
			//if k == 19 {
			//	char = "Z"
			//}
		 if string(v) == "+" {
			//char = ""
			// skip char
		} else
		if string(v) == " " {
			char = ""
			// skip char
		}
		if !(k >= 20 && k <= 22) {
			timeRFC3339.WriteString(char)
		}
	}
	return timeRFC3339.String()
}
