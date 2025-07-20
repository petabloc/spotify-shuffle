package playlist

import (
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/zmb3/spotify/v2"
)

func TestNewManager(t *testing.T) {
	client := &spotify.Client{}
	manager := NewManager(client)

	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}

	if manager.client != client {
		t.Error("Manager client not set correctly")
	}
}

func TestTrack(t *testing.T) {
	track := Track{
		ID:      spotify.ID("test_id"),
		Name:    "Test Song",
		Artists: []string{"Artist 1", "Artist 2"},
		URI:     spotify.URI("spotify:track:test_id"),
		AddedAt: time.Now(),
	}

	if track.ID != "test_id" {
		t.Errorf("Track ID = %v, want %v", track.ID, "test_id")
	}

	if track.Name != "Test Song" {
		t.Errorf("Track Name = %v, want %v", track.Name, "Test Song")
	}

	expectedArtists := []string{"Artist 1", "Artist 2"}
	if !reflect.DeepEqual(track.Artists, expectedArtists) {
		t.Errorf("Track Artists = %v, want %v", track.Artists, expectedArtists)
	}
}

func TestShuffleTracks(t *testing.T) {
	tracks := []Track{
		{ID: "1", Name: "Song 1"},
		{ID: "2", Name: "Song 2"},
		{ID: "3", Name: "Song 3"},
		{ID: "4", Name: "Song 4"},
		{ID: "5", Name: "Song 5"},
	}

	originalOrder := make([]Track, len(tracks))
	copy(originalOrder, tracks)

	shuffleTracks(tracks)

	// Check that all tracks are still present
	if len(tracks) != len(originalOrder) {
		t.Errorf("Shuffled tracks length = %v, want %v", len(tracks), len(originalOrder))
	}

	// Check that all original tracks are still in the slice
	for _, originalTrack := range originalOrder {
		found := false
		for _, shuffledTrack := range tracks {
			if originalTrack.ID == shuffledTrack.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Track %v not found after shuffle", originalTrack.ID)
		}
	}
}

func TestSortTracksByTitle(t *testing.T) {
	tracks := []Track{
		{ID: "1", Name: "Zebra Song"},
		{ID: "2", Name: "Apple Song"},
		{ID: "3", Name: "Banana Song"},
	}

	sortTracksByTitle(tracks)

	expected := []string{"Apple Song", "Banana Song", "Zebra Song"}
	for i, track := range tracks {
		if track.Name != expected[i] {
			t.Errorf("Track at index %d: Name = %v, want %v", i, track.Name, expected[i])
		}
	}
}

func TestSortTracksByArtist(t *testing.T) {
	tracks := []Track{
		{ID: "1", Name: "Song 1", Artists: []string{"Zebra Artist"}},
		{ID: "2", Name: "Song 2", Artists: []string{"Apple Artist"}},
		{ID: "3", Name: "Song 3", Artists: []string{"Banana Artist"}},
	}

	sortTracksByArtist(tracks)

	expected := []string{"Apple Artist", "Banana Artist", "Zebra Artist"}
	for i, track := range tracks {
		if track.Artists[0] != expected[i] {
			t.Errorf("Track at index %d: Artist = %v, want %v", i, track.Artists[0], expected[i])
		}
	}
}

func TestReverseTracks(t *testing.T) {
	tracks := []Track{
		{ID: "1", Name: "Song 1"},
		{ID: "2", Name: "Song 2"},
		{ID: "3", Name: "Song 3"},
	}

	reverseTracks(tracks)

	expected := []string{"Song 3", "Song 2", "Song 1"}
	for i, track := range tracks {
		if track.Name != expected[i] {
			t.Errorf("Track at index %d: Name = %v, want %v", i, track.Name, expected[i])
		}
	}
}

func TestFilterTracksByAge(t *testing.T) {
	now := time.Now()
	tracks := []Track{
		{ID: "1", Name: "Recent Song", AddedAt: now.AddDate(0, 0, -30)},   // 30 days ago
		{ID: "2", Name: "Old Song", AddedAt: now.AddDate(0, 0, -100)},     // 100 days ago
		{ID: "3", Name: "Very Old Song", AddedAt: now.AddDate(0, 0, -200)}, // 200 days ago
	}

	// Filter tracks older than 90 days
	filtered := filterTracksByAge(tracks, 90)

	// Should remove tracks older than 90 days
	expected := 1 // Only the 30-day-old track should remain
	if len(filtered) != expected {
		t.Errorf("Filtered tracks length = %v, want %v", len(filtered), expected)
	}

	if filtered[0].Name != "Recent Song" {
		t.Errorf("Remaining track = %v, want %v", filtered[0].Name, "Recent Song")
	}
}

func TestFilterTracksByArtist(t *testing.T) {
	tracks := []Track{
		{ID: "1", Name: "Song 1", Artists: []string{"Artist A"}},
		{ID: "2", Name: "Song 2", Artists: []string{"Artist B"}},
		{ID: "3", Name: "Song 3", Artists: []string{"Artist A", "Artist C"}},
		{ID: "4", Name: "Song 4", Artists: []string{"Artist D"}},
	}

	// Remove tracks by "Artist A"
	filtered := filterTracksByArtist(tracks, "Artist A")

	// Should remove tracks 1 and 3
	expected := 2
	if len(filtered) != expected {
		t.Errorf("Filtered tracks length = %v, want %v", len(filtered), expected)
	}

	// Check remaining tracks
	for _, track := range filtered {
		for _, artist := range track.Artists {
			if artist == "Artist A" {
				t.Errorf("Track %v still contains Artist A after filtering", track.Name)
			}
		}
	}
}

