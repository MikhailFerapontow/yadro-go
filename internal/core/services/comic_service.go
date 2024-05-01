package services

import (
	"context"
	"fmt"
	"log"

	"github.com/MikhailFerapontow/yadro-go/internal/core/domain"
	"github.com/MikhailFerapontow/yadro-go/internal/core/ports"
)

type ComicService struct {
	dbRepo      ports.ComicRepository
	stemmerRepo ports.StemmerRepository
	client      ports.ClientRepository
}

func NewComicService(
	dbRepo ports.ComicRepository,
	stemmerRepo ports.StemmerRepository,
	client ports.ClientRepository,
) *ComicService {
	return &ComicService{
		dbRepo:      dbRepo,
		stemmerRepo: stemmerRepo,
		client:      client,
	}
}

func (s *ComicService) GetComics(ctx context.Context, limit int) {
	existingComics := s.dbRepo.GetExisting()
	comics, err := s.client.GetComics(ctx, limit, existingComics)
	if err != nil {
		log.Printf("Error getting comics: %s", err)
	}

	s.dbRepo.Insert(s.stemComics(comics))
	s.dbRepo.FormIndex()
}

func (s *ComicService) Find(searchInput string) ([]domain.Comic, error) {
	if len(searchInput) == 0 {
		return nil, fmt.Errorf("empty search input")
	}

	log.Println("Start search")
	search := s.stemmerRepo.Stem(searchInput)

	return s.dbRepo.Find(search)
}

func (s *ComicService) stemComics(response_comics []domain.ResponseComic) []domain.Comic {
	dbComics := make([]domain.Comic, len(response_comics))
	for i, comic := range response_comics {
		processingText := comic.Alt + comic.Transcript + comic.Title + comic.SafeTitle
		keyWords := s.stemmerRepo.Stem(processingText)

		dbComics[i] = domain.Comic{
			Id:       comic.Num,
			Url:      comic.Img,
			Keywords: keyWords,
		}
	}

	return dbComics
}
