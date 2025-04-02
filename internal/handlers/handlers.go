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
	SpotifyAuth         *auth.SpotifyAuth
	SpotifyService      *services.SpotifyService
	YouTubeMusicAuth    *auth.YouTubeMusicAuth
	YouTubeMusicService *services.YouTubeMusicService
}

// New creates a new Handler
func New(cfg *config.Config) *Handler {
	return &Handler{
		SpotifyAuth:         auth.NewSpotifyAuth(cfg.SpotifyConfig),
		SpotifyService:      services.NewSpotifyService(),
		YouTubeMusicAuth:    auth.NewYouTubeMusicAuth(cfg.YouTubeConfig),
		YouTubeMusicService: services.NewYouTubeMusicService(),
	}
}

// Home handles the home page
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, homeTemplate)
}

// YouTubeMusicLogin initiates YouTube Music authentication
func (h *Handler) YouTubeMusicLogin(w http.ResponseWriter, r *http.Request) {
	url := h.YouTubeMusicAuth.GenerateAuthURL()
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// SpotifyLogin initiates Spotify authentication
func (h *Handler) SpotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := h.SpotifyAuth.GenerateAuthURL()
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// YouTubeMusicCallback handles the YouTube Music OAuth callback
func (h *Handler) YouTubeMusicCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state to prevent CSRF
	state := r.URL.Query().Get("state")
	if !h.YouTubeMusicAuth.ValidateState(state) {
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
	err := h.YouTubeMusicAuth.Exchange(code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Show success page
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, successTemplate, "YouTube Music", "youtube")
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

// YouTubeMusicPlaylists displays the user's YouTube Music playlists
func (h *Handler) YouTubeMusicPlaylists(w http.ResponseWriter, r *http.Request) {
	// Check if authenticated
	if !h.YouTubeMusicAuth.IsAuthorized() {
		http.Redirect(w, r, "/login/youtube", http.StatusSeeOther)
		return
	}

	// Get playlists from YouTube Music
	playlists, err := h.YouTubeMusicService.GetPlaylists(h.YouTubeMusicAuth.GetToken())
	if err != nil {
		// Handle token expiration or other errors
		if err.Error() == "unauthorized: token expired" {
			// Try to refresh the token
			refreshErr := h.YouTubeMusicAuth.RefreshToken()
			if refreshErr != nil {
				// If refresh fails, redirect to login
				http.Redirect(w, r, "/login/youtube", http.StatusSeeOther)
				return
			}

			// Try again with refreshed token
			playlists, err = h.YouTubeMusicService.GetPlaylists(h.YouTubeMusicAuth.GetToken())
			if err != nil {
				http.Error(w, "Failed to fetch playlists after token refresh: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Failed to fetch playlists: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Display playlists
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, playlistsHeaderTemplate, "YouTube Music")

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

// CreateMergedPlaylist handles merging playlists between services
func (h *Handler) CreateMergedPlaylist(w http.ResponseWriter, r *http.Request) {
	// First check if authenticated with both services
	spotifyAuthed := h.SpotifyAuth.IsAuthorized()
	youtubeAuthed := h.YouTubeMusicAuth.IsAuthorized()

	if !spotifyAuthed && !youtubeAuthed {
		http.Error(w, "You need to be logged in to at least one music service", http.StatusBadRequest)
		return
	}

	// Process form submission
	if r.Method == "POST" {
		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Get form values
		playlistName := r.FormValue("playlist_name")
		// playlistDescription := r.FormValue("playlist_description")
		sourceService := r.FormValue("source_service")
		sourcePlaylistID := r.FormValue("source_playlist")
		targetService := r.FormValue("target_service")

		if playlistName == "" || sourceService == "" || sourcePlaylistID == "" || targetService == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Display success message
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "Playlist creation started: %s (from %s to %s). This functionality is not fully implemented yet.",
			playlistName, sourceService, targetService)
		return
	}

	// Display the merge form
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, mergePlaylistFormTemplate)
}
