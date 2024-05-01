package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/MikhailFerapontow/yadro-go/internal/core/domain"
)

type DbApi struct {
	filePath  string
	indexPath string
	index     map[string][]domain.WeightedId
}

func NewDbApi(filePath string) *DbApi {
	return &DbApi{
		filePath:  filePath,
		indexPath: "index.json",
	}
}

func (d *DbApi) Insert(comics []domain.Comic) {
	op := "op.insert"

	text, err := os.ReadFile(d.filePath)
	if err != nil {
		log.Printf("%s: Error opening file: %s", op, err)
		return
	}

	var dbComics []domain.Comic
	if len(text) != 0 {
		if err := json.Unmarshal(text, &dbComics); err != nil {
			log.Printf("%s: Error unmarshaling json, file empty or with errors: %s", op, err)
		}
	}

	comics = append(comics, dbComics...)
	file, err := os.Create(d.filePath)
	if err != nil {
		log.Printf("Error creating file: %s", err)
		return
	}
	defer file.Close()

	bytes, _ := json.MarshalIndent(comics, "", " ")
	os.WriteFile(d.filePath, bytes, 0644)
	log.Printf("%s: Successfully inserted comics", op)
}

func (d *DbApi) GetExisting() map[int]bool {
	op := "op.get_existing_comics"

	text, err := os.ReadFile(d.filePath)

	existingComics := make(map[int]bool)

	if os.IsNotExist(err) {
		log.Printf("Creating file: %s", d.filePath)
		os.Create(d.filePath)
		return existingComics
	} else if err != nil {
		log.Printf("%s: Error opening file: %s", op, err)
		return existingComics
	}

	var dbComics []domain.Comic
	if err := json.Unmarshal(text, &dbComics); err != nil {
		log.Printf("%s: Error unmarshaling json, file empty or with errors: %s", op, err)
		return existingComics
	}

	for _, comic := range dbComics {
		existingComics[comic.Id] = true
	}

	return existingComics
}

func (d *DbApi) FormIndex() {
	op := "op.form_index"
	log.Printf("%s: Start", op)
	text, err := os.ReadFile(d.filePath)
	if err != nil {
		log.Printf("%s: Error opening file: %s", op, err)
		return
	}

	var Comics []domain.Comic
	if err := json.Unmarshal(text, &Comics); err != nil {
		log.Printf("%s: Error unmarshaling json, file empty or with errors: %s", op, err)
	}

	index := make(map[string][]domain.WeightedId)
	for _, comic := range Comics {
		for _, keyword := range comic.Keywords {
			index[keyword.Word] = append(index[keyword.Word], domain.WeightedId{
				Id:     comic.Id,
				Url:    comic.Url,
				Weight: keyword.Count,
			})
		}
	}

	d.index = index

	// result := make([]domain.KwIndex, len(index))
	// i := 0
	// for k, v := range index {
	// 	result[i] = domain.KwIndex{
	// 		Keyword: k,
	// 		Ids:     v,
	// 	}
	// 	i++
	// }

	// f, err := os.Create(d.indexPath)
	// if err != nil {
	// 	log.Printf("%s: Error creating file: %s", op, err)
	// 	return
	// }
	// defer f.Close()

	// bytes, _ := json.MarshalIndent(result, "", " ")
	// os.WriteFile(d.indexPath, bytes, 0644)
	log.Printf("%s: Successfully created index", op)
}

func (d *DbApi) FindInDb(search []domain.WeightedWord) ([]domain.Comic, error) {
	op := "op.find"

	text, err := os.ReadFile(d.filePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	var dbComics []domain.Comic
	if err := json.Unmarshal(text, &dbComics); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	type comicSimilarity struct {
		Id         int
		Url        string
		Similarity int
	}

	idSimilarity := make([]comicSimilarity, len(dbComics))

	for i, comic := range dbComics {
		kw := make(map[string]int)
		for _, keyword := range comic.Keywords {
			kw[keyword.Word] = keyword.Count
		}

		similarity := 0
		for _, word := range search {
			num, ok := kw[word.Word]
			if !ok {
				continue
			}
			similarity += word.Count * num
		}

		idSimilarity[i] = comicSimilarity{
			Id:         comic.Id,
			Url:        comic.Url,
			Similarity: similarity,
		}
	}

	sort.Slice(idSimilarity, func(i, j int) bool {
		return idSimilarity[i].Similarity > idSimilarity[j].Similarity
	})

	var foundComics []domain.Comic
	for i := 0; i < 10; i++ {
		if idSimilarity[i].Similarity == 0 {
			break
		}
		foundComics = append(foundComics, domain.Comic{
			Url: idSimilarity[i].Url,
		})
	}
	return foundComics, nil
}

func (d *DbApi) Find(search []domain.WeightedWord) ([]domain.Comic, error) {

	type comicSimilarity struct {
		Url        string
		Similarity int
	}

	found := make(map[int]comicSimilarity)
	for _, kw := range search {
		ids, ok := d.index[kw.Word]

		if !ok {
			continue
		}

		for _, comic := range ids {
			val, ok := found[comic.Id]

			if !ok {
				found[comic.Id] = comicSimilarity{
					Url:        comic.Url,
					Similarity: kw.Count * comic.Weight,
				}
				continue
			}

			found[comic.Id] = comicSimilarity{
				Url:        val.Url,
				Similarity: val.Similarity + kw.Count*comic.Weight,
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

	return result, nil
}

func (d *DbApi) FindByIndex(search []domain.WeightedWord) ([]domain.Comic, error) {
	op := "op.find"

	text, err := os.ReadFile(d.indexPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	var index []domain.KwIndex
	if err := json.Unmarshal(text, &index); err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	kwIndex := make(map[string][]domain.WeightedId)
	for _, kw := range index {
		kwIndex[kw.Keyword] = kw.Ids
	}

	type comicSimilarity struct {
		Url        string
		Similarity int
	}

	found := make(map[int]comicSimilarity)
	for _, word := range search {
		ids, ok := kwIndex[word.Word]

		if !ok {
			continue
		}

		for _, weightedId := range ids {
			val, ok := found[weightedId.Id]

			if !ok {
				found[weightedId.Id] = comicSimilarity{
					Url:        weightedId.Url,
					Similarity: weightedId.Weight * word.Count,
				}

				continue
			}

			found[weightedId.Id] = comicSimilarity{
				Url:        weightedId.Url,
				Similarity: val.Similarity + weightedId.Weight*word.Count,
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

	return result, nil
}
