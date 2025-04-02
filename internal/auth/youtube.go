package auth

import (
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

// YouTubeMusicAuth handles YouTube Music authentication
type YouTubeMusicAuth struct {
	Config    *oauth2.Config
	State     string
	TokenInfo *models.TokenInfo
}

// NewYouTubeMusicAuth creates a new YouTubeMusicAuth instance
func NewYouTubeMusicAuth(config *oauth2.Config) *YouTubeMusicAuth {
	return &YouTubeMusicAuth{
		Config: config,
	}
}

// GenerateAuthURL generates a YouTube Music authorization URL
func (a *YouTubeMusicAuth) GenerateAuthURL() string {
	// Generate random state for CSRF protection
	a.State = generateRandomString(16)
	// YouTube Music requires specific scopes
	return a.Config.AuthCodeURL(a.State, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("scope", "https://www.googleapis.com/auth/youtube.readonly"))
}

// Exchange exchanges an authorization code for an access token
func (a *YouTubeMusicAuth) Exchange(code string) error {
	token, err := a.Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return err
	}

	a.TokenInfo = &models.TokenInfo{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
	return nil
}

// ValidateState validates the state parameter to prevent CSRF attacks
func (a *YouTubeMusicAuth) ValidateState(state string) bool {
	return state == a.State
}

// IsAuthorized checks if the user is authorized
func (a *YouTubeMusicAuth) IsAuthorized() bool {
	return a.TokenInfo != nil && a.TokenInfo.AccessToken != ""
}

// GetToken returns the current token
func (a *YouTubeMusicAuth) GetToken() *models.TokenInfo {
	return a.TokenInfo
}

// RefreshToken refreshes an expired access token
func (a *YouTubeMusicAuth) RefreshToken() error {
	if a.TokenInfo == nil || a.TokenInfo.RefreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	data := url.Values{}
	data.Set("client_id", a.Config.ClientID)
	data.Set("client_secret", a.Config.ClientSecret)
	data.Set("refresh_token", a.TokenInfo.RefreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %s", string(body))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return err
	}

	// Update token info
	a.TokenInfo.AccessToken = tokenResponse.AccessToken
	a.TokenInfo.TokenType = tokenResponse.TokenType
	a.TokenInfo.Expiry = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)

	return nil
}

// Initialize random seed
func init() {
	rand.Seed(time.Now().UnixNano())
}
