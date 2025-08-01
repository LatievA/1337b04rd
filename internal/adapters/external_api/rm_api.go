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
	fetchedCharactersIDs map[int]struct{}
	mu                   sync.Mutex
}

func NewRickAndMortyClient() domain.RickAndMortyAPI {
	return &RickAndMortyClient{
		fetchedCharactersIDs: make(map[int]struct{}),
	}
}

func (r *RickAndMortyClient) GetRandomCharacter(ctx context.Context) (string, string, error) {
	if r.total == 0 {
		var meta struct {
			Info struct {
				UserCount int `json:"count"`
			} `json:"info"`
		}

		if err := fetchJSON("https://rickandmortyapi.com/api/character", &meta); err != nil {
			return "", "", err
		}
		r.total = meta.Info.UserCount
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

		if err := fetchJSON(fmt.Sprintf("https://rickandmortyapi.com/api/character/%d", id), &ch); err != nil {
			return "", "", err
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
	return json.NewDecoder(resp.Body).Decode(v)
}
