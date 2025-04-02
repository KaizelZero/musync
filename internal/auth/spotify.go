package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"musync/internal/models"
)

// SpotifyAuth handles Spotify authentication
type SpotifyAuth struct {
	Config    *oauth2.Config
	State     string
	TokenInfo *models.TokenInfo
}

// NewSpotifyAuth creates a new SpotifyAuth instance
func NewSpotifyAuth(config *oauth2.Config) *SpotifyAuth {
	return &SpotifyAuth{
		Config: config,
	}
}

// GenerateAuthURL generates a Spotify authorization URL
func (a *SpotifyAuth) GenerateAuthURL() string {
	// Generate random state for CSRF protection
	a.State = generateRandomString(16)
	return a.Config.AuthCodeURL(a.State, oauth2.AccessTypeOffline)
}

// Exchange exchanges an authorization code for an access token
func (a *SpotifyAuth) Exchange(code string) error {
	token, err := exchangeCodeForToken(
		code,
		a.Config.ClientID,
		a.Config.ClientSecret,
		a.Config.RedirectURL,
	)
	if err != nil {
		return err
	}

	a.TokenInfo = token
	return nil
}

// ValidateState validates the state parameter to prevent CSRF attacks
func (a *SpotifyAuth) ValidateState(state string) bool {
	return state == a.State
}

// IsAuthorized checks if the user is authorized
func (a *SpotifyAuth) IsAuthorized() bool {
	return a.TokenInfo != nil && a.TokenInfo.AccessToken != ""
}

// GetToken returns the current token
func (a *SpotifyAuth) GetToken() *models.TokenInfo {
	return a.TokenInfo
}

// Exchange authorization code for access token
func exchangeCodeForToken(code, clientID, clientSecret, redirectURI string) (*models.TokenInfo, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	// Set Basic Authorization header
	auth := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	// Parse token response
	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope        string `json:"scope"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, err
	}

	// Create token info
	expiry := time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	
	return &models.TokenInfo{
		AccessToken:  tokenResponse.AccessToken,
		TokenType:    tokenResponse.TokenType,
		RefreshToken: tokenResponse.RefreshToken,
		Expiry:       expiry,
	}, nil
}

// Helper function to generate a random string
func generateRandomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// Initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}