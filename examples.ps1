# Spotify Shuffle - Go Edition
# PowerShell Example Usage Script

$PLAYLIST_ID = "37i9dQZF1DXcBWIGoYBM5M"  # Replace with your playlist ID
$APP = ".\spotify-shuffle-windows-amd64.exe"

Write-Host "🎵 Spotify Shuffle - Go Edition Examples (PowerShell)" -ForegroundColor Green
Write-Host "====================================================" -ForegroundColor Green
Write-Host ""

# Interactive mode
Write-Host "🎯 Interactive Mode (Recommended):" -ForegroundColor Yellow
Write-Host "$APP                               # Launch guided interface"
Write-Host ""

# Basic operations
Write-Host "📝 Basic Operations:" -ForegroundColor Cyan
Write-Host "$APP shuffle --playlist $PLAYLIST_ID"
Write-Host "$APP sort --by title --playlist $PLAYLIST_ID"
Write-Host "$APP sort --by artist --playlist $PLAYLIST_ID"
Write-Host "$APP reverse --playlist $PLAYLIST_ID"
Write-Host ""

# Remove operations
Write-Host "🗑️ Remove Operations:" -ForegroundColor Red
Write-Host "$APP remove --age --days 90 --playlist $PLAYLIST_ID"
Write-Host "$APP remove --age --days 180 --playlist $PLAYLIST_ID"
Write-Host "$APP remove --artist --interactive --playlist $PLAYLIST_ID"
Write-Host "$APP remove --artist --name 'Taylor Swift' --playlist $PLAYLIST_ID"
Write-Host ""

# Create operations
Write-Host "➕ Create Playlist Operations:" -ForegroundColor Green
Write-Host ""

Write-Host "🆕 Fresh Playlists (Recent Tracks):" -ForegroundColor Magenta
Write-Host "$APP create --type fresh --days 30 --name 'Last 30 Days' --playlist $PLAYLIST_ID"
Write-Host "$APP create --type fresh --days 90 --name 'Recent Hits' --playlist $PLAYLIST_ID"
Write-Host "$APP create --type fresh --days 180 --name 'Half Year Collection' --playlist $PLAYLIST_ID"
Write-Host ""

Write-Host "📦 Chunk Playlists (Split Large Playlists):" -ForegroundColor Blue
Write-Host "$APP create --type chunk --name 'MyBigPlaylist' --size 250 --playlist $PLAYLIST_ID"
Write-Host "$APP create --type chunk --name 'SmallChunks' --size 100 --playlist $PLAYLIST_ID"
Write-Host "$APP create --type chunk --name 'MegaList' --size 500 --overwrite --playlist $PLAYLIST_ID"
Write-Host ""

Write-Host "🎵 Genre Playlists:" -ForegroundColor DarkYellow
Write-Host "$APP create --type genre --interactive --playlist $PLAYLIST_ID"
Write-Host "$APP create --type genre --genre 'rock' --name 'Rock Collection' --playlist $PLAYLIST_ID"
Write-Host "$APP create --type genre --genre 'pop' --name 'Pop Hits' --overwrite --playlist $PLAYLIST_ID"
Write-Host "$APP create --type genre --genre 'electronic' --name 'Electronic Vibes' --playlist $PLAYLIST_ID"
Write-Host ""

# Configuration examples
Write-Host "⚙️ Configuration:" -ForegroundColor DarkCyan
Write-Host ""
Write-Host "Using environment variables:"
Write-Host "`$env:SPOTIFY_CLIENT_ID='your_client_id'"
Write-Host "`$env:SPOTIFY_CLIENT_SECRET='your_client_secret'"
Write-Host "`$env:SPOTIFY_REDIRECT_URI='http://127.0.0.1:8080/callback'"
Write-Host ""

Write-Host "Using config file (~/.spotify-shuffle.yaml):"
Write-Host "spotify:"
Write-Host "  client_id: 'your_client_id'"
Write-Host "  client_secret: 'your_client_secret'"
Write-Host "  redirect_uri: 'http://127.0.0.1:8080/callback'"
Write-Host ""

# Help and troubleshooting
Write-Host "❓ Help & Info:" -ForegroundColor Gray
Write-Host "$APP --help"
Write-Host "$APP create --help"
Write-Host "$APP remove --help"
Write-Host ""

Write-Host "🎯 PowerShell Pro Tips:" -ForegroundColor Yellow
Write-Host "• Use single quotes for strings with spaces: 'My Playlist Name'"
Write-Host "• PowerShell has better Unicode support than Command Prompt"
Write-Host "• Use Tab completion for commands and file names"
Write-Host "• Use --interactive flag for guided prompts"
Write-Host "• Use --overwrite to replace existing playlists"
Write-Host "• Set execution policy if needed: Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser"
Write-Host ""

Write-Host "📁 Example PowerShell Workflow:" -ForegroundColor Green
Write-Host "1. $APP"
Write-Host "2. Follow interactive prompts to authenticate"
Write-Host "3. Browse and select your playlists"
Write-Host "4. Choose operations from guided menus"
Write-Host "5. Confirm changes with safety prompts"
Write-Host ""

Write-Host "Ready to try it? Press any key to launch interactive mode..." -ForegroundColor White
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
& $APP