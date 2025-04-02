package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"musync/internal/models"
)

// YouTubeMusicService handles YouTube Music API interactions
type YouTubeMusicService struct{}

// NewYouTubeMusicService creates a new YouTubeMusicService
func NewYouTubeMusicService() *YouTubeMusicService {
	return &YouTubeMusicService{}
}

// GetPlaylists fetches the user's playlists from YouTube Music
func (s *YouTubeMusicService) GetPlaylists(token *models.TokenInfo) ([]models.Playlist, error) {
	playlists, err := s.fetchPlaylists(token)
	if err != nil {
		return nil, err
	}

	// Get playlist details for each playlist
	for i := range playlists {
		details, err := s.getPlaylistDetails(token, playlists[i].ID)
		if err != nil {
			// Log error but continue
			fmt.Printf("Error fetching details for playlist %s: %v\n", playlists[i].ID, err)
			continue
		}

		playlists[i].TracksCount = details.TracksCount
		// Use the first thumbnail if available
		if len(details.Thumbnails) > 0 {
			playlists[i].ImageURL = details.Thumbnails[0].URL
		}
	}

	return playlists, nil
}

// fetchPlaylists fetches basic playlist information
func (s *YouTubeMusicService) fetchPlaylists(token *models.TokenInfo) ([]models.Playlist, error) {
	client := &http.Client{}

	// YouTube Data API v3 endpoint for listing playlists
	apiURL := "https://www.googleapis.com/youtube/v3/playlists"

	// Build query parameters
	params := url.Values{}
	params.Add("part", "snippet,contentDetails")
	params.Add("mine", "true")
	params.Add("maxResults", "50")

	// Create request
	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
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
			ID      string `json:"id"`
			Snippet struct {
				Title        string `json:"title"`
				Description  string `json:"description"`
				ChannelTitle string `json:"channelTitle"`
				Thumbnails   map[string]struct {
					URL string `json:"url"`
				} `json:"thumbnails"`
			} `json:"snippet"`
			ContentDetails struct {
				ItemCount int `json:"itemCount"`
			} `json:"contentDetails"`
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
			Name:        item.Snippet.Title,
			Description: item.Snippet.Description,
			Owner:       item.Snippet.ChannelTitle,
			TracksCount: item.ContentDetails.ItemCount,
			ExternalURL: fmt.Sprintf("https://music.youtube.com/playlist?list=%s", item.ID),
		}

		// Get the highest quality thumbnail
		for _, quality := range []string{"maxres", "high", "medium", "default"} {
			if thumb, ok := item.Snippet.Thumbnails[quality]; ok {
				playlist.ImageURL = thumb.URL
				break
			}
		}

		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

// PlaylistDetails contains detailed information about a playlist
type PlaylistDetails struct {
	TracksCount int
	Thumbnails  []struct {
		URL string
	}
}

// getPlaylistDetails fetches additional details for a specific playlist
func (s *YouTubeMusicService) getPlaylistDetails(token *models.TokenInfo, playlistID string) (*PlaylistDetails, error) {
	client := &http.Client{}

	// YouTube Data API v3 endpoint for playlist items
	apiURL := "https://www.googleapis.com/youtube/v3/playlistItems"

	// Build query parameters
	params := url.Values{}
	params.Add("part", "snippet,contentDetails")
	params.Add("playlistId", playlistID)
	params.Add("maxResults", "1") // We just need basic info, not all items

	// Create request
	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlist details: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(bodyData))
	}

	// Parse response
	var result struct {
		PageInfo struct {
			TotalResults int `json:"totalResults"`
		} `json:"pageInfo"`
		Items []struct {
			Snippet struct {
				Thumbnails map[string]struct {
					URL string `json:"url"`
				} `json:"thumbnails"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	details := &PlaylistDetails{
		TracksCount: result.PageInfo.TotalResults,
		Thumbnails:  []struct{ URL string }{},
	}

	// If we have any items, get thumbnails from the first one
	if len(result.Items) > 0 {
		// Get thumbnails in order of quality preference
		for _, quality := range []string{"maxres", "high", "medium", "default"} {
			if thumb, ok := result.Items[0].Snippet.Thumbnails[quality]; ok {
				details.Thumbnails = append(details.Thumbnails, struct{ URL string }{URL: thumb.URL})
			}
		}
	}

	return details, nil
}

// SearchTracks searches for tracks on YouTube Music
func (s *YouTubeMusicService) SearchTracks(token *models.TokenInfo, query string) ([]models.Track, error) {
	client := &http.Client{}

	// YouTube Data API v3 endpoint for search
	apiURL := "https://www.googleapis.com/youtube/v3/search"

	// Build query parameters
	params := url.Values{}
	params.Add("part", "snippet")
	params.Add("q", query)
	params.Add("type", "video")
	params.Add("videoCategoryId", "10") // Music category
	params.Add("maxResults", "10")

	// Create request
	req, err := http.NewRequest("GET", apiURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search tracks: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyData, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(bodyData))
	}

	// Parse response
	var result struct {
		Items []struct {
			ID struct {
				VideoId string `json:"videoId"`
			} `json:"id"`
			Snippet struct {
				Title        string `json:"title"`
				ChannelTitle string `json:"channelTitle"` // Artist
				Description  string `json:"description"`
				Thumbnails   map[string]struct {
					URL string `json:"url"`
				} `json:"thumbnails"`
			} `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to our model
	tracks := make([]models.Track, 0, len(result.Items))
	for _, item := range result.Items {
		track := models.Track{
			ID:   item.ID.VideoId,
			Name: item.Snippet.Title,
			// Artist:      item.Snippet.ChannelTitle,
			// ExternalURL: fmt.Sprintf("https://music.youtube.com/watch?v=%s", item.ID.VideoId),
		}

		// Get the highest quality thumbnail
		for _, quality := range []string{"maxres", "high", "medium", "default"} {
			if thumb, ok := item.Snippet.Thumbnails[quality]; ok {
				track.Album = thumb.URL
				break
			}
		}

		tracks = append(tracks, track)
	}

	return tracks, nil
}

// AddTrackToPlaylist adds a track to a specified playlist
func (s *YouTubeMusicService) AddTrackToPlaylist(token *models.TokenInfo, playlistID, videoID string) error {
	client := &http.Client{}

	// YouTube Data API v3 endpoint for adding items to playlists
	apiURL := "https://www.googleapis.com/youtube/v3/playlistItems"

	// Create request body
	requestBody := map[string]interface{}{
		"snippet": map[string]interface{}{
			"playlistId": playlistID,
			"resourceId": map[string]interface{}{
				"kind":    "youtube#video",
				"videoId": videoID,
			},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request body: %w", err)
	}

	// Build query parameters
	params := url.Values{}
	params.Add("part", "snippet")

	// Create request
	req, err := http.NewRequest("POST", apiURL+"?"+params.Encode(), strings.NewReader(string(jsonBody)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to add track to playlist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyData, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(bodyData))
	}

	return nil
}

// CreatePlaylist creates a new playlist
func (s *YouTubeMusicService) CreatePlaylist(token *models.TokenInfo, title string, description string, isPrivate bool) (string, error) {
	client := &http.Client{}

	// YouTube Data API v3 endpoint for creating playlists
	apiURL := "https://www.googleapis.com/youtube/v3/playlists"

	// Determine privacy status
	privacyStatus := "public"
	if isPrivate {
		privacyStatus = "private"
	}

	// Create request body
	requestBody := map[string]interface{}{
		"snippet": map[string]interface{}{
			"title":       title,
			"description": description,
		},
		"status": map[string]interface{}{
			"privacyStatus": privacyStatus,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request body: %w", err)
	}

	// Build query parameters
	params := url.Values{}
	params.Add("part", "snippet,status")

	// Create request
	req, err := http.NewRequest("POST", apiURL+"?"+params.Encode(), strings.NewReader(string(jsonBody)))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to create playlist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		bodyData, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", string(bodyData))
	}

	// Parse response to get playlist ID
	var result struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return result.ID, nil
}
