package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Spotify SpotifyConfig `mapstructure:"spotify"`
}

type SpotifyConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURI  string `mapstructure:"redirect_uri"`
}

var cfg Config

// SetConfigFile sets the config file explicitly
func SetConfigFile(file string) {
	viper.SetConfigFile(file)
}

// SetConfigPaths sets the config search paths
func SetConfigPaths(home string) {
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.SetConfigName(".spotify-shuffle")
}

// ReadConfig reads the configuration
func ReadConfig() error {
	// Environment variables
	viper.SetEnvPrefix("SPOTIFY")
	viper.AutomaticEnv()

	// Default values
	viper.SetDefault("spotify.redirect_uri", "http://127.0.0.1:8080/callback")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		// Config file not found, create default
		if err := createDefaultConfig(); err != nil {
			return err
		}
	}

	return viper.Unmarshal(&cfg)
}

// Get returns the current configuration
func Get() Config {
	return cfg
}

// GetSpotify returns Spotify configuration
func GetSpotify() SpotifyConfig {
	return cfg.Spotify
}

// createDefaultConfig creates a default config file
func createDefaultConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".spotify-shuffle.yaml")
	
	defaultConfig := `# Spotify Shuffle Configuration
# Get your credentials from: https://developer.spotify.com/dashboard

spotify:
  client_id: "your_spotify_client_id"
  client_secret: "your_spotify_client_secret"
  redirect_uri: "http://127.0.0.1:8080/callback"

# You can also set these as environment variables:
# export SPOTIFY_CLIENT_ID="your_client_id"
# export SPOTIFY_CLIENT_SECRET="your_client_secret"
# export SPOTIFY_REDIRECT_URI="http://127.0.0.1:8080/callback"
`

	return os.WriteFile(configPath, []byte(defaultConfig), 0644)
}