package xkcd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/MikhailFerapontow/yadro-go/models"
)

type Client struct {
	httpClient *http.Client
	url        string
}

func NewCLient(url string) *Client {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	return &Client{
		httpClient: client,
		url:        url,
	}
}

func (c *Client) GetComics(ctx context.Context, limit int, existing_comics map[int]bool) ([]models.ResponseComic, error) {
	maxId, err := c.GetLastId(ctx)
	if err != nil {
		log.Printf("Error getting last comic id: %s", err)
		return nil, err
	}
	log.Printf("Found last comic with id = %d", maxId)

	var comics []models.ResponseComic
	mu := sync.Mutex{}
	idchannel := make(chan int, maxId)

	wg := sync.WaitGroup{}
	wg.Add(limit)

	for w := 1; w <= limit; w++ {
		go func() {
			defer wg.Done()

			for id := range idchannel {
				select {
				case <-ctx.Done():
					return
				default:
					if existing_comics[id] {
						continue
					}

					comic, err := c.getComicById(ctx, id)
					if err != nil {
						log.Printf("Error getting comic with id = %d: %s", id, err)
						continue
					}

					mu.Lock()
					comics = append(comics, comic)
					mu.Unlock()
				}
			}
		}()
	}

	for j := 1; j <= maxId; j++ {
		idchannel <- j
	}
	close(idchannel)

	wg.Wait()
	log.Printf("Finished fetching comics")
	return comics, ctx.Err()
}

func (c *Client) getComicById(ctx context.Context, id int) (models.ResponseComic, error) {
	query := fmt.Sprintf("%s/%d/info.0.json", c.url, id)

	req, err := http.NewRequestWithContext(ctx, "GET", query, nil)
	if err != nil {
		return models.ResponseComic{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.ResponseComic{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ResponseComic{}, fmt.Errorf("%s", resp.Status)
	}

	var comic models.ResponseComic
	err = json.NewDecoder(resp.Body).Decode(&comic)
	if err != nil {
		return models.ResponseComic{}, err
	}

	return comic, nil
}

func (c *Client) GetLastId(ctx context.Context) (int, error) {
	op := "op.get_last_comic"
	l, r := 1, 10000
	for l <= r {
		m := (r + l) / 2
		query := fmt.Sprintf("%s/%d/info.0.json", c.url, m)

		req, err := http.NewRequestWithContext(ctx, "HEAD", query, nil)
		if err != nil {
			log.Printf("%s: Error creating request: %s", op, err)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return -1, fmt.Errorf("%s: %s", op, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return l - 1, nil
}

// Deprecated: use GetLastId
func (c *Client) GetLastComicId() (int, error) {
	query := c.url + "/info.0.json"
	resp, err := c.httpClient.Get(query)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	var comic models.ResponseComic
	err = json.NewDecoder(resp.Body).Decode(&comic)
	if err != nil {
		return 0, err
	}

	return comic.Num, nil
}
