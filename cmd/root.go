package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/spotify-shuffle/internal/config"
)

var (
	cfgFile    string
	playlistID string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spotify-shuffle",
	Short: "A powerful Spotify playlist manager",
	Long: `Spotify Shuffle is a CLI tool for managing your Spotify playlists.
Features include shuffling, sorting, reversing, removing tracks, and creating new playlists.

Examples:
  spotify-shuffle --playlist 37i9dQZF1DXcBWIGoYBM5M
  spotify-shuffle shuffle --playlist 37i9dQZF1DXcBWIGoYBM5M
  spotify-shuffle sort --by title --playlist 37i9dQZF1DXcBWIGoYBM5M`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spotify-shuffle.yaml)")
	rootCmd.PersistentFlags().StringVarP(&playlistID, "playlist", "p", "", "Spotify playlist ID or URL (required)")
	rootCmd.MarkPersistentFlagRequired("playlist")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		config.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		config.SetConfigPaths(home)
	}

	if err := config.ReadConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading config:", err)
	}
}