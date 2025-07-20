package cmd

import (
	"fmt"
	"os"

	"github.com/petabloc/spotify-shuffle/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile         string
	playlistID      string
	interactiveMode bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "spotify-shuffle",
	Short: "A powerful Spotify playlist manager",
	Long: `Spotify Shuffle is a CLI tool for managing your Spotify playlists.
Features include shuffling, sorting, reversing, removing tracks, and creating new playlists.

Examples:
  spotify-shuffle interactive                                    # Interactive mode
  spotify-shuffle --playlist 37i9dQZF1DXcBWIGoYBM5M
  spotify-shuffle shuffle --playlist 37i9dQZF1DXcBWIGoYBM5M
  spotify-shuffle sort --by title --playlist 37i9dQZF1DXcBWIGoYBM5M`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If interactive flag is set or no arguments provided, run interactive mode
		if interactiveMode || len(args) == 0 {
			return runInteractive(cmd, args)
		}

		// Otherwise show help
		return cmd.Help()
	},
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
	rootCmd.PersistentFlags().StringVarP(&playlistID, "playlist", "p", "", "Spotify playlist ID or URL (required for non-interactive commands)")
	rootCmd.PersistentFlags().BoolVarP(&interactiveMode, "interactive", "i", false, "Run in interactive mode")
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
