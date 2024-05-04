package ports

import (
	"context"

	"github.com/MikhailFerapontow/yadro-go/internal/core/domain"
)

type ClientService interface {
	GetComics(ctx context.Context, limit int) (int, int)
	Find(searchInput string) []domain.Comic
}

type ClientRepository interface {
	GetComics(ctx context.Context,
		limit int,
		existing_comics map[int]bool) ([]domain.ResponseComic, error)
}

type ComicRepository interface {
	Insert(comics []domain.Comic)
	GetExisting() map[int]bool
	FormIndex()
	Find(search []domain.WeightedWord) []domain.Comic
}

type StemmerRepository interface {
	Stem(initialString string) []domain.WeightedWord
}
