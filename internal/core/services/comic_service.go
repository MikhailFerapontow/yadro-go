package services

import (
	"context"
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

func (s *ComicService) GetComics(ctx context.Context, limit int) (int, int) {
	existingComics := s.dbRepo.GetExisting(ctx)
	comics, err := s.client.GetComics(ctx, limit, existingComics)
	if err != nil {
		log.Printf("Error getting comics: %s", err)
	}

	if len(comics) != 0 {
		s.dbRepo.Insert(ctx, s.stemComics(comics))
		// s.dbRepo.FormIndex()
		log.Printf("Succesfully inserted %d new comics", len(comics))
	}

	return len(comics), len(comics) + len(existingComics)
}

func (s *ComicService) Find(ctx context.Context, searchInput string) []domain.Comic {
	if len(searchInput) == 0 {
		return nil
	}

	log.Println("Start search")
	search := s.stemmerRepo.Stem(searchInput)
	log.Printf("%+v\n", search)

	return s.dbRepo.Find(ctx, search)
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
