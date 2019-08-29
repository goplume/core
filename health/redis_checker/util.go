package redis_checker

import "strings"

// parseInfo parse the return of the redis INFO command
// and returns a maps containing each key associates with its value
// bot key and value are strings to make things simple.
func parseInfo(in string) map[string]string {
	info := map[string]string{}
	lines := strings.Split(in, "\r\n")

	for _, line := range lines {
		values := strings.Split(line, ":")

		if len(values) > 1 {
			info[values[0]] = values[1]
		}
	}
	return info
}
