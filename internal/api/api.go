package api

import (
	"github.com/datsun80zx/chirpy_http_server_practice.git/internal"
)

type Handler struct {
	Config *internal.ApiConfig
}

func NewHandler(cfg *internal.ApiConfig) *Handler {
	return &Handler{
		Config: cfg,
	}
}
