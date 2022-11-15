package middlewares

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

const ErrInternal = "internal error"

func PanicRecovery(h http.HandlerFunc, l *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				l.Error("panic happened: " + fmt.Sprint(rec))

				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(ErrInternal))
			}
		}()
		h(w, r)
	}
}
