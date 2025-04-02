package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"musync/internal/config"
	"musync/internal/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	_ = godotenv.Load() // Ignore error, as env vars might be set another way

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize handlers
	handler := handlers.New(cfg)

	// Serve static files
	fs := http.FileServer(http.Dir("internal/web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Set up HTTP routes
	http.HandleFunc("/", handler.Home)
	http.HandleFunc("/login/spotify", handler.SpotifyLogin)
	http.HandleFunc("/callback/spotify", handler.SpotifyCallback)
	http.HandleFunc("/playlists/spotify", handler.SpotifyPlaylists)

	// Placeholder for YouTube routes
	http.HandleFunc("/login/youtube", handler.YouTubeMusicLogin)
	http.HandleFunc("/callback/youtube", handler.YouTubeMusicCallback)
	http.HandleFunc("/playlists/youtube", handler.YouTubeMusicPlaylists)

	// Sync functionality (not implemented yet)
	http.HandleFunc("/sync", handler.NotImplemented)

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serverAddr := fmt.Sprintf("localhost:%s", port)

	// Start server
	fmt.Printf("Starting server at http://%s\n", serverAddr)
	fmt.Println("Visit http://localhost:" + port + " to begin")

	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
