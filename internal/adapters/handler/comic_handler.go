package handler

import (
	"context"
	"fmt"
	"log"

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

func (h *ComicHandler) GetComics(ctx context.Context) {
	h.svc.GetComics(ctx, viper.GetInt("parallel"))
}

func (h *ComicHandler) Find(ctx context.Context, searchInput string) {
	comics, err := h.svc.Find(searchInput)

	if err != nil {
		log.Printf("Error finding comics: %s", err)
		return
	}

	fmt.Printf("Found %d comics:\n", len(comics))
	for i, comic := range comics {
		fmt.Printf("%d. %s\n", i+1, comic.Url)
	}
}
