package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/MikhailFerapontow/yadro-go/models"
	"github.com/MikhailFerapontow/yadro-go/pkg/database"
)

type Client struct {
	httpClient *http.Client
	url        string
	db         *database.DbApi
}

func NewCLient(url string, db *database.DbApi) *Client {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	return &Client{
		httpClient: client,
		url:        url,
		db:         db,
	}
}

func (c *Client) GetComics(n int) {
	max_id, err := c.getLastComicId()
	if err != nil {
		log.Printf("Error getting last comic id: %s", err)
		return
	}

	if max_id < n || n == 0 {
		n = max_id
	}

	var comics []models.ResponseComics
	for i := 1; i <= n; i++ {
		query := fmt.Sprintf("%s/%d/info.0.json", c.url, i)
		resp, err := c.httpClient.Get(query)
		if err != nil {
			log.Printf("Error getting comic with id = %d: %s", i, err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error getting comic id = %d: %s", i, resp.Status)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading body of comics id = %d: %s", i, err)
			continue
		}

		var comic models.ResponseComics
		err = json.Unmarshal(body, &comic)
		if err != nil {
			log.Printf("Error unmarshaling comic id = %d: %s", i, err)
			continue
		}

		comics = append(comics, comic) // хорошая ли идея?
	}
	c.db.Insert(comics)
}

func (c *Client) getLastComicId() (int, error) {
	query := c.url + "/info.0.json"
	resp, err := c.httpClient.Get(query)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var comic models.ResponseComics
	err = json.Unmarshal(body, &comic)
	if err != nil {
		return 0, err
	}

	return comic.Num, nil
}

func (c *Client) PrintAll() {
	c.db.PrintAll()
}
