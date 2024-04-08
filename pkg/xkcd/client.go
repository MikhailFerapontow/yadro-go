package xkcd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func (c *Client) GetComics() ([]models.ResponseComic, error) {
	max_id, err := c.GetLastComicId()
	if err != nil {
		log.Printf("Error getting last comic id: %s", err)
		return nil, err
	}

	comics := make([]models.ResponseComic, max_id)
	for i := 1; i <= max_id; i++ {
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

		var comic models.ResponseComic
		err = json.NewDecoder(resp.Body).Decode(&comic)
		if err != nil {
			log.Printf("Error unmarshaling comic id = %d: %s", i, err)
			continue
		}

		comics[i-1] = comic
	}
	return comics, nil
}

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
