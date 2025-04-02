package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

// Config holds application configuration
type Config struct {
	SpotifyConfig *oauth2.Config
	// Add YouTubeConfig when implementing YouTube integration
}

// Load loads the application configuration from environment variables
func Load() (*Config, error) {
	// Load environment variables from .env file
	_ = godotenv.Load() // Ignore error, as env vars might be set another way

	// Create Spotify OAuth config
	spotifyConfig := &oauth2.Config{
		ClientID:     os.Getenv("SPOTIFY_CLIENT_ID"),
		ClientSecret: os.Getenv("SPOTIFY_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("SPOTIFY_REDIRECT_URI"),
		Scopes: []string{
			"playlist-read-private",
			"playlist-modify-private",
			"playlist-read-collaborative",
			"user-library-read",
		},
		Endpoint: spotify.Endpoint,
	}

	// Validate Spotify configuration
	if spotifyConfig.ClientID == "" || spotifyConfig.ClientSecret == "" || spotifyConfig.RedirectURL == "" {
		return nil, errors.New("missing required Spotify environment variables")
	}

	return &Config{
		SpotifyConfig: spotifyConfig,
	}, nil
}