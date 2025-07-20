package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/petabloc/spotify-shuffle/internal/auth"
	"github.com/petabloc/spotify-shuffle/internal/config"
	"github.com/petabloc/spotify-shuffle/internal/playlist"
	"github.com/zmb3/spotify/v2"
)

// PlaylistCommandFunc represents a function that operates on a playlist
type PlaylistCommandFunc func(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error

// runPlaylistCommand is a helper that sets up auth and runs a playlist command
func runPlaylistCommand(fn PlaylistCommandFunc) error {
	// Extract playlist ID from URL if needed
	pid := extractPlaylistID(playlistID)
	if pid == "" {
		return fmt.Errorf("playlist ID or URL is required. Use --playlist flag or run in interactive mode with 'spotify-shuffle interactive'")
	}
	
	// Get authenticated client
	client, err := getAuthenticatedClient()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	
	ctx := context.Background()
	
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

// getAuthenticatedClient creates and returns an authenticated Spotify client
func getAuthenticatedClient() (*spotify.Client, error) {
	// Get Spotify configuration
	spotifyConfig := config.GetSpotify()
	if spotifyConfig.ClientID == "" || spotifyConfig.ClientSecret == "" {
		return nil, fmt.Errorf("Spotify credentials not configured. Please run 'spotify-shuffle interactive' to set up your credentials")
	}
	
	// Check for placeholder values
	if spotifyConfig.ClientID == "your_spotify_client_id" || spotifyConfig.ClientSecret == "your_spotify_client_secret" {
		return nil, fmt.Errorf("please update your Spotify credentials in the config file or run 'spotify-shuffle interactive' for guided setup")
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
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}
	
	return client, nil
}