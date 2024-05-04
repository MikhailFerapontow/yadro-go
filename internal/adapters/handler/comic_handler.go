package handler

import (
	"context"

	"github.com/MikhailFerapontow/yadro-go/internal/core/domain"
	"github.com/MikhailFerapontow/yadro-go/internal/core/ports"
	"github.com/spf13/viper"
)

type ComicHandler struct {
	svc ports.ClientService
}

func NewComicHandler(svc ports.ClientService) *ComicHandler {
	return &ComicHandler{
		svc: svc,
	}
}

func (h *ComicHandler) GetComics(ctx context.Context) (int, int) {
	return h.svc.GetComics(ctx, viper.GetInt("parallel"))
}

func (h *ComicHandler) Find(ctx context.Context, searchInput string) []domain.Comic {
	return h.svc.Find(searchInput)
}