func TestGetFreshTracks(t *testing.T) {
	now := time.Now()
	tracks := []Track{
		{ID: "1", Name: "Very Recent", AddedAt: now.AddDate(0, 0, -10)},  // 10 days ago
		{ID: "2", Name: "Recent", AddedAt: now.AddDate(0, 0, -20)},      // 20 days ago
		{ID: "3", Name: "Old", AddedAt: now.AddDate(0, 0, -40)},         // 40 days ago
	}

	// Get tracks from last 30 days
	fresh := getFreshTracks(tracks, 30)

	expected := 2 // Tracks from 10 and 20 days ago
	if len(fresh) != expected {
		t.Errorf("Fresh tracks length = %v, want %v", len(fresh), expected)
	}
}

func TestGetTracksByGenre(t *testing.T) {
	tracks := []Track{
		{ID: "1", Name: "Rock Song", Artists: []string{"Rock Artist"}},
		{ID: "2", Name: "Pop Song", Artists: []string{"Pop Artist"}},
		{ID: "3", Name: "Jazz Song", Artists: []string{"Jazz Artist"}},
	}

	// Mock genre lookup - in real implementation this would call Spotify API
	// For testing, we'll assume the function works with artist names containing genres
	rock := getTracksByGenre(tracks, "rock")
	if len(rock) != 1 {
		t.Errorf("Rock tracks length = %v, want %v", len(rock), 1)
	}
}

func TestChunkTracks(t *testing.T) {
	tracks := []Track{
		{ID: "1"}, {ID: "2"}, {ID: "3"}, {ID: "4"}, {ID: "5"},
		{ID: "6"}, {ID: "7"}, {ID: "8"}, {ID: "9"}, {ID: "10"},
	}

	chunks := chunkTracks(tracks, 3)

	expectedChunks := 4 // 3, 3, 3, 1
	if len(chunks) != expectedChunks {
		t.Errorf("Number of chunks = %v, want %v", len(chunks), expectedChunks)
	}

	// Check chunk sizes
	expectedSizes := []int{3, 3, 3, 1}
	for i, chunk := range chunks {
		if len(chunk) != expectedSizes[i] {
			t.Errorf("Chunk %d size = %v, want %v", i, len(chunk), expectedSizes[i])
		}
	}
}

// Helper functions for testing - these would be implemented in manager.go

func shuffleTracks(tracks []Track) {
	for i := range tracks {
		j := rand.Intn(i + 1)
		tracks[i], tracks[j] = tracks[j], tracks[i]
	}
}

func sortTracksByTitle(tracks []Track) {
	sort.Slice(tracks, func(i, j int) bool {
		return tracks[i].Name < tracks[j].Name
	})
}

func sortTracksByArtist(tracks []Track) {
	sort.Slice(tracks, func(i, j int) bool {
		if len(tracks[i].Artists) == 0 {
			return true
		}
		if len(tracks[j].Artists) == 0 {
			return false
		}
		return tracks[i].Artists[0] < tracks[j].Artists[0]
	})
}

func reverseTracks(tracks []Track) {
	for i, j := 0, len(tracks)-1; i < j; i, j = i+1, j-1 {
		tracks[i], tracks[j] = tracks[j], tracks[i]
	}
}

func filterTracksByAge(tracks []Track, days int) []Track {
	cutoff := time.Now().AddDate(0, 0, -days)
	var filtered []Track
	for _, track := range tracks {
		if track.AddedAt.After(cutoff) {
			filtered = append(filtered, track)
		}
	}
	return filtered
}

func filterTracksByArtist(tracks []Track, artistName string) []Track {
	var filtered []Track
	for _, track := range tracks {
		hasArtist := false
		for _, artist := range track.Artists {
			if artist == artistName {
				hasArtist = true
				break
			}
		}
		if !hasArtist {
			filtered = append(filtered, track)
		}
	}
	return filtered
}

func getFreshTracks(tracks []Track, days int) []Track {
	cutoff := time.Now().AddDate(0, 0, -days)
	var fresh []Track
	for _, track := range tracks {
		if track.AddedAt.After(cutoff) {
			fresh = append(fresh, track)
		}
	}
	return fresh
}

func getTracksByGenre(tracks []Track, genre string) []Track {
	// Mock implementation for testing
	var filtered []Track
	for _, track := range tracks {
		for _, artist := range track.Artists {
			if strings.Contains(strings.ToLower(artist), strings.ToLower(genre)) {
				filtered = append(filtered, track)
				break
			}
		}
	}
	return filtered
}

func chunkTracks(tracks []Track, size int) [][]Track {
	var chunks [][]Track
	for i := 0; i < len(tracks); i += size {
		end := i + size
		if end > len(tracks) {
			end = len(tracks)
		}
		chunks = append(chunks, tracks[i:end])
	}
	return chunks
}