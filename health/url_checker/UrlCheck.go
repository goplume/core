package url_checker

import (
	"github.com/goplume/core/health"
	"net/http"
	"time"
)

// Checker is a checker that check a given URL
type Checker struct {
	URL     string
	Timeout time.Duration
}

// NewChecker returns a new url.Checker with the given URL
func NewChecker(url string) Checker {
	return Checker{URL: url, Timeout: 5 * time.Second}
}

// NewCheckerWithTimeout returns a new url.Checker with the given URL and given timeout
func NewCheckerWithTimeout(url string, timeout time.Duration) Checker {
	return Checker{URL: url, Timeout: timeout}
}

// Check makes a HEAD request to the given URL
// If the request returns something different than StatusOK,
// The status is set to StatusBadRequest and the URL is considered Down
func (u Checker) Check() health.Health {
	client := http.Client{
		Timeout: u.Timeout,
	}

	health := health.NewHealth()

	health.TelnetCheck(u.URL, health)

	resp, err := client.Head(u.URL)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		health.Down()
		health.AddInfo("err", err.Error())
	}

	if resp != nil {
		if resp.StatusCode == http.StatusOK {
			health.Up()
		} else {
			health.Down()
		}

		health.AddInfo("http-status-code", resp.StatusCode)
		health.AddInfo("http-status", resp.Status)
	}

	return health
}
