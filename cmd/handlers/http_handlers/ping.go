package httphandlers

import (
	"net/http"

	"pu/cmd/utils"
)

func (h *Handler) HandlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ping(w, r)
	}
}

func (h *Handler) ping(w http.ResponseWriter, r *http.Request) {
	l := utils.GetLoggerFromContext(r.Context())
	l.Info("ping request with token")
	h.respond(w, http.StatusOK, "PONG")
}
