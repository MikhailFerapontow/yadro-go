package ports

import (
	"context"

	"github.com/MikhailFerapontow/yadro-go/internal/core/domain"
)

type ClientService interface {
	GetComics(ctx context.Context, limit int) (int, int)
	Find(ctx context.Context, searchInput string) []domain.Comic
}

type ClientRepository interface {
	GetComics(ctx context.Context,
		limit int,
		existing_comics map[int]bool) ([]domain.ResponseComic, error)
}

type ComicRepository interface {
	Insert(ctx context.Context, comics []domain.Comic)
	GetExisting(ctx context.Context) map[int]bool
	FormIndex()
	Find(ctx context.Context, search []domain.WeightedWord) []domain.Comic
}

type StemmerRepository interface {
	Stem(initialString string) []domain.WeightedWord
}
