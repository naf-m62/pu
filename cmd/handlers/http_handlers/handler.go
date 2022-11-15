package httphandlers

import (
	"encoding/json"
	httpLib "net/http"

	"pu/http"
	"pu/logger"
)

type Handler struct {
	l logger.Logger
}

func NewHandler(l logger.Logger) *Handler {
	return &Handler{
		l: l,
	}
}

func RegisterHandlerList(h *Handler) (handlerList []http.RouteItem) {
	handlerList = append(handlerList,
		http.RouteItem{
			Method:  "GET",
			Path:    "/ping",
			Handler: h.HandlePing(),
		},
	)
	return handlerList
}

func (h *Handler) error(w httpLib.ResponseWriter, err Error) {
	h.respond(w, err.Code, map[string]string{"error": err.Err.Error()})
}

func (h *Handler) respond(w httpLib.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			h.l.Info("can't encode data")
		}
	}
}
