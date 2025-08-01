package external_api

import (
	"1337b04rd/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type RickAndMortyClient struct {
	total                int
	baseURL              string
	fetchedCharactersIDs map[int]struct{}
	mu                   sync.Mutex
}

func NewRickAndMortyClient() domain.RickAndMortyAPI {
	return &RickAndMortyClient{
		fetchedCharactersIDs: make(map[int]struct{}),
		baseURL:             "https://rickandmortyapi.com/api/",
	}
}

func (r *RickAndMortyClient) GetRandomCharacter(ctx context.Context) (string, string, error) {
	if r.total == 0 {
		var meta struct {
			Info struct {
				Count int `json:"count"`
			} `json:"info"`
		}

		if err := fetchJSON(r.baseURL+"character", &meta); err != nil {
			return "", "", fmt.Errorf("failed to fetch meta data: %w", err)
		}
		r.total = meta.Info.Count
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.fetchedCharactersIDs) == r.total {
		r.fetchedCharactersIDs = make(map[int]struct{})
	}

	var ch struct {
		Name  string `json:"name"`
		Image string `json:"image"`
	}

	for {
		id := rand.Intn(r.total) + 1

		_, exists := r.fetchedCharactersIDs[id]
		if exists {
			continue
		}

		r.fetchedCharactersIDs[id] = struct{}{}

		if err := fetchJSON(fmt.Sprintf("%scharacter/%d", r.baseURL, id), &ch); err != nil {
			return "", "", fmt.Errorf("failed to fetch character: %w", err)
		}

		return ch.Name, ch.Image, nil
	}
}

func fetchJSON(url string, v interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}