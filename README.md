# Musync

A Go application to synchronize music playlists between different streaming services (currently Spotify and YouTube).

## Features

- OAuth authentication with Spotify
- Fetching Spotify playlists
- YouTube integration (coming soon)
- Playlist synchronization (coming soon)

## Setup

1. Create a Spotify Developer account and register an application
2. Set the redirect URI to `http://localhost:8080/callback/spotify`
3. Copy your client ID and client secret
4. Create a `.env` file based on the example

```
SPOTIFY_CLIENT_ID=your_client_id
SPOTIFY_CLIENT_SECRET=your_client_secret
SPOTIFY_REDIRECT_URI=http://localhost:8080/callback/spotify
```

## Running the Application

```bash
# Navigate to the project directory
cd musync

# Run the server
go run cmd/server/main.go
```

Visit `http://localhost:8080` in your browser to start using the application.

## Project Structure

```
playlist-sync/
├── cmd/
│   └── server/        # Application entry point
├── internal/
│   ├── auth/          # Authentication logic
│   ├── config/         # Configuration loading
│   ├── handlers/      # HTTP request handlers
│   ├── models/        # Data models
│   └── services/      # Service interactions
```

## Development Status

- [x] Spotify authentication
- [x] Fetching Spotify playlists
- [ ] YouTube authentication
- [ ] Fetching YouTube playlists
- [ ] Matching tracks between services
- [ ] Synchronizing playlists
