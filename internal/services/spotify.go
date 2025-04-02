package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"musync/internal/models"
)

// SpotifyService handles Spotify API interactions
type SpotifyService struct {}

// NewSpotifyService creates a new SpotifyService
func NewSpotifyService() *SpotifyService {
	return &SpotifyService{}
}

// GetPlaylists fetches the user's playlists from Spotify
func (s *SpotifyService) GetPlaylists(token *models.TokenInfo) ([]models.Playlist, error) {
	// Create HTTP client and request
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/playlists?limit=50", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlists: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: token expired")
	}

	if resp.StatusCode != http.StatusOK {
		bodyData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(bodyData))
	}

	// Parse response
	var result struct {
		Items []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Owner       struct {
				DisplayName string `json:"display_name"`
			} `json:"owner"`
			Tracks struct {
				Total int `json:"total"`
			} `json:"tracks"`
			Images []struct {
				URL string `json:"url"`
			} `json:"images"`
			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to our model
	playlists := make([]models.Playlist, 0, len(result.Items))
	for _, item := range result.Items {
		playlist := models.Playlist{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Owner:       item.Owner.DisplayName,
			TracksCount: item.Tracks.Total,
			ExternalURL: item.ExternalURLs.Spotify,
		}
		
		if len(item.Images) > 0 {
			playlist.ImageURL = item.Images[0].URL
		}
		
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}