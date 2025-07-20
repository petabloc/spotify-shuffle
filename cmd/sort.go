package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/spotify-shuffle/internal/playlist"
	"github.com/zmb3/spotify/v2"
)

var sortBy string

// sortCmd represents the sort command
var sortCmd = &cobra.Command{
	Use:   "sort",
	Short: "Sort playlist tracks",
	Long:  `Sort playlist tracks alphabetically by title or artist name.`,
	RunE:  runSort,
}

func runSort(cmd *cobra.Command, args []string) error {
	return runPlaylistCommand(func(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
		switch sortBy {
		case "title":
			fmt.Println("🔤 Sorting playlist by title...")
			if err := manager.SortPlaylistByTitle(ctx, playlistID); err != nil {
				return fmt.Errorf("failed to sort playlist by title: %w", err)
			}
			fmt.Println("✅ Playlist sorted by title successfully!")
			
		case "artist":
			fmt.Println("👨‍🎤 Sorting playlist by artist...")
			if err := manager.SortPlaylistByArtist(ctx, playlistID); err != nil {
				return fmt.Errorf("failed to sort playlist by artist: %w", err)
			}
			fmt.Println("✅ Playlist sorted by artist successfully!")
			
		default:
			return fmt.Errorf("invalid sort option: %s (use 'title' or 'artist')", sortBy)
		}
		
		return nil
	})
}

func init() {
	rootCmd.AddCommand(sortCmd)
	sortCmd.Flags().StringVar(&sortBy, "by", "title", "Sort by 'title' or 'artist'")
}