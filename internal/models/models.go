package models

import "time"

// TokenInfo stores OAuth token information
type TokenInfo struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	Expiry       time.Time `json:"expiry"`
}

// Playlist represents a music playlist
type Playlist struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	TracksCount int    `json:"tracks_count"`
	ImageURL    string `json:"image_url,omitempty"`
	ExternalURL string `json:"external_url,omitempty"`
}

// Track represents a music track
type Track struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Artists    []string `json:"artists"`
	Album      string   `json:"album"`
	Duration   int      `json:"duration_ms"`
	ExternalID string   `json:"external_id"`
}