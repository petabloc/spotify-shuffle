package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/petabloc/spotify-shuffle/internal/playlist"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

var (
	removeByAge    bool
	removeByArtist bool
	removeDays     int
	artistName     string
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove tracks from playlist",
	Long:  `Remove tracks from playlist by age or artist name.`,
	RunE:  runRemove,
}

func runRemove(cmd *cobra.Command, args []string) error {
	return runPlaylistCommand(func(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
		if removeByAge && removeByArtist {
			return fmt.Errorf("cannot use both --age and --artist flags")
		}

		if !removeByAge && !removeByArtist {
			return fmt.Errorf("must specify either --age or --artist")
		}

		if removeByAge {
			return removeByTrackAge(ctx, manager, playlistID)
		}

		if removeByArtist {
			return removeByTrackArtist(ctx, manager, playlistID)
		}

		return nil
	})
}

func removeByTrackAge(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
	if removeDays <= 0 {
		return fmt.Errorf("days must be greater than 0")
	}

	fmt.Printf("🔍 Removing tracks older than %d days...\n", removeDays)

	// Ask for confirmation
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("⚠️  This will permanently remove tracks from your playlist. Continue? (y/N): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		fmt.Println("❌ Operation cancelled")
		return nil
	}

	removedCount, err := manager.RemoveOldTracks(ctx, playlistID, removeDays)
	if err != nil {
		return fmt.Errorf("failed to remove old tracks: %w", err)
	}

	if removedCount > 0 {
		fmt.Printf("✅ Removed %d old tracks!\n", removedCount)
	} else {
		fmt.Println("ℹ️  No tracks found older than specified time period")
	}

	return nil
}

func removeByTrackArtist(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
	var artist string

	if artistName != "" {
		artist = artistName
	} else {
		// Show available artists
		fmt.Println("👨‍🎤 Getting artists from playlist...")
		artists, err := manager.GetUniqueArtists(ctx, playlistID)
		if err != nil {
			return fmt.Errorf("failed to get artists: %w", err)
		}

		if len(artists) == 0 {
			fmt.Println("ℹ️  No artists found in playlist")
			return nil
		}

		fmt.Printf("\n👨‍🎤 Found %d unique artists:\n", len(artists))
		maxShow := 20
		for i, a := range artists {
			if i >= maxShow {
				fmt.Printf("... and %d more\n", len(artists)-maxShow)
				break
			}
			fmt.Printf("%2d. %s\n", i+1, a)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nEnter artist number or name: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("❌ Operation cancelled")
			return nil
		}

		// Check if input is a number
		if num, err := strconv.Atoi(input); err == nil {
			if num >= 1 && num <= len(artists) && num <= maxShow {
				artist = artists[num-1]
			} else {
				return fmt.Errorf("invalid artist number: %d", num)
			}
		} else {
			artist = input
		}
	}

	fmt.Printf("🔍 Removing all tracks by '%s'...\n", artist)

	// Ask for confirmation
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("⚠️  This will permanently remove all tracks by this artist. Continue? (y/N): ")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))

	if response != "y" && response != "yes" {
		fmt.Println("❌ Operation cancelled")
		return nil
	}

	removedCount, err := manager.RemoveTracksByArtist(ctx, playlistID, artist)
	if err != nil {
		return fmt.Errorf("failed to remove tracks by artist: %w", err)
	}

	if removedCount > 0 {
		fmt.Printf("✅ Removed %d tracks by '%s'!\n", removedCount, artist)
	} else {
		fmt.Printf("ℹ️  No tracks found by '%s'\n", artist)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVar(&removeByAge, "age", false, "Remove tracks by age")
	removeCmd.Flags().BoolVar(&removeByArtist, "artist", false, "Remove tracks by artist")
	removeCmd.Flags().IntVar(&removeDays, "days", 0, "Number of days (use with --age)")
	removeCmd.Flags().StringVar(&artistName, "name", "", "Artist name (use with --artist)")
}
