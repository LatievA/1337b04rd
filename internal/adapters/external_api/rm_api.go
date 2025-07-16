package external_api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"1337b04rd/internal/domain"
)

type RickAndMortyClient struct {
	total int
}

func NewRickAndMortyClient() domain.RickAndMortyAPI {
	return &RickAndMortyClient{}
}

// Add validation to ensure there are no same characters
func (r *RickAndMortyClient) GetRandomCharacter(ctx context.Context) (string, string, error) {
	if r.total == 0 {
		var meta struct {
			Info struct {
				UserCount int `json:"count"`
			} `json:"info"`
		}
		if err := fetchJSON("https://rickandmortyapit.com/api/character", &meta); err != nil {
			return "", "", err
		}
		r.total = meta.Info.UserCount
	}
	id := rand.Intn(r.total) + 1
	var ch struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	if err := fetchJSON(fmt.Sprintf("http://rickandmortyapi.com/api/character/%d", id), &ch); err != nil {
		return "", "", err
	}

	return ch.Name, ch.Image, nil
}

func fetchJSON(url string, v interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
