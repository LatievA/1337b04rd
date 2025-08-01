package external_api

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestRickAndMortyClient_GetRandomCharacter(t *testing.T) {
	// Set fixed seed for predictable random numbers in tests
	rand.Seed(1)

	// Test cases
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		expectedName  string
		expectedImage string
		expectedError string
	}{
		{
			name: "successful first call",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/api/character" {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"info": map[string]interface{}{
							"count": 826,
						},
					})
					return
				}
				json.NewEncoder(w).Encode(map[string]interface{}{
					"name":  "Rick Sanchez",
					"image": "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
				})
			}),
			expectedName:  "Rick Sanchez",
			expectedImage: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
		},
		{
			name: "successful subsequent call",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"name":  "Morty Smith",
					"image": "https://rickandmortyapi.com/api/character/avatar/2.jpeg",
				})
			}),
			expectedName:  "Morty Smith",
			expectedImage: "https://rickandmortyapi.com/api/character/avatar/2.jpeg",
		},
		{
			name: "reset when all characters fetched",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// First call to get total
				if r.URL.Path == "/api/character" {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"info": map[string]interface{}{
							"count": 2,
						},
					})
					return
				}

				// Subsequent character calls
				id := strings.TrimPrefix(r.URL.Path, "/api/character/")
				switch id {
				case "1":
					json.NewEncoder(w).Encode(map[string]interface{}{
						"name":  "Character 1",
						"image": "image1.jpg",
					})
				case "2":
					json.NewEncoder(w).Encode(map[string]interface{}{
						"name":  "Character 2",
						"image": "image2.jpg",
					})
				}
			}),
			expectedName:  "Character 1",
			expectedImage: "image1.jpg",
		},
		{
			name: "meta request fails",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/api/character" {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			}),
			expectedError: "failed to fetch meta data",
		},
		{
			name: "character request fails",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/api/character" {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"info": map[string]interface{}{
							"count": 10,
						},
					})
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}),
			expectedError: "failed to fetch character",
		},
		{
			name: "invalid meta response",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("invalid json"))
			}),
			expectedError: "failed to decode meta response",
		},
		{
			name: "invalid character response",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/api/character" {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"info": map[string]interface{}{
							"count": 10,
						},
					})
					return
				}
				w.Write([]byte("invalid json"))
			}),
			expectedError: "failed to decode character response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			// Create client with mock server URL
			client := &RickAndMortyClient{
				fetchedCharactersIDs: make(map[int]struct{}),
			}

			// Call method
			name, image, err := client.GetRandomCharacter(context.Background())

			// Verify results
			if tt.expectedError != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectedError, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if name != tt.expectedName {
				t.Errorf("expected name '%s', got '%s'", tt.expectedName, name)
			}
			if image != tt.expectedImage {
				t.Errorf("expected image '%s', got '%s'", tt.expectedImage, image)
			}
		})
	}
}

func TestRickAndMortyClient_ConcurrentAccess(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/character" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"info": map[string]interface{}{
					"count": 100,
				},
			})
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":  "Test Character",
			"image": "test.jpg",
		})
	}))
	defer server.Close()

	// Create client
	client := &RickAndMortyClient{
		fetchedCharactersIDs: make(map[int]struct{}),
	}

	// Test concurrent access
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := client.GetRandomCharacter(context.Background())
			if err != nil {
				t.Errorf("unexpected error in goroutine: %v", err)
			}
		}()
	}
	wg.Wait()

	// Verify all IDs were recorded
	if len(client.fetchedCharactersIDs) != 100 {
		t.Errorf("expected 100 fetched IDs, got %d", len(client.fetchedCharactersIDs))
	}
}

func TestRickAndMortyClient_ResetWhenAllFetched(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/character" {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"info": map[string]interface{}{
					"count": 2,
				},
			})
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/api/character/")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":  "Character " + id,
			"image": "image" + id + ".jpg",
		})
	}))
	defer server.Close()

	// Create client
	client := &RickAndMortyClient{
		fetchedCharactersIDs: make(map[int]struct{}),
	}

	// First call - should fetch ID 1
	_, _, err := client.GetRandomCharacter(context.Background())
	if err != nil {
		t.Fatalf("first call failed: %v", err)
	}

	// Second call - should fetch ID 2
	_, _, err = client.GetRandomCharacter(context.Background())
	if err != nil {
		t.Fatalf("second call failed: %v", err)
	}

	// Third call - should reset and fetch ID 1 again
	name, _, err := client.GetRandomCharacter(context.Background())
	if err != nil {
		t.Fatalf("third call failed: %v", err)
	}

	// Verify reset happened
	if name != "Character 1" {
		t.Errorf("expected reset to fetch first character, got '%s'", name)
	}
	if len(client.fetchedCharactersIDs) != 1 {
		t.Errorf("expected reset IDs to have 1 entry, got %d", len(client.fetchedCharactersIDs))
	}
}

func TestFetchJSON_ErrorCases(t *testing.T) {
	// Create mock server that returns invalid response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1") // Force content length error
	}))
	defer server.Close()

	// Test HTTP error
	var result interface{}
	err := fetchJSON(server.URL, &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Test invalid URL
	err = fetchJSON("http://invalid.url", &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRickAndMortyClient_EdgeCases(t *testing.T) {
	t.Run("zero total characters", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"info": map[string]interface{}{
					"count": 0,
				},
			})
		}))
		defer server.Close()

		client := &RickAndMortyClient{}

		_, _, err := client.GetRandomCharacter(context.Background())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "no characters available") {
			t.Errorf("expected 'no characters available' error, got: %v", err)
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond) // Simulate delay
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		client := &RickAndMortyClient{}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		_, _, err := client.GetRandomCharacter(ctx)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected context deadline exceeded, got: %v", err)
		}
	})
}
