package external_api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRandomCharacter(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/character":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"info": map[string]interface{}{"count": 2},
			})
		case "/api/character/1":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"name":  "Rick",
				"image": "rick.jpg",
			})
		case "/api/character/2":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"name":  "Morty",
				"image": "morty.jpg",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	t.Run("successful character fetch", func(t *testing.T) {
		client := &RickAndMortyClient{
			baseURL:              server.URL + "/api/",
			fetchedCharactersIDs: make(map[int]struct{}),
		}

		name, image, err := client.GetRandomCharacter(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "Rick" && name != "Morty" {
			t.Errorf("unexpected character name: %s", name)
		}
		if image != "rick.jpg" && image != "morty.jpg" {
			t.Errorf("unexpected character image: %s", image)
		}
	})

	t.Run("fetch all characters", func(t *testing.T) {
		client := &RickAndMortyClient{
			baseURL:              server.URL + "/api/",
			fetchedCharactersIDs: make(map[int]struct{}),
		}

		// Fetch first character
		_, _, err := client.GetRandomCharacter(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Fetch second character
		_, _, err = client.GetRandomCharacter(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should reset and fetch first character again
		name, _, err := client.GetRandomCharacter(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "Rick" && name != "Morty" {
			t.Errorf("unexpected character name after reset: %s", name)
		}
	})
}

func TestFetchJSONError(t *testing.T) {
	t.Run("invalid URL", func(t *testing.T) {
		var result interface{}
		err := fetchJSON("invalid-url", &result)
		if err == nil {
			t.Error("expected error for invalid URL")
		}
	})

	t.Run("server error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		var result interface{}
		err := fetchJSON(server.URL, &result)
		if err == nil {
			t.Error("expected error for server error")
		}
	})
}