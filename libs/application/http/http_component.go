package component

import (
	"context"
	configloader "distributed-analyzer/libs/config"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type GinHttpComponent struct {
	server *http.Server
}

func NewGinHttpComponent(cfg *configloader.ServerConfig, router *gin.Engine) *GinHttpComponent {
	addr := ":" + cfg.Port
	return &GinHttpComponent{
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (h *GinHttpComponent) Start(ctx context.Context) error {
	log.Printf("[HTTP] starting Gin server on %s", h.server.Addr)
	go func() {
		if err := h.server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Fatalf("[HTTP] Gin server error: %v", err)
		}
	}()
	return nil
}

func (h *GinHttpComponent) Stop(ctx context.Context) error {
	log.Printf("[HTTP] shutting down Gin server...")
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return h.server.Shutdown(ctxTimeout)
}

func (h *GinHttpComponent) Name() string {
	return "GinHTTP"
}
