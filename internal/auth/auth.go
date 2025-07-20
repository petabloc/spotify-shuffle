package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type SpotifyAuth struct {
	auth         *spotifyauth.Authenticator
	state        string
	tokenFile    string
	redirectURI  string
	clientID     string
	clientSecret string
}

// NewSpotifyAuth creates a new Spotify authenticator
func NewSpotifyAuth(clientID, clientSecret, redirectURI string) *SpotifyAuth {
	home, _ := os.UserHomeDir()
	tokenFile := filepath.Join(home, ".spotify-shuffle-token.json")

	auth := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistReadPrivate,
		),
		spotifyauth.WithClientID(clientID),
		spotifyauth.WithClientSecret(clientSecret),
	)

	// Generate random state
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	return &SpotifyAuth{
		auth:         auth,
		state:        state,
		tokenFile:    tokenFile,
		redirectURI:  redirectURI,
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// GetClient returns an authenticated Spotify client
func (sa *SpotifyAuth) GetClient(ctx context.Context) (*spotify.Client, error) {
	// Try to load existing token
	if token, err := sa.loadToken(); err == nil {
		client := spotify.New(sa.auth.Client(ctx, token))
		return client, nil
	}

	// Need to authenticate
	return sa.authenticate(ctx)
}

// authenticate performs the OAuth flow
func (sa *SpotifyAuth) authenticate(ctx context.Context) (*spotify.Client, error) {
	// Start local server to handle callback
	ch := make(chan *oauth2.Token)
	errCh := make(chan error)

	// Create HTTP server
	server := &http.Server{Addr: ":8080"}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		token, err := sa.auth.Token(ctx, sa.state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusForbidden)
			errCh <- err
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<html>
			<head><title>Spotify Authentication</title></head>
			<body style="font-family: Arial, sans-serif; text-align: center; padding: 50px;">
				<h1 style="color: #1DB954;">✅ Authentication Successful!</h1>
				<p>You can now close this window and return to the terminal.</p>
			</body>
			</html>
		`))

		ch <- token
	})

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Get auth URL and display to user
	url := sa.auth.AuthURL(sa.state)
	fmt.Printf("\n🔐 Please open this URL in your browser to authenticate:\n%s\n\n", url)
	fmt.Println("Waiting for authentication...")

	// Wait for token or timeout
	var token *oauth2.Token
	select {
	case token = <-ch:
		fmt.Println("✅ Authentication successful!")
	case err := <-errCh:
		return nil, fmt.Errorf("authentication error: %w", err)
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("authentication timeout")
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	// Save token
	if err := sa.saveToken(token); err != nil {
		log.Printf("Warning: failed to save token: %v", err)
	}

	// Create client
	client := spotify.New(sa.auth.Client(context.Background(), token))
	return client, nil
}

// loadToken loads a saved token from file
func (sa *SpotifyAuth) loadToken() (*oauth2.Token, error) {
	// For now, return error to force new authentication
	// In a real implementation, you'd load from a file
	return nil, fmt.Errorf("no saved token")
}

// saveToken saves a token to file
func (sa *SpotifyAuth) saveToken(token *oauth2.Token) error {
	// For now, just log that we would save it
	log.Printf("Token would be saved to: %s", sa.tokenFile)
	// In a real implementation, you'd save to a file
	return nil
}
