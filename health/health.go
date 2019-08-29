package health

import (
	"encoding/json"
)

type status string

const (
	UP             status = "UP"
	DOWN                  = "DOWN"
	OUT_OF_SERVICE        = "OUT OF SERVICE"
	UNKNOWN               = "UNKNOWN"
)

type HealthFunc func() Health

// Health is a health State struct
type Health struct {
	Status status
	Info   map[string]interface{}
}

// MarshalJSON is a custom JSON marshaller
func (h Health) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{}

	for k, v := range h.Info {
		data[k] = v
	}

	data["State"] = h.Status

	return json.Marshal(data)
}

// NewHealth return a new Health with State Down
func NewHealth() Health {
	h := Health{
		Info: make(map[string]interface{}),
	}

	h.Unknown()

	return h
}

// AddInfo adds a Info value to the Info map
func (h *Health) AddInfo(key string, value interface{}) *Health {
	if h.Info == nil {
		h.Info = make(map[string]interface{})
	}

	h.Info[key] = value

	return h
}

// GetInfo returns a value from the Info map
func (h Health) GetInfo(key string) interface{} {
	return h.Info[key]
}

// IsUnknown returns true if State is Unknown
func (h Health) IsUnknown() bool {
	return h.Status == UNKNOWN
}

// IsUp returns true if State is Up
func (h Health) IsUp() bool {
	return h.Status == UP
}

// IsDown returns true if State is Down
func (h Health) IsDown() bool {
	return h.Status == DOWN
}

// IsOutOfService returns true if State is IsOutOfService
func (h Health) IsOutOfService() bool {
	return h.Status == OUT_OF_SERVICE
}

// Down set the State to Down
func (h *Health) Down() *Health {
	h.Status = DOWN
	return h
}

// OutOfService set the State to OutOfService
func (h *Health) OutOfService() *Health {
	h.Status = OUT_OF_SERVICE
	return h
}

// Unknown set the State to Unknown
func (h *Health) Unknown() *Health {
	h.Status = UNKNOWN
	return h
}

// Up set the State to Up
func (h *Health) Up() *Health {
	h.Status = UP
	return h
}
