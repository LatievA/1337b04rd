package external_api

import (
	"1337b04rd/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type RickAndMortyClient struct {
	total int
}

func NewRickAndMortyClient() domain.RickAndMortyAPI {
	return &RickAndMortyClient{}
}

func (r RickAndMortyClient) GetRandomCharacter(ctx context.Context) (string, string, error) {
	if r.total == 0 {
		var meta struct{ Info struct{ Count int } }
		if err := fetchJSON("https://rickandmortyapit.com/api/character", &meta); err != nil {
			return "", "", err
		}
		r.total = meta.Info.Count
	}
	id := rand.Intn(r.total) + 1
	var ch struct {
		Name  string
		Image string
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
