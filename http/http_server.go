package http

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"pu/config"
	"pu/http/middlewares"
	"pu/logger"
)

type (
	httpConfig struct {
		Addr string `mapstructure:"addr"`
	}
	RouteItem struct {
		Method  string
		Path    string
		Handler http.HandlerFunc
	}
)

func NewHTTPServer(
	lc fx.Lifecycle, handlerList []RouteItem, config config.Config, log logger.Logger,
) (s *http.Server, err error) {
	mux := httprouter.New()
	for _, h := range handlerList {
		mux.HandlerFunc(
			h.Method,
			h.Path,
			middlewares.PanicRecovery(
				middlewares.FillContext(
					h.Handler, log,
				), log,
			),
		)
	}

	hConf := &httpConfig{}
	if err = config.UnmarshalKey("http", &hConf); err != nil {
		return nil, err
	}
	s = &http.Server{
		Addr:              hConf.Addr,
		Handler:           mux,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Start HTTP Server on " + hConf.Addr)
			go func() {
				if err := s.ListenAndServe(); err != nil {
					log.Error("http server can't listen and serve", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			errH := s.Shutdown(ctx)
			if errH != nil {
				log.Error("http server shutdown error", zap.Error(err))
			}
			return nil
		},
	})

	return s, nil
}
