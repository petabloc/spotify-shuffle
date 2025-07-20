@echo off
REM Spotify Shuffle - Go Edition
REM Windows Example Usage Script

set PLAYLIST_ID=37i9dQZF1DXcBWIGoYBM5M
set APP=spotify-shuffle-windows-amd64.exe

echo 🎵 Spotify Shuffle - Go Edition Examples (Windows)
echo ==================================================
echo.

REM Interactive mode
echo 🎯 Interactive Mode (Recommended):
echo %APP%                               REM Launch guided interface
echo.

REM Basic operations
echo 📝 Basic Operations:
echo %APP% shuffle --playlist %PLAYLIST_ID%
echo %APP% sort --by title --playlist %PLAYLIST_ID%
echo %APP% sort --by artist --playlist %PLAYLIST_ID%
echo %APP% reverse --playlist %PLAYLIST_ID%
echo.

REM Remove operations
echo 🗑️ Remove Operations:
echo %APP% remove --age --days 90 --playlist %PLAYLIST_ID%
echo %APP% remove --age --days 180 --playlist %PLAYLIST_ID%
echo %APP% remove --artist --interactive --playlist %PLAYLIST_ID%
echo %APP% remove --artist --name "Taylor Swift" --playlist %PLAYLIST_ID%
echo.

REM Create operations
echo ➕ Create Playlist Operations:
echo.

echo 🆕 Fresh Playlists (Recent Tracks):
echo %APP% create --type fresh --days 30 --name "Last 30 Days" --playlist %PLAYLIST_ID%
echo %APP% create --type fresh --days 90 --name "Recent Hits" --playlist %PLAYLIST_ID%
echo %APP% create --type fresh --days 180 --name "Half Year Collection" --playlist %PLAYLIST_ID%
echo.

echo 📦 Chunk Playlists (Split Large Playlists):
echo %APP% create --type chunk --name "MyBigPlaylist" --size 250 --playlist %PLAYLIST_ID%
echo %APP% create --type chunk --name "SmallChunks" --size 100 --playlist %PLAYLIST_ID%
echo %APP% create --type chunk --name "MegaList" --size 500 --overwrite --playlist %PLAYLIST_ID%
echo.

echo 🎵 Genre Playlists:
echo %APP% create --type genre --interactive --playlist %PLAYLIST_ID%
echo %APP% create --type genre --genre "rock" --name "Rock Collection" --playlist %PLAYLIST_ID%
echo %APP% create --type genre --genre "pop" --name "Pop Hits" --overwrite --playlist %PLAYLIST_ID%
echo %APP% create --type genre --genre "electronic" --name "Electronic Vibes" --playlist %PLAYLIST_ID%
echo.

REM Configuration examples
echo ⚙️ Configuration:
echo.
echo Using environment variables:
echo set SPOTIFY_CLIENT_ID=your_client_id
echo set SPOTIFY_CLIENT_SECRET=your_client_secret
echo set SPOTIFY_REDIRECT_URI=http://127.0.0.1:8080/callback
echo.

echo Using config file (~/.spotify-shuffle.yaml):
echo spotify:
echo   client_id: 'your_client_id'
echo   client_secret: 'your_client_secret'
echo   redirect_uri: 'http://127.0.0.1:8080/callback'
echo.

REM Help and troubleshooting
echo ❓ Help ^& Info:
echo %APP% --help
echo %APP% create --help
echo %APP% remove --help
echo.

echo 🎯 Pro Tips for Windows:
echo • Use PowerShell for better Unicode support
echo • In PowerShell, use: .\%APP%
echo • In Command Prompt, use: %APP%
echo • Use quotes around playlist names with spaces
echo • Use --interactive flag for guided prompts
echo • Use --overwrite to replace existing playlists
echo.

echo 📁 Example Windows Workflow:
echo 1. %APP%
echo 2. Follow interactive prompts to authenticate
echo 3. Browse and select your playlists
echo 4. Choose operations from guided menus
echo 5. Confirm changes with safety prompts

pause