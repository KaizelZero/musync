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