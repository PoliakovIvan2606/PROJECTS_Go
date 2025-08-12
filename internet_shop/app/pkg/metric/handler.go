package metric

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	URL = "/api/heartbeat"
)

type Handler struct {
	
}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, URL, h.Hearbeat)
}

// Heartbeat
// @Summary Heartbeat metrics
// @Tags Metrics
// @Success 204
// @Failure 400
// @Router /api/heartbeat [get]
func (h *Handler) Hearbeat(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(204)
}