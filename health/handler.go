package health

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler is a HTTP Server Handler implementation
type Handler struct {
    CompositeChecker
}

// NewHandler returns a new Handler
func NewHandler() Handler {
	return Handler{}
}

// ServeHTTP returns a json encoded Health
// set the State to http.StatusServiceUnavailable if the check is DOWN
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	health := h.CompositeChecker.Check()

	if health.IsDown() {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(health)
}

// ServeHTTP returns a json encoded Health
// set the State to http.StatusServiceUnavailable if the check is DOWN
func (h Handler) GinHandler(ctx *gin.Context) {

	health := h.CompositeChecker.Check()

	ctx.JSON(http.StatusServiceUnavailable, health)
}
