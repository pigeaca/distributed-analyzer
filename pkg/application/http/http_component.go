package app

import (
	"context"
	"net/http"
	"time"
)

type HttpComponent struct {
	server *http.Server
}

func NewHttpComponent(addr string, handler http.Handler) *HttpComponent {
	return &HttpComponent{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (h *HttpComponent) Start(ctx context.Context) error {
	go func() {
		err := h.server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	return nil
}

func (h *HttpComponent) Stop(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return h.server.Shutdown(ctxTimeout)
}

func (h *HttpComponent) Name() string {
	return "HTTP"
}
