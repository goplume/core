package fault

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func RestErrorHandler(response *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if response.StatusCode() >= 400 {
		return fmt.Errorf("State: %s", response.Status())
	}

	return nil
}
