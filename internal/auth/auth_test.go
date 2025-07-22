package auth

import (
	"context"
	"testing"
	"time"
)

func TestNewSpotifyAuth(t *testing.T) {
	clientID := "test_client_id"
	clientSecret := "test_client_secret"
	redirectURI := "http://127.0.0.1:8080/callback"

	auth := NewSpotifyAuth(clientID, clientSecret, redirectURI)

	if auth == nil {
		t.Fatal("NewSpotifyAuth() returned nil")
	}

	if auth.clientID != clientID {
		t.Errorf("clientID = %v, want %v", auth.clientID, clientID)
	}

	if auth.clientSecret != clientSecret {
		t.Errorf("clientSecret = %v, want %v", auth.clientSecret, clientSecret)
	}

	if auth.redirectURI != redirectURI {
		t.Errorf("redirectURI = %v, want %v", auth.redirectURI, redirectURI)
	}

	if auth.state == "" {
		t.Error("state should not be empty")
	}

	if auth.tokenFile == "" {
		t.Error("tokenFile should not be empty")
	}

	if auth.auth == nil {
		t.Error("auth should not be nil")
	}
}

func TestSpotifyAuth_GetClient_NoSavedToken(t *testing.T) {
	// This test verifies that GetClient fails when no token is saved
	// We use a very short timeout to avoid actually starting the server
	auth := NewSpotifyAuth("test_id", "test_secret", "http://127.0.0.1:8080/callback")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// This should fail quickly due to timeout, avoiding the server startup
	_, err := auth.GetClient(ctx)
	if err == nil {
		t.Error("Expected error when no saved token and context times out quickly")
	}
}

func TestSpotifyAuth_loadToken(t *testing.T) {
	auth := NewSpotifyAuth("test_id", "test_secret", "http://127.0.0.1:8080/callback")

	// Current implementation always returns error
	_, err := auth.loadToken()
	if err == nil {
		t.Error("Expected error from loadToken() in current implementation")
	}
}

func TestSpotifyAuth_saveToken(t *testing.T) {
	auth := NewSpotifyAuth("test_id", "test_secret", "http://127.0.0.1:8080/callback")

	// Current implementation doesn't actually save, just logs
	err := auth.saveToken(nil)
	if err != nil {
		t.Errorf("saveToken() returned unexpected error: %v", err)
	}
}

func TestSpotifyAuth_StateGeneration(t *testing.T) {
	auth1 := NewSpotifyAuth("test_id", "test_secret", "http://127.0.0.1:8080/callback")
	auth2 := NewSpotifyAuth("test_id", "test_secret", "http://127.0.0.1:8080/callback")

	if auth1.state == auth2.state {
		t.Error("State should be different for different auth instances")
	}

	if len(auth1.state) == 0 {
		t.Error("State should not be empty")
	}
}
