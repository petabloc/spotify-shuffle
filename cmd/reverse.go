package cmd

import (
	"context"
	"fmt"

	"github.com/petabloc/spotify-shuffle/internal/playlist"
	"github.com/spf13/cobra"
	"github.com/zmb3/spotify/v2"
)

// reverseCmd represents the reverse command
var reverseCmd = &cobra.Command{
	Use:   "reverse",
	Short: "Reverse the order of tracks in a playlist",
	Long:  `Reverses the current order of all tracks in the specified playlist.`,
	RunE:  runReverse,
}

func runReverse(cmd *cobra.Command, args []string) error {
	return runPlaylistCommand(func(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
		fmt.Println("🔄 Reversing playlist order...")

		if err := manager.ReversePlaylist(ctx, playlistID); err != nil {
			return fmt.Errorf("failed to reverse playlist: %w", err)
		}

		fmt.Println("✅ Playlist reversed successfully!")
		return nil
	})
}

func init() {
	rootCmd.AddCommand(reverseCmd)
}
