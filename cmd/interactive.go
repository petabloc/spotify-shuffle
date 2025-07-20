package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/petabloc/spotify-shuffle/internal/config"
	"github.com/petabloc/spotify-shuffle/internal/playlist"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

// interactiveCmd represents the interactive mode command
var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Run in interactive mode",
	Long: `Interactive mode provides a guided interface for managing your Spotify playlists.
You can select playlists, choose operations, and get real-time feedback.`,
	RunE: runInteractive,
}

func runInteractive(cmd *cobra.Command, args []string) error {
	fmt.Println("🎵 Welcome to Spotify Shuffle Interactive Mode!")
	fmt.Println("===============================================")

	// Check if configuration exists and is valid
	if !config.IsConfigured() {
		if err := interactiveSetup(); err != nil {
			return fmt.Errorf("setup failed: %w", err)
		}
	}

	// Get authenticated client
	client, err := getAuthenticatedClient()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	manager := playlist.NewManager(client)
	reader := bufio.NewReader(os.Stdin)

	for {
		// Select playlist
		selectedPlaylist, err := selectPlaylist(cmd.Context(), client, reader)
		if err != nil {
			return err
		}
		if selectedPlaylist == nil {
			fmt.Println("👋 Goodbye!")
			return nil
		}

		// Show playlist operations menu
		if err := showPlaylistMenu(cmd.Context(), manager, *selectedPlaylist, reader); err != nil {
			return err
		}

		// Ask if user wants to continue with another playlist
		fmt.Print("\n🔄 Would you like to work with another playlist? (y/N): ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println("👋 Goodbye!")
			return nil
		}
		fmt.Println()
	}
}

func selectPlaylist(ctx context.Context, client *spotify.Client, reader *bufio.Reader) (*spotify.SimplePlaylist, error) {
	fmt.Println("\n📋 Select a playlist:")
	fmt.Println("1. Enter playlist ID/URL manually")
	fmt.Println("2. Choose from your playlists")
	fmt.Println("3. Exit")

	fmt.Print("\nChoose option (1-3): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return selectPlaylistManually(reader)
	case "2":
		return selectFromUserPlaylists(ctx, client, reader)
	case "3":
		return nil, nil
	default:
		fmt.Println("❌ Invalid choice. Please enter 1, 2, or 3.")
		return selectPlaylist(ctx, client, reader)
	}
}

func selectPlaylistManually(reader *bufio.Reader) (*spotify.SimplePlaylist, error) {
	fmt.Print("🔗 Enter playlist ID or URL: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return nil, fmt.Errorf("playlist ID/URL cannot be empty")
	}

	// Extract ID from URL if needed
	playlistID := extractPlaylistID(input)

	// Create a simple playlist object (we'll get full details later)
	return &spotify.SimplePlaylist{
		ID:   spotify.ID(playlistID),
		Name: "Selected Playlist",
	}, nil
}

func selectFromUserPlaylists(ctx context.Context, client *spotify.Client, reader *bufio.Reader) (*spotify.SimplePlaylist, error) {
	return selectFromUserPlaylistsWithOffset(ctx, client, reader, 0)
}

func selectFromUserPlaylistsWithOffset(ctx context.Context, client *spotify.Client, reader *bufio.Reader, offset int) (*spotify.SimplePlaylist, error) {
	fmt.Println("🔍 Loading your playlists...")

	// Get current user
	user, err := client.CurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	// Get user's playlists with a higher limit to support pagination
	playlists, err := client.GetPlaylistsForUser(ctx, user.ID, spotify.Limit(50), spotify.Offset(offset))
	if err != nil {
		return nil, fmt.Errorf("failed to get playlists: %w", err)
	}

	if len(playlists.Playlists) == 0 {
		if offset == 0 {
			fmt.Println("ℹ️  No playlists found")
			return selectPlaylistManually(reader)
		} else {
			fmt.Println("ℹ️  No more playlists found")
			return selectFromUserPlaylistsWithOffset(ctx, client, reader, 0) // Go back to first page
		}
	}

	totalPlaylists := playlists.Total
	fmt.Printf("\n👤 Found %d playlists for %s:\n", totalPlaylists, user.DisplayName)

	// Show playlists
	pageSize := 20
	endIndex := min(len(playlists.Playlists), pageSize)
	
	for i := 0; i < endIndex; i++ {
		playlist := playlists.Playlists[i]
		fmt.Printf("%2d. %s (%d tracks)\n", offset+i+1, playlist.Name, playlist.Tracks.Total)
	}

	// Build options
	options := []string{}
	if endIndex > 0 {
		options = append(options, fmt.Sprintf("1-%d", offset+endIndex))
	}
	
	hasMore := offset+len(playlists.Playlists) < totalPlaylists
	if hasMore {
		options = append(options, "n (next)")
	}
	
	if offset > 0 {
		options = append(options, "p (previous)")
	}
	
	options = append(options, "Enter (manual entry)")

	fmt.Printf("\n📋 Showing %d-%d of %d playlists\n", offset+1, offset+endIndex, totalPlaylists)
	fmt.Printf("Choose: %s: ", strings.Join(options, ", "))
	
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "" {
		return selectPlaylistManually(reader)
	}

	// Handle navigation
	if input == "n" && hasMore {
		return selectFromUserPlaylistsWithOffset(ctx, client, reader, offset+pageSize)
	}
	
	if input == "p" && offset > 0 {
		newOffset := offset - pageSize
		if newOffset < 0 {
			newOffset = 0
		}
		return selectFromUserPlaylistsWithOffset(ctx, client, reader, newOffset)
	}

	// Handle playlist selection
	num, err := strconv.Atoi(input)
	if err != nil || num < offset+1 || num > offset+endIndex {
		fmt.Printf("❌ Invalid choice. Please enter a number between %d and %d, or use navigation options.\n", offset+1, offset+endIndex)
		return selectFromUserPlaylistsWithOffset(ctx, client, reader, offset)
	}

	// Convert to 0-based index within current page
	localIndex := num - offset - 1
	selected := playlists.Playlists[localIndex]
	return &selected, nil
}

func showPlaylistMenu(ctx context.Context, manager *playlist.Manager, playlist spotify.SimplePlaylist, reader *bufio.Reader) error {
	fmt.Printf("\n🎵 Working with playlist: %s\n", playlist.Name)
	fmt.Println("==========================================")

	for {
		fmt.Println("\n📋 Choose an operation:")
		fmt.Println("1. 🔀 Shuffle tracks")
		fmt.Println("2. 🔤 Sort tracks")
		fmt.Println("3. 🔄 Reverse tracks")
		fmt.Println("4. 🗑️  Remove tracks")
		fmt.Println("5. ➕ Create new playlist")
		fmt.Println("6. ℹ️  Show playlist info")
		fmt.Println("7. 🔙 Go back to playlist selection")

		fmt.Print("\nChoose operation (1-7): ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if err := interactiveShuffle(ctx, manager, playlist.ID); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
		case "2":
			if err := interactiveSort(ctx, manager, playlist.ID, reader); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
		case "3":
			if err := interactiveReverse(ctx, manager, playlist.ID); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
		case "4":
			if err := interactiveRemove(ctx, manager, playlist.ID, reader); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
		case "5":
			if err := interactiveCreate(ctx, manager, playlist.ID, reader); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
		case "6":
			if err := showPlaylistInfo(ctx, manager, playlist.ID); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
		case "7":
			return nil
		default:
			fmt.Println("❌ Invalid choice. Please enter a number between 1 and 7.")
		}
	}
}

func interactiveShuffle(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
	fmt.Println("🔀 Shuffling playlist...")

	if !confirmAction("shuffle the playlist") {
		fmt.Println("❌ Operation cancelled")
		return nil
	}

	if err := manager.ShufflePlaylist(ctx, playlistID); err != nil {
		return err
	}

	fmt.Println("✅ Playlist shuffled successfully!")
	return nil
}

func interactiveSort(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID, reader *bufio.Reader) error {
	fmt.Println("\n🔤 Sort options:")
	fmt.Println("1. Sort by title (A-Z)")
	fmt.Println("2. Sort by artist (A-Z)")
	fmt.Println("3. Cancel")

	fmt.Print("Choose sort method (1-3): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var sortBy string
	switch choice {
	case "1":
		sortBy = "title"
	case "2":
		sortBy = "artist"
	case "3":
		fmt.Println("❌ Operation cancelled")
		return nil
	default:
		fmt.Println("❌ Invalid choice")
		return nil
	}

	if !confirmAction(fmt.Sprintf("sort the playlist by %s", sortBy)) {
		fmt.Println("❌ Operation cancelled")
		return nil
	}

	fmt.Printf("🔤 Sorting playlist by %s...\n", sortBy)
	if err := manager.SortPlaylist(ctx, playlistID, sortBy); err != nil {
		return err
	}

	fmt.Printf("✅ Playlist sorted by %s successfully!\n", sortBy)
	return nil
}

func interactiveReverse(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
	fmt.Println("🔄 Reversing playlist...")

	if !confirmAction("reverse the playlist order") {
		fmt.Println("❌ Operation cancelled")
		return nil
	}

	if err := manager.ReversePlaylist(ctx, playlistID); err != nil {
		return err
	}

	fmt.Println("✅ Playlist reversed successfully!")
	return nil
}

func interactiveRemove(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID, reader *bufio.Reader) error {
	fmt.Println("\n🗑️  Remove options:")
	fmt.Println("1. Remove tracks by age")
	fmt.Println("2. Remove tracks by artist")
	fmt.Println("3. Cancel")

	fmt.Print("Choose remove method (1-3): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return interactiveRemoveByAge(ctx, manager, playlistID, reader)
	case "2":
		return interactiveRemoveByArtist(ctx, manager, playlistID, reader)
	case "3":
		fmt.Println("❌ Operation cancelled")
		return nil
	default:
		fmt.Println("❌ Invalid choice")
		return nil
	}
}

func interactiveRemoveByAge(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID, reader *bufio.Reader) error {
	fmt.Println("\n📅 Age options:")
	fmt.Println("1. 90 days")
	fmt.Println("2. 180 days")
	fmt.Println("3. 1 year")
	fmt.Println("4. 2 years")
	fmt.Println("5. 3 years")
	fmt.Println("6. Custom")
	fmt.Println("7. Cancel")

	fmt.Print("Choose age threshold (1-7): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var days int
	switch choice {
	case "1":
		days = 90
	case "2":
		days = 180
	case "3":
		days = 365
	case "4":
		days = 730
	case "5":
		days = 1095
	case "6":
		fmt.Print("Enter number of days: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		var err error
		days, err = strconv.Atoi(input)
		if err != nil || days <= 0 {
			fmt.Println("❌ Invalid number of days")
			return nil
		}
	case "7":
		fmt.Println("❌ Operation cancelled")
		return nil
	default:
		fmt.Println("❌ Invalid choice")
		return nil
	}

	if !confirmAction(fmt.Sprintf("remove tracks older than %d days", days)) {
		fmt.Println("❌ Operation cancelled")
		return nil
	}

	fmt.Printf("🗑️  Removing tracks older than %d days...\n", days)
	count, err := manager.RemoveOldTracks(ctx, playlistID, days)
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Printf("✅ Removed %d old tracks!\n", count)
	} else {
		fmt.Println("ℹ️  No tracks found older than specified time period")
	}
	return nil
}

func interactiveRemoveByArtist(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID, reader *bufio.Reader) error {
	return interactiveRemoveByArtistWithOffset(ctx, manager, playlistID, reader, 0)
}

func interactiveRemoveByArtistWithOffset(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID, reader *bufio.Reader, offset int) error {
	if offset == 0 {
		fmt.Println("👨‍🎤 Getting artists from playlist...")
	}
	
	artists, err := manager.GetUniqueArtists(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to get artists: %w", err)
	}

	if len(artists) == 0 {
		fmt.Println("ℹ️  No artists found in playlist")
		return nil
	}

	pageSize := 20
	startIndex := offset
	endIndex := min(offset+pageSize, len(artists))

	if startIndex >= len(artists) {
		fmt.Println("ℹ️  No more artists found")
		return interactiveRemoveByArtistWithOffset(ctx, manager, playlistID, reader, 0) // Go back to first page
	}

	fmt.Printf("\n👨‍🎤 Found %d unique artists:\n", len(artists))
	
	for i := startIndex; i < endIndex; i++ {
		fmt.Printf("%2d. %s\n", i+1, artists[i])
	}

	// Build options
	options := []string{}
	if endIndex > startIndex {
		options = append(options, fmt.Sprintf("1-%d", endIndex))
	}
	
	hasMore := endIndex < len(artists)
	if hasMore {
		options = append(options, "n (next)")
	}
	
	if offset > 0 {
		options = append(options, "p (previous)")
	}
	
	options = append(options, "name (enter artist name)", "Enter (cancel)")

	fmt.Printf("\n📋 Showing %d-%d of %d artists\n", startIndex+1, endIndex, len(artists))
	fmt.Printf("Choose: %s: ", strings.Join(options, ", "))
	
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		fmt.Println("❌ Operation cancelled")
		return nil
	}

	// Handle navigation
	inputLower := strings.ToLower(input)
	if inputLower == "n" && hasMore {
		return interactiveRemoveByArtistWithOffset(ctx, manager, playlistID, reader, offset+pageSize)
	}
	
	if inputLower == "p" && offset > 0 {
		newOffset := offset - pageSize
		if newOffset < 0 {
			newOffset = 0
		}
		return interactiveRemoveByArtistWithOffset(ctx, manager, playlistID, reader, newOffset)
	}

	var artistName string
	if num, err := strconv.Atoi(input); err == nil {
		if num >= 1 && num <= len(artists) {
			artistName = artists[num-1]
		} else {
			fmt.Printf("❌ Invalid artist number: %d\n", num)
			return interactiveRemoveByArtistWithOffset(ctx, manager, playlistID, reader, offset)
		}
	} else {
		artistName = input
	}

	if !confirmAction(fmt.Sprintf("remove all tracks by '%s'", artistName)) {
		fmt.Println("❌ Operation cancelled")
		return nil
	}

	fmt.Printf("🗑️  Removing tracks by '%s'...\n", artistName)
	count, err := manager.RemoveTracksByArtist(ctx, playlistID, artistName)
	if err != nil {
		return err
	}

	if count > 0 {
		fmt.Printf("✅ Removed %d tracks by '%s'!\n", count, artistName)
	} else {
		fmt.Printf("ℹ️  No tracks found by '%s'\n", artistName)
	}
	return nil
}

func interactiveCreate(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID, reader *bufio.Reader) error {
	fmt.Println("\n➕ Create playlist options:")
	fmt.Println("1. Fresh playlist (recent tracks)")
	fmt.Println("2. Chunk playlists (split large playlist)")
	fmt.Println("3. Genre playlist")
	fmt.Println("4. Cancel")

	fmt.Print("Choose creation method (1-4): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return interactiveCreateFresh(ctx, manager, playlistID, reader)
	case "2":
		return interactiveCreateChunk(ctx, manager, playlistID, reader)
	case "3":
		return interactiveCreateGenre(ctx, manager, playlistID, reader)
	case "4":
		fmt.Println("❌ Operation cancelled")
		return nil
	default:
		fmt.Println("❌ Invalid choice")
		return nil
	}
}

func interactiveCreateFresh(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID, reader *bufio.Reader) error {
	fmt.Println("\n📅 Fresh playlist time range:")
	fmt.Println("1. Last 30 days")
	fmt.Println("2. Last 90 days")
	fmt.Println("3. Last 180 days")
	fmt.Println("4. Custom")
	fmt.Println("5. Cancel")

	fmt.Print("Choose time range (1-5): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	var days int
	switch choice {
	case "1":
		days = 30
	case "2":
		days = 90
	case "3":
		days = 180
	case "4":
		fmt.Print("Enter number of days: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		var err error
		days, err = strconv.Atoi(input)
		if err != nil || days <= 0 {
			fmt.Println("❌ Invalid number of days")
			return nil
		}
	case "5":
		fmt.Println("❌ Operation cancelled")
		return nil
	default:
		fmt.Println("❌ Invalid choice")
		return nil
	}

	fmt.Print("Enter name for the fresh playlist: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = fmt.Sprintf("Fresh - Last %d days", days)
	}

	fmt.Printf("➕ Creating fresh playlist '%s' with tracks from last %d days...\n", name, days)
	count, err := manager.CreateFreshPlaylist(ctx, playlistID, name, days, false)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Created fresh playlist '%s' with %d tracks!\n", name, count)
	return nil
}

func interactiveCreateChunk(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID, reader *bufio.Reader) error {
	fmt.Print("Enter chunk size (default 250): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	chunkSize := 250
	if input != "" {
		var err error
		chunkSize, err = strconv.Atoi(input)
		if err != nil || chunkSize <= 0 {
			fmt.Println("❌ Invalid chunk size")
			return nil
		}
	}

	fmt.Print("Enter base name for chunk playlists: ")
	baseName, _ := reader.ReadString('\n')
	baseName = strings.TrimSpace(baseName)
	if baseName == "" {
		baseName = "Chunk"
	}

	fmt.Printf("➕ Creating chunk playlists with %d tracks each...\n", chunkSize)
	count, err := manager.CreateChunkPlaylists(ctx, playlistID, baseName, chunkSize, false)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Created %d chunk playlists!\n", count)
	return nil
}

func interactiveCreateGenre(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID, reader *bufio.Reader) error {
	fmt.Print("Enter genre (e.g., rock, pop, jazz): ")
	genre, _ := reader.ReadString('\n')
	genre = strings.TrimSpace(genre)
	if genre == "" {
		fmt.Println("❌ Genre cannot be empty")
		return nil
	}

	fmt.Print("Enter name for the genre playlist: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name == "" {
		name = fmt.Sprintf("%s Collection", strings.Title(genre))
	}

	fmt.Printf("➕ Creating genre playlist '%s' for %s...\n", name, genre)
	count, err := manager.CreateGenrePlaylist(ctx, playlistID, name, genre, false)
	if err != nil {
		return err
	}

	fmt.Printf("✅ Created genre playlist '%s' with %d tracks!\n", name, count)
	return nil
}

func showPlaylistInfo(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
	fmt.Println("ℹ️  Getting playlist information...")

	tracks, err := manager.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return fmt.Errorf("failed to get playlist tracks: %w", err)
	}

	fmt.Printf("\n📊 Playlist Statistics:\n")
	fmt.Printf("📍 Total tracks: %d\n", len(tracks))

	if len(tracks) > 0 {
		// Get unique artists
		artistMap := make(map[string]int)
		for _, track := range tracks {
			for _, artist := range track.Artists {
				artistMap[artist]++
			}
		}

		fmt.Printf("👨‍🎤 Unique artists: %d\n", len(artistMap))

		// Show top artists
		if len(artistMap) > 0 {
			fmt.Println("\n🔥 Top artists:")
			type artistCount struct {
				name  string
				count int
			}

			var topArtists []artistCount
			for name, count := range artistMap {
				topArtists = append(topArtists, artistCount{name, count})
			}

			// Simple sort by count (descending)
			for i := 0; i < len(topArtists); i++ {
				for j := i + 1; j < len(topArtists); j++ {
					if topArtists[j].count > topArtists[i].count {
						topArtists[i], topArtists[j] = topArtists[j], topArtists[i]
					}
				}
			}

			maxShow := min(5, len(topArtists))
			for i := 0; i < maxShow; i++ {
				fmt.Printf("   %d. %s (%d tracks)\n", i+1, topArtists[i].name, topArtists[i].count)
			}
		}
	}

	return nil
}

func confirmAction(action string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("⚠️  This will %s. Continue? (y/N): ", action)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// Note: extractPlaylistID is defined in utils.go

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func interactiveSetup() error {
	fmt.Println("\n🔧 First-Time Setup")
	fmt.Println("===================")
	fmt.Println()
	fmt.Println("To use Spotify Shuffle, you need to create a Spotify app and get your credentials.")
	fmt.Println()
	fmt.Println("📋 Setup Steps:")
	fmt.Println("1. Go to: https://developer.spotify.com/dashboard")
	fmt.Println("2. Click 'Create an app'")
	fmt.Println("3. Fill in app name and description (anything you want)")
	fmt.Println("4. In 'Redirect URIs', add: http://127.0.0.1:8080/callback")
	fmt.Println("5. Copy your Client ID and Client Secret")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Check if user wants to proceed or has already done this
	fmt.Print("Have you already created a Spotify app? (y/N): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		fmt.Println()
		fmt.Println("🌐 Opening Spotify Developer Dashboard...")
		fmt.Println("Please complete the setup steps above, then come back here.")
		fmt.Println()
		fmt.Print("Press Enter when you've created your Spotify app and have your credentials ready...")
		reader.ReadString('\n')
		fmt.Println()
	}

	// Get Client ID
	fmt.Print("🔑 Enter your Spotify Client ID: ")
	clientID, _ := reader.ReadString('\n')
	clientID = strings.TrimSpace(clientID)

	if clientID == "" {
		return fmt.Errorf("Client ID cannot be empty")
	}

	// Get Client Secret
	fmt.Print("🔐 Enter your Spotify Client Secret: ")
	clientSecret, _ := reader.ReadString('\n')
	clientSecret = strings.TrimSpace(clientSecret)

	if clientSecret == "" {
		return fmt.Errorf("Client Secret cannot be empty")
	}

	// Set redirect URI (default)
	redirectURI := "http://127.0.0.1:8080/callback"
	fmt.Printf("🔗 Redirect URI (default: %s): ", redirectURI)
	customRedirectURI, _ := reader.ReadString('\n')
	customRedirectURI = strings.TrimSpace(customRedirectURI)

	if customRedirectURI != "" {
		redirectURI = customRedirectURI
	}

	// Update configuration
	config.SetSpotifyConfig(clientID, clientSecret, redirectURI)

	// Save configuration
	if err := config.SaveConfig(); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println()
	fmt.Println("✅ Configuration saved successfully!")
	fmt.Printf("📁 Config file: %s\n", getConfigPath())
	fmt.Println()
	fmt.Println("🔐 Next: You'll be prompted to authenticate with Spotify...")
	fmt.Println()

	return nil
}

func getConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.spotify-shuffle.yaml"
	}
	return home + "/.spotify-shuffle.yaml"
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
