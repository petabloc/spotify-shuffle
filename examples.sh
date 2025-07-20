#!/bin/bash

# Spotify Shuffle - Go Edition
# Example Usage Script

PLAYLIST_ID="37i9dQZF1DXcBWIGoYBM5M"  # Replace with your playlist ID
APP="./spotify-shuffle"

echo "🎵 Spotify Shuffle - Go Edition Examples"
echo "========================================"
echo ""

# Basic operations
echo "📝 Basic Operations:"
echo "$APP shuffle --playlist $PLAYLIST_ID"
echo "$APP sort --by title --playlist $PLAYLIST_ID"
echo "$APP sort --by artist --playlist $PLAYLIST_ID"
echo "$APP reverse --playlist $PLAYLIST_ID"
echo ""

# Remove operations
echo "🗑️ Remove Operations:"
echo "$APP remove --age --days 90 --playlist $PLAYLIST_ID"
echo "$APP remove --age --days 180 --playlist $PLAYLIST_ID"
echo "$APP remove --artist --interactive --playlist $PLAYLIST_ID"
echo "$APP remove --artist --name 'Taylor Swift' --playlist $PLAYLIST_ID"
echo ""

# Create operations
echo "➕ Create Playlist Operations:"
echo ""

echo "🆕 Fresh Playlists (Recent Tracks):"
echo "$APP create --type fresh --days 30 --name 'Last 30 Days' --playlist $PLAYLIST_ID"
echo "$APP create --type fresh --days 90 --name 'Recent Hits' --playlist $PLAYLIST_ID"
echo "$APP create --type fresh --days 180 --name 'Half Year Collection' --playlist $PLAYLIST_ID"
echo ""

echo "📦 Chunk Playlists (Split Large Playlists):"
echo "$APP create --type chunk --name 'MyBigPlaylist' --size 250 --playlist $PLAYLIST_ID"
echo "$APP create --type chunk --name 'SmallChunks' --size 100 --playlist $PLAYLIST_ID"
echo "$APP create --type chunk --name 'MegaList' --size 500 --overwrite --playlist $PLAYLIST_ID"
echo ""

echo "🎵 Genre Playlists:"
echo "$APP create --type genre --interactive --playlist $PLAYLIST_ID"
echo "$APP create --type genre --genre 'rock' --name 'Rock Collection' --playlist $PLAYLIST_ID"
echo "$APP create --type genre --genre 'pop' --name 'Pop Hits' --overwrite --playlist $PLAYLIST_ID"
echo "$APP create --type genre --genre 'electronic' --name 'Electronic Vibes' --playlist $PLAYLIST_ID"
echo ""

# Advanced examples
echo "🚀 Advanced Examples:"
echo ""

echo "Chain operations (using bash):"
echo "# First shuffle, then create fresh playlist"
echo "$APP shuffle --playlist $PLAYLIST_ID && \\"
echo "$APP create --type fresh --days 30 --name 'Shuffled Recent' --playlist $PLAYLIST_ID"
echo ""

echo "Batch processing multiple playlists:"
echo "for id in '37i9dQZF1DXcBWIGoYBM5M' '37i9dQZF1DX0XUsuxWHRQV'; do"
echo "  echo \"Processing playlist: \$id\""
echo "  $APP shuffle --playlist \$id"
echo "done"
echo ""

echo "Create genre playlists for all major genres:"
echo "for genre in 'rock' 'pop' 'hip hop' 'electronic' 'indie'; do"
echo "  $APP create --type genre --genre \"\$genre\" --name \"\${genre^} Collection\" --overwrite --playlist $PLAYLIST_ID"
echo "done"
echo ""

# Configuration examples
echo "⚙️ Configuration:"
echo ""
echo "Using environment variables:"
echo "export SPOTIFY_CLIENT_ID='your_client_id'"
echo "export SPOTIFY_CLIENT_SECRET='your_client_secret'"
echo "export SPOTIFY_REDIRECT_URI='http://127.0.0.1:8080/callback'"
echo ""

echo "Using config file (~/.spotify-shuffle.yaml):"
echo "spotify:"
echo "  client_id: 'your_client_id'"
echo "  client_secret: 'your_client_secret'"
echo "  redirect_uri: 'http://127.0.0.1:8080/callback'"
echo ""

# Help and troubleshooting
echo "❓ Help & Info:"
echo "$APP --help"
echo "$APP create --help"
echo "$APP remove --help"
echo ""

echo "🔍 Get playlist info without modifying:"
echo "# The app shows playlist info before any operation"
echo "$APP shuffle --playlist $PLAYLIST_ID --help  # Shows info then help"
echo ""

echo "🎯 Pro Tips:"
echo "• Use full Spotify URLs: https://open.spotify.com/playlist/37i9dQZF1DXcBWIGoYBM5M"
echo "• Use --interactive flag for guided prompts"
echo "• Use --overwrite to replace existing playlists"
echo "• Chunk size of 250 works well for most use cases"
echo "• Fresh playlists work best with frequently updated playlists"
echo "• Genre matching is case-insensitive and supports partial matches"
echo ""

echo "📁 Example Workflow:"
echo "1. $APP shuffle --playlist $PLAYLIST_ID"
echo "2. $APP create --type fresh --days 30 --name 'Recent Finds' --playlist $PLAYLIST_ID"
echo "3. $APP create --type genre --genre 'rock' --name 'Rock Only' --playlist $PLAYLIST_ID"
echo "4. $APP remove --age --days 365 --playlist $PLAYLIST_ID"