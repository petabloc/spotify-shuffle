package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/spotify-shuffle/internal/playlist"
	"github.com/zmb3/spotify/v2"
)

var (
	createType   string
	name         string
	days         int
	chunkSize    int
	genre        string
	overwrite    bool
	interactive  bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new playlists from existing playlist",
	Long:  `Create new playlists from an existing playlist using various methods: fresh (recent tracks), chunk (split into smaller playlists), or genre-based filtering.`,
	RunE:  runCreate,
}

func runCreate(cmd *cobra.Command, args []string) error {
	return runPlaylistCommand(func(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
		switch createType {
		case "fresh":
			return createFreshPlaylist(ctx, manager, playlistID)
		case "chunk":
			return createChunkPlaylists(ctx, manager, playlistID)
		case "genre":
			return createGenrePlaylist(ctx, manager, playlistID)
		default:
			return fmt.Errorf("invalid create type: %s (use 'fresh', 'chunk', or 'genre')", createType)
		}
	})
}

func createFreshPlaylist(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
	if days <= 0 {
		return fmt.Errorf("days must be greater than 0")
	}
	
	if name == "" {
		if interactive {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter name for the fresh playlist: ")
			input, _ := reader.ReadString('\n')
			name = strings.TrimSpace(input)
		}
		if name == "" {
			return fmt.Errorf("playlist name is required (use --name or --interactive)")
		}
	}
	
	// Check for existing playlist if not overwriting
	if !overwrite {
		if _, err := manager.FindPlaylistByName(ctx, name); err == nil {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("⚠️  Playlist '%s' already exists. Overwrite? (y/N): ", name)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			overwrite = (response == "y" || response == "yes")
			if !overwrite {
				fmt.Println("❌ Operation cancelled")
				return nil
			}
		}
	}
	
	fmt.Printf("🔍 Creating fresh playlist with tracks from last %d days...\n", days)
	
	trackCount, err := manager.CreateFreshPlaylist(ctx, playlistID, name, days, overwrite)
	if err != nil {
		return fmt.Errorf("failed to create fresh playlist: %w", err)
	}
	
	fmt.Printf("✅ Created fresh playlist '%s' with %d tracks!\n", name, trackCount)
	return nil
}

func createChunkPlaylists(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
	if chunkSize <= 0 {
		chunkSize = 250 // Default chunk size
	}
	
	if name == "" {
		if interactive {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter base name for chunk playlists: ")
			input, _ := reader.ReadString('\n')
			name = strings.TrimSpace(input)
		}
		if name == "" {
			return fmt.Errorf("base name is required (use --name or --interactive)")
		}
	}
	
	// Check for existing chunk playlists if not overwriting
	if !overwrite && interactive {
		// Check for existing chunks
		existingCount := 0
		for i := 0; i < 100; i++ { // Check first 100 possible chunks
			chunkName := fmt.Sprintf("%s-%02d", name, i)
			if _, err := manager.FindPlaylistByName(ctx, chunkName); err == nil {
				existingCount++
			}
		}
		
		if existingCount > 0 {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("⚠️  Found %d existing chunk playlists. Overwrite? (y/N): ", existingCount)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			overwrite = (response == "y" || response == "yes")
			if !overwrite {
				fmt.Println("❌ Operation cancelled")
				return nil
			}
		}
	}
	
	fmt.Printf("🔍 Creating chunk playlists with %d tracks per chunk...\n", chunkSize)
	
	createdCount, err := manager.CreateChunkPlaylists(ctx, playlistID, name, chunkSize, overwrite)
	if err != nil {
		return fmt.Errorf("failed to create chunk playlists: %w", err)
	}
	
	fmt.Printf("✅ Created %d chunk playlists!\n", createdCount)
	return nil
}

func createGenrePlaylist(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
	var selectedGenre string
	
	if genre != "" {
		selectedGenre = genre
	} else if interactive {
		// Show available genres
		fmt.Println("🎵 Getting genres from playlist...")
		genres, err := manager.GetPlaylistGenres(ctx, playlistID)
		if err != nil {
			return fmt.Errorf("failed to get genres: %w", err)
		}
		
		if len(genres) == 0 {
			fmt.Println("ℹ️  No genres found in playlist")
			return nil
		}
		
		// Sort genres by track count (descending)
		type genreCount struct {
			name  string
			count int
		}
		var sortedGenres []genreCount
		for g, c := range genres {
			sortedGenres = append(sortedGenres, genreCount{g, c})
		}
		sort.Slice(sortedGenres, func(i, j int) bool {
			return sortedGenres[i].count > sortedGenres[j].count
		})
		
		fmt.Printf("\n🎵 Found %d genres in playlist:\n", len(sortedGenres))
		maxShow := 20
		for i, gc := range sortedGenres {
			if i >= maxShow {
				fmt.Printf("... and %d more\n", len(sortedGenres)-maxShow)
				break
			}
			fmt.Printf("%2d. %s (%d tracks)\n", i+1, gc.name, gc.count)
		}
		
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter genre number or name: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		if input == "" {
			fmt.Println("❌ Operation cancelled")
			return nil
		}
		
		// Check if input is a number
		if num, err := strconv.Atoi(input); err == nil {
			if num >= 1 && num <= len(sortedGenres) && num <= maxShow {
				selectedGenre = sortedGenres[num-1].name
			} else {
				return fmt.Errorf("invalid genre number: %d", num)
			}
		} else {
			selectedGenre = input
		}
	} else {
		return fmt.Errorf("genre is required (use --genre or --interactive)")
	}
	
	if name == "" {
		if interactive {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Enter name for the '%s' playlist: ", selectedGenre)
			input, _ := reader.ReadString('\n')
			name = strings.TrimSpace(input)
		}
		if name == "" {
			return fmt.Errorf("playlist name is required (use --name or --interactive)")
		}
	}
	
	// Check for existing playlist if not overwriting
	if !overwrite {
		if _, err := manager.FindPlaylistByName(ctx, name); err == nil {
			if interactive {
				reader := bufio.NewReader(os.Stdin)
				fmt.Printf("⚠️  Playlist '%s' already exists. Overwrite? (y/N): ", name)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))
				overwrite = (response == "y" || response == "yes")
				if !overwrite {
					fmt.Println("❌ Operation cancelled")
					return nil
				}
			} else {
				return fmt.Errorf("playlist '%s' already exists (use --overwrite to replace)", name)
			}
		}
	}
	
	fmt.Printf("🔍 Creating genre playlist for '%s'...\n", selectedGenre)
	
	trackCount, err := manager.CreateGenrePlaylist(ctx, playlistID, name, selectedGenre, overwrite)
	if err != nil {
		return fmt.Errorf("failed to create genre playlist: %w", err)
	}
	
	fmt.Printf("✅ Created genre playlist '%s' with %d tracks!\n", name, trackCount)
	return nil
}

func init() {
	rootCmd.AddCommand(createCmd)
	
	createCmd.Flags().StringVar(&createType, "type", "", "Creation type: 'fresh', 'chunk', or 'genre' (required)")
	createCmd.Flags().StringVar(&name, "name", "", "Playlist name (or base name for chunks)")
	createCmd.Flags().IntVar(&days, "days", 30, "Number of days for fresh playlist (default: 30)")
	createCmd.Flags().IntVar(&chunkSize, "size", 250, "Tracks per chunk for chunk playlists (default: 250)")
	createCmd.Flags().StringVar(&genre, "genre", "", "Genre name for genre playlist")
	createCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing playlists")
	createCmd.Flags().BoolVar(&interactive, "interactive", false, "Use interactive mode for prompts")
	
	createCmd.MarkFlagRequired("type")
}