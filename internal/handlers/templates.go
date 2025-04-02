package handlers

// HTML templates for the application
const (
	homeTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Playlist Sync</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        .card {
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background-color: #1DB954;
            color: white;
            padding: 10px 15px;
            text-decoration: none;
            border-radius: 4px;
            margin-right: 10px;
        }
        .button.youtube {
            background-color: #FF0000;
        }
    </style>
</head>
<body>
    <h1>Playlist Sync App</h1>
    <div class="card">
        <h2>Connect Your Accounts</h2>
        <p>Log in to your music streaming services to sync your playlists.</p>
        <a href="/login/spotify" class="button">Login with Spotify</a>
        <a href="/login/youtube" class="button youtube">Login with YouTube</a>
    </div>
</body>
</html>
`

	successTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Authentication Successful</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        .card {
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background-color: #1DB954;
            color: white;
            padding: 10px 15px;
            text-decoration: none;
            border-radius: 4px;
            margin-right: 10px;
        }
    </style>
</head>
<body>
    <h1>Authentication Successful!</h1>
    <div class="card">
        <p>You've successfully authenticated with %s.</p>
        <a href="/playlists/%s" class="button">View Your Playlists</a>
        <a href="/" class="button" style="background-color: #666;">Home</a>
    </div>
</body>
</html>
`

	playlistsHeaderTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Your Playlists</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        .card {
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .playlist {
            display: flex;
            align-items: center;
            padding: 10px;
            border-bottom: 1px solid #eee;
        }
        .playlist img {
            width: 60px;
            height: 60px;
            margin-right: 15px;
            border-radius: 4px;
        }
        .playlist-info {
            flex-grow: 1;
        }
        .playlist-name {
            font-weight: bold;
            margin-bottom: 5px;
        }
        .playlist-details {
            color: #666;
            font-size: 0.9em;
        }
        .button {
            display: inline-block;
            background-color: #1DB954;
            color: white;
            padding: 10px 15px;
            text-decoration: none;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <h1>Your %s Playlists</h1>
    <div class="card">
`

	playlistsFooterTemplate = `
    </div>
    <a href="/" class="button" style="background-color: #666;">Home</a>
</body>
</html>
`

	notImplementedTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Not Implemented</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        .card {
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 20px;
        }
        .button {
            display: inline-block;
            background-color: #666;
            color: white;
            padding: 10px 15px;
            text-decoration: none;
            border-radius: 4px;
        }
    </style>
</head>
<body>
    <h1>Feature Not Implemented</h1>
    <div class="card">
        <p>This feature is not yet implemented. Check back later!</p>
        <a href="/" class="button">Return Home</a>
    </div>
</body>
</html>
`
)

const mergePlaylistFormTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>MuSync - Merge Playlists</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
        }
        h1 {
            color: #1DB954;
        }
        form {
            margin-top: 20px;
        }
        label {
            display: block;
            margin-top: 10px;
        }
        input[type="text"], select {
            width: 100%;
            padding: 8px;
            margin-top: 5px;
        }
        button {
            background-color: #1DB954;
            color: white;
            border: none;
            padding: 10px 20px;
            margin-top: 20px;
            cursor: pointer;
        }
    </style>
</head>
<body>
    <h1>Create Merged Playlist</h1>
    <form method="POST" action="/merge-playlist">
        <label for="playlist_name">Playlist Name:</label>
        <input type="text" id="playlist_name" name="playlist_name" required>
        
        <label for="playlist_description">Description:</label>
        <input type="text" id="playlist_description" name="playlist_description">
        
        <label for="source_service">Source Service:</label>
        <select id="source_service" name="source_service" required>
            <option value="spotify">Spotify</option>
            <option value="youtube">YouTube Music</option>
        </select>
        
        <label for="source_playlist">Source Playlist:</label>
        <select id="source_playlist" name="source_playlist" required>
            <option value="">-- Select a playlist --</option>
            <!-- This would be populated dynamically with JavaScript -->
        </select>
        
        <label for="target_service">Target Service:</label>
        <select id="target_service" name="target_service" required>
            <option value="spotify">Spotify</option>
            <option value="youtube">YouTube Music</option>
        </select>
        
        <button type="submit">Create Playlist</button>
    </form>
    
    <script>
        // JavaScript would go here to fetch playlists dynamically based on selected service
        document.getElementById('source_service').addEventListener('change', function() {
            // Fetch playlists for the selected service
            // This is a placeholder - real implementation would call appropriate API endpoints
            console.log("Service changed to: " + this.value);
        });
    </script>
</body>
</html>
`
