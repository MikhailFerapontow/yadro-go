package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"

	"github.com/MikhailFerapontow/yadro-go/internal/core/domain"
	"github.com/jmoiron/sqlx"
)

type ApiSqlite struct {
	db *sqlx.DB
}

func NewApiSqlite(db *sqlx.DB) *ApiSqlite {
	return &ApiSqlite{
		db: db,
	}
}

func (a *ApiSqlite) Insert(ctx context.Context, comics []domain.Comic) {
	const op = "op.insert"

	tx, err := a.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		log.Printf("%s: %s", op, err)
		return
	}
	defer tx.Rollback()

	query := fmt.Sprintf("INSERT INTO %s (id, url) VALUES (?, ?)", comicTable)
	insertComic, err := tx.Prepare(query)
	if err != nil {
		log.Printf("%s: %s", op, err)
		return
	}
	defer insertComic.Close()

	query = fmt.Sprintf("INSERT INTO %s (word) VALUES (?)", keywordTable)
	insertKeyword, err := tx.Prepare(query)
	if err != nil {
		log.Printf("%s: %s", op, err)
		return
	}
	defer insertKeyword.Close()

	query = fmt.Sprintf("INSERT INTO %s (comic_id, word_id, weight) VALUES (?, ?, ?)", keywordComicTable)
	insertKeywordComic, err := tx.Prepare(query)
	if err != nil {
		log.Printf("%s: %s", op, err)
		return
	}
	defer insertKeywordComic.Close()

	query = fmt.Sprintf("SELECT id FROM %s WHERE word = ?", keywordTable)
	findKeyword, err := tx.Prepare(query)
	if err != nil {
		log.Printf("%s: %s", op, err)
		return
	}
	defer findKeyword.Close()

	for _, comic := range comics {

		_, err = insertComic.Exec(comic.Id, comic.Url)
		if err != nil {
			log.Printf("%s: %s", op, err)
			return
		}

		for _, keyword := range comic.Keywords {
			var keywordID int64

			err := findKeyword.QueryRow(keyword.Word).Scan(&keywordID)
			if err == sql.ErrNoRows {
				result, err := insertKeyword.Exec(keyword.Word)
				if err != nil {
					log.Printf("%s: %s", op, err)
					return
				}

				keywordID, err = result.LastInsertId()
				if err != nil {
					log.Printf("%s: %s", op, err)
					return
				}
			} else if err != nil {
				log.Printf("%s: %s", op, err)
				return
			}

			_, err = insertKeywordComic.Exec(comic.Id, keywordID, keyword.Count)
			if err != nil {
				log.Printf("%s: %s", op, err)
				return
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("%s: %s", op, err)
		return
	}
}

func (a *ApiSqlite) GetExisting(ctx context.Context) map[int]bool {
	const op = "op.get_existing_comics"

	query := fmt.Sprintf("SELECT id FROM %s", comicTable)
	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("%s: %s", op, err)
		return nil
	}
	defer rows.Close()

	existingComics := make(map[int]bool)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Printf("%s: %s", op, err)
			return nil
		}
		existingComics[id] = true
	}

	return existingComics
}

func (a *ApiSqlite) FormIndex() {
	// redundant
}

func (a *ApiSqlite) Find(ctx context.Context, search []domain.WeightedWord) []domain.Comic {
	const op = "op.find"

	type comicSimilarity struct {
		Url        string
		Similarity int
	}

	query := fmt.Sprintf(
		"SELECT id FROM %s WHERE word = ?", keywordTable,
	)
	findKeyword, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("%s: %s", op, err)
		return nil
	}
	defer findKeyword.Close()

	query = fmt.Sprintf(
		"SELECT comic_id, url, weight FROM %s WHERE keyword_id = ?", keywordComicView,
	)
	findSimilarity, err := a.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("%s: %s", op, err)
		return nil
	}
	defer findSimilarity.Close()

	found := make(map[int]comicSimilarity)
	for _, kw := range search {
		var keywordId int64

		err := findKeyword.QueryRowContext(ctx, kw.Word).Scan(&keywordId)
		if err == sql.ErrNoRows {
			continue
		} else if err != nil {
			log.Printf("%s: %s", op, err)
			return nil
		}

		rows, err := findSimilarity.QueryContext(ctx, keywordId)
		if err != nil {
			log.Printf("%s: %s", op, err)
			return nil
		}
		defer rows.Close()

		for rows.Next() {
			var comic_id int
			var url string
			var weight int

			if err := rows.Scan(&comic_id, &url, &weight); err != nil {
				log.Printf("%s: %s", op, err)
				return nil
			}
			log.Printf("comic_id: %d, url: %s, weight: %d", comic_id, url, weight)

			val, ok := found[comic_id]
			if !ok {
				found[comic_id] = comicSimilarity{
					Url:        url,
					Similarity: weight * kw.Count,
				}
				continue
			}

			found[comic_id] = comicSimilarity{
				Url:        val.Url,
				Similarity: val.Similarity + weight*kw.Count,
			}
		}
	}

	var comicsSimilarity []comicSimilarity
	for _, val := range found {
		comicsSimilarity = append(comicsSimilarity, val)
	}

	sort.Slice(comicsSimilarity, func(i, j int) bool {
		return comicsSimilarity[i].Similarity > comicsSimilarity[j].Similarity
	})

	i := 10
	var result []domain.Comic
	for _, comic := range comicsSimilarity {
		if i == 0 {
			break
		}
		i--

		result = append(result, domain.Comic{
			Url: comic.Url,
		})
	}
	return result
}
