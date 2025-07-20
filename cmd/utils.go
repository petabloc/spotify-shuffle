package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/user/spotify-shuffle/internal/auth"
	"github.com/user/spotify-shuffle/internal/config"
	"github.com/user/spotify-shuffle/internal/playlist"
	"github.com/zmb3/spotify/v2"
)

// PlaylistCommandFunc represents a function that operates on a playlist
type PlaylistCommandFunc func(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error

// runPlaylistCommand is a helper that sets up auth and runs a playlist command
func runPlaylistCommand(fn PlaylistCommandFunc) error {
	// Extract playlist ID from URL if needed
	pid := extractPlaylistID(playlistID)
	if pid == "" {
		return fmt.Errorf("invalid playlist ID or URL")
	}
	
	// Get Spotify configuration
	spotifyConfig := config.GetSpotify()
	if spotifyConfig.ClientID == "" || spotifyConfig.ClientSecret == "" {
		return fmt.Errorf("Spotify credentials not configured. Please check your config file or environment variables")
	}
	
	// Create authenticator
	spotifyAuth := auth.NewSpotifyAuth(
		spotifyConfig.ClientID,
		spotifyConfig.ClientSecret,
		spotifyConfig.RedirectURI,
	)
	
	// Get authenticated client
	ctx := context.Background()
	client, err := spotifyAuth.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	
	// Get playlist info
	playlistInfo, err := client.GetPlaylist(ctx, spotify.ID(pid))
	if err != nil {
		return fmt.Errorf("failed to access playlist: %w", err)
	}
	
	fmt.Printf("\n📱 Playlist: %s\n", playlistInfo.Name)
	fmt.Printf("📊 Total tracks: %d\n", playlistInfo.Tracks.Total)
	
	// Create playlist manager and run command
	manager := playlist.NewManager(client)
	return fn(ctx, manager, spotify.ID(pid))
}

// extractPlaylistID extracts the playlist ID from a URL or returns the ID as-is
func extractPlaylistID(input string) string {
	input = strings.TrimSpace(input)
	
	// Handle Spotify URLs
	if strings.Contains(input, "playlist/") {
		parts := strings.Split(input, "playlist/")
		if len(parts) > 1 {
			id := parts[1]
			// Remove query parameters
			if idx := strings.Index(id, "?"); idx != -1 {
				id = id[:idx]
			}
			return id
		}
	}
	
	// Handle Spotify URIs
	if strings.HasPrefix(input, "spotify:playlist:") {
		return strings.TrimPrefix(input, "spotify:playlist:")
	}
	
	// Assume it's already a playlist ID
	return input
}