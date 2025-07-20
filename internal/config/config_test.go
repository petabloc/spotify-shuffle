package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestConfig(t *testing.T) {
	// Reset viper state between tests
	viper.Reset()

	tests := []struct {
		name     string
		envVars  map[string]string
		expected SpotifyConfig
	}{
		{
			name: "environment variables",
			envVars: map[string]string{
				"SPOTIFY_CLIENT_ID":     "test_client_id",
				"SPOTIFY_CLIENT_SECRET": "test_client_secret",
				"SPOTIFY_REDIRECT_URI":  "http://127.0.0.1:9999/callback",
			},
			expected: SpotifyConfig{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://127.0.0.1:9999/callback",
			},
		},
		{
			name: "default redirect URI",
			envVars: map[string]string{
				"SPOTIFY_CLIENT_ID":     "test_client_id",
				"SPOTIFY_CLIENT_SECRET": "test_client_secret",
			},
			expected: SpotifyConfig{
				ClientID:     "test_client_id",
				ClientSecret: "test_client_secret",
				RedirectURI:  "http://127.0.0.1:8080/callback",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper
			viper.Reset()

			// Clean up environment first
			os.Unsetenv("SPOTIFY_CLIENT_ID")
			os.Unsetenv("SPOTIFY_CLIENT_SECRET")
			os.Unsetenv("SPOTIFY_REDIRECT_URI")

			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Use a temporary directory for config
			tempDir := t.TempDir()
			SetConfigPaths(tempDir)

			err := ReadConfig()
			if err != nil {
				t.Fatalf("ReadConfig() error = %v", err)
			}

			spotify := GetSpotify()
			if spotify.ClientID != tt.expected.ClientID {
				t.Errorf("ClientID = %v, want %v", spotify.ClientID, tt.expected.ClientID)
			}
			if spotify.ClientSecret != tt.expected.ClientSecret {
				t.Errorf("ClientSecret = %v, want %v", spotify.ClientSecret, tt.expected.ClientSecret)
			}
			if spotify.RedirectURI != tt.expected.RedirectURI {
				t.Errorf("RedirectURI = %v, want %v", spotify.RedirectURI, tt.expected.RedirectURI)
			}
		})
	}
}

func TestSetConfigFile(t *testing.T) {
	// Reset viper
	viper.Reset()

	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test-config.yaml")

	// Create test config file
	configContent := `spotify:
  client_id: "file_client_id"
  client_secret: "file_client_secret"
  redirect_uri: "http://127.0.0.1:7777/callback"
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	SetConfigFile(configFile)
	err = ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() error = %v", err)
	}

	spotify := GetSpotify()
	expected := SpotifyConfig{
		ClientID:     "file_client_id",
		ClientSecret: "file_client_secret",
		RedirectURI:  "http://127.0.0.1:7777/callback",
	}

	if spotify != expected {
		t.Errorf("Config from file = %+v, want %+v", spotify, expected)
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	// Reset viper
	viper.Reset()

	tempDir := t.TempDir()

	// Mock home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	SetConfigPaths(tempDir)
	err := ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() error = %v", err)
	}

	// Check that default config file was created
	configPath := filepath.Join(tempDir, ".spotify-shuffle.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Default config file was not created at %s", configPath)
	}

	// Verify default values
	spotify := GetSpotify()
	if spotify.RedirectURI != "http://127.0.0.1:8080/callback" {
		t.Errorf("Default RedirectURI = %v, want %v", spotify.RedirectURI, "http://127.0.0.1:8080/callback")
	}
}

func TestGet(t *testing.T) {
	// Reset viper
	viper.Reset()

	tempDir := t.TempDir()
	SetConfigPaths(tempDir)

	err := ReadConfig()
	if err != nil {
		t.Fatalf("ReadConfig() error = %v", err)
	}

	config := Get()
	if config.Spotify.RedirectURI != "http://127.0.0.1:8080/callback" {
		t.Errorf("Get().Spotify.RedirectURI = %v, want %v", config.Spotify.RedirectURI, "http://127.0.0.1:8080/callback")
	}
}

func TestIsConfigured(t *testing.T) {
	// Reset viper
	viper.Reset()

	tests := []struct {
		name     string
		envVars  map[string]string
		expected bool
	}{
		{
			name: "properly configured",
			envVars: map[string]string{
				"SPOTIFY_CLIENT_ID":     "real_client_id",
				"SPOTIFY_CLIENT_SECRET": "real_client_secret",
			},
			expected: true,
		},
		{
			name: "placeholder values",
			envVars: map[string]string{
				"SPOTIFY_CLIENT_ID":     "your_spotify_client_id",
				"SPOTIFY_CLIENT_SECRET": "your_spotify_client_secret",
			},
			expected: false,
		},
		{
			name:     "empty values",
			envVars:  map[string]string{},
			expected: false,
		},
		{
			name: "missing client secret",
			envVars: map[string]string{
				"SPOTIFY_CLIENT_ID": "real_client_id",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset viper and environment
			viper.Reset()
			os.Unsetenv("SPOTIFY_CLIENT_ID")
			os.Unsetenv("SPOTIFY_CLIENT_SECRET")
			os.Unsetenv("SPOTIFY_REDIRECT_URI")

			// Set environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Use a temporary directory for config
			tempDir := t.TempDir()
			SetConfigPaths(tempDir)

			err := ReadConfig()
			if err != nil {
				t.Fatalf("ReadConfig() error = %v", err)
			}

			result := IsConfigured()
			if result != tt.expected {
				t.Errorf("IsConfigured() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSetSpotifyConfig(t *testing.T) {
	// Reset configuration
	viper.Reset()
	tempDir := t.TempDir()
	SetConfigPaths(tempDir)
	ReadConfig()

	clientID := "test_client_id"
	clientSecret := "test_client_secret"
	redirectURI := "http://127.0.0.1:9999/callback"

	SetSpotifyConfig(clientID, clientSecret, redirectURI)

	spotify := GetSpotify()
	if spotify.ClientID != clientID {
		t.Errorf("ClientID = %v, want %v", spotify.ClientID, clientID)
	}
	if spotify.ClientSecret != clientSecret {
		t.Errorf("ClientSecret = %v, want %v", spotify.ClientSecret, clientSecret)
	}
	if spotify.RedirectURI != redirectURI {
		t.Errorf("RedirectURI = %v, want %v", spotify.RedirectURI, redirectURI)
	}
}

func TestSaveConfig(t *testing.T) {
	// Reset configuration
	viper.Reset()
	tempDir := t.TempDir()
	SetConfigPaths(tempDir)
	ReadConfig()

	// Mock home directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	clientID := "test_client_id"
	clientSecret := "test_client_secret"
	redirectURI := "http://127.0.0.1:9999/callback"

	SetSpotifyConfig(clientID, clientSecret, redirectURI)

	err := SaveConfig()
	if err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Check that config file was created
	configPath := filepath.Join(tempDir, ".spotify-shuffle.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created at %s", configPath)
	}

	// Read and verify config file content
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	configStr := string(content)
	if !strings.Contains(configStr, clientID) {
		t.Errorf("Config file does not contain client ID")
	}
	if !strings.Contains(configStr, clientSecret) {
		t.Errorf("Config file does not contain client secret")
	}
	if !strings.Contains(configStr, redirectURI) {
		t.Errorf("Config file does not contain redirect URI")
	}
}
