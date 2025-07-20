package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/petabloc/spotify-shuffle/internal/playlist"
	"github.com/zmb3/spotify/v2"
)

// shuffleCmd represents the shuffle command
var shuffleCmd = &cobra.Command{
	Use:   "shuffle",
	Short: "Shuffle the order of tracks in a playlist",
	Long:  `Randomly reorders all tracks in the specified Spotify playlist.`,
	RunE:  runShuffle,
}

func runShuffle(cmd *cobra.Command, args []string) error {
	return runPlaylistCommand(func(ctx context.Context, manager *playlist.Manager, playlistID spotify.ID) error {
		fmt.Println("🔀 Shuffling playlist...")
		
		if err := manager.ShufflePlaylist(ctx, playlistID); err != nil {
			return fmt.Errorf("failed to shuffle playlist: %w", err)
		}
		
		fmt.Println("✅ Playlist shuffled successfully!")
		return nil
	})
}

func init() {
	rootCmd.AddCommand(shuffleCmd)
}