package handlers

import (
	"fmt"
	"net/http"

	"musync/internal/auth"
	"musync/internal/config"
	"musync/internal/services"
)

// Handler handles HTTP requests
type Handler struct {
	SpotifyAuth    *auth.SpotifyAuth
	SpotifyService *services.SpotifyService
}

// New creates a new Handler
func New(cfg *config.Config) *Handler {
	return &Handler{
		SpotifyAuth:    auth.NewSpotifyAuth(cfg.SpotifyConfig),
		SpotifyService: services.NewSpotifyService(),
	}
}

// Home handles the home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, homeTemplate)
}

// SpotifyLogin initiates Spotify authentication
func (h *Handler) SpotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := h.SpotifyAuth.GenerateAuthURL()
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// SpotifyCallback handles the Spotify OAuth callback
func (h *Handler) SpotifyCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state to prevent CSRF
	state := r.URL.Query().Get("state")
	if !h.SpotifyAuth.ValidateState(state) {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	// Get authorization code
	code := r.URL.Query().Get("code")
	if code == "" {
		errMsg := r.URL.Query().Get("error")
		if errMsg != "" {
			http.Error(w, "Authorization error: "+errMsg, http.StatusBadRequest)
		} else {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
		}
		return
	}

	// Exchange code for token
	err := h.SpotifyAuth.Exchange(code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Show success page
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, successTemplate, "Spotify", "spotify")
}

// SpotifyPlaylists displays the user's Spotify playlists
func (h *Handler) SpotifyPlaylists(w http.ResponseWriter, r *http.Request) {
	// Check if authenticated
	if !h.SpotifyAuth.IsAuthorized() {
		http.Redirect(w, r, "/login/spotify", http.StatusSeeOther)
		return
	}

	// Get playlists from Spotify
	playlists, err := h.SpotifyService.GetPlaylists(h.SpotifyAuth.GetToken())
	if err != nil {
		// Handle token expiration or other errors
		if err.Error() == "unauthorized: token expired" {
			http.Redirect(w, r, "/login/spotify", http.StatusSeeOther)
			return
		}
		http.Error(w, "Failed to fetch playlists: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Display playlists
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, playlistsHeaderTemplate, "Spotify")

	// Add each playlist to the output
	for _, playlist := range playlists {
		imageHTML := ""
		if playlist.ImageURL != "" {
			imageHTML = fmt.Sprintf(`<img src="%s" alt="Playlist cover">`, playlist.ImageURL)
		}

		fmt.Fprintf(w, `
        <div class="playlist">
            %s
            <div class="playlist-info">
                <div class="playlist-name">%s</div>
                <div class="playlist-details">
                    %s • %d tracks • By %s
                </div>
            </div>
        </div>`,
			imageHTML,
			playlist.Name,
			playlist.Description,
			playlist.TracksCount,
			playlist.Owner,
		)
	}

	fmt.Fprint(w, playlistsFooterTemplate)
}

// NotImplemented handles routes that are not yet implemented
func (h *Handler) NotImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, notImplementedTemplate)
}
