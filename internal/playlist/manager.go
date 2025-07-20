package playlist

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/zmb3/spotify/v2"
)

type Manager struct {
	client *spotify.Client
}

// NewManager creates a new playlist manager
func NewManager(client *spotify.Client) *Manager {
	return &Manager{client: client}
}

// Track represents a track with metadata
type Track struct {
	ID      spotify.ID
	Name    string
	Artists []string
	URI     spotify.URI
	AddedAt time.Time
}

// GetPlaylistTracks retrieves all tracks from a playlist
func (m *Manager) GetPlaylistTracks(ctx context.Context, playlistID spotify.ID) ([]Track, error) {
	var tracks []Track
	limit := 50
	offset := 0

	for {
		page, err := m.client.GetPlaylistTracks(ctx, playlistID, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return nil, fmt.Errorf("failed to get playlist tracks: %w", err)
		}

		for _, item := range page.Tracks {
			if item.Track.ID == "" {
				continue
			}

			track := item.Track
			var artistNames []string
			for _, artist := range track.Artists {
				artistNames = append(artistNames, artist.Name)
			}

			addedAt := time.Time{}
			if item.AddedAt != "" {
				if parsed, err := time.Parse(time.RFC3339, item.AddedAt); err == nil {
					addedAt = parsed
				}
			}

			tracks = append(tracks, Track{
				ID:      track.ID,
				Name:    track.Name,
				Artists: artistNames,
				URI:     track.URI,
				AddedAt: addedAt,
			})
		}

		if len(page.Tracks) < limit {
			break
		}
		offset += limit
	}

	return tracks, nil
}

// ShufflePlaylist randomizes the order of tracks in a playlist
func (m *Manager) ShufflePlaylist(ctx context.Context, playlistID spotify.ID) error {
	tracks, err := m.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return fmt.Errorf("playlist is empty")
	}

	// Extract URIs and shuffle
	uris := make([]spotify.URI, len(tracks))
	for i, track := range tracks {
		uris[i] = track.URI
	}
	rand.Shuffle(len(uris), func(i, j int) {
		uris[i], uris[j] = uris[j], uris[i]
	})

	// Replace playlist with shuffled tracks
	return m.replacePlaylistTracks(ctx, playlistID, uris)
}

// SortPlaylist sorts playlist tracks by the specified criteria
func (m *Manager) SortPlaylist(ctx context.Context, playlistID spotify.ID, sortBy string) error {
	switch sortBy {
	case "title":
		return m.SortPlaylistByTitle(ctx, playlistID)
	case "artist":
		return m.SortPlaylistByArtist(ctx, playlistID)
	default:
		return fmt.Errorf("invalid sort criteria: %s. Use 'title' or 'artist'", sortBy)
	}
}

// SortPlaylistByTitle sorts playlist tracks alphabetically by title
func (m *Manager) SortPlaylistByTitle(ctx context.Context, playlistID spotify.ID) error {
	tracks, err := m.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return fmt.Errorf("playlist is empty")
	}

	// Sort by title (case insensitive)
	sort.Slice(tracks, func(i, j int) bool {
		return strings.ToLower(tracks[i].Name) < strings.ToLower(tracks[j].Name)
	})

	// Extract URIs
	uris := make([]spotify.URI, len(tracks))
	for i, track := range tracks {
		uris[i] = track.URI
	}

	// Replace playlist with sorted tracks
	return m.replacePlaylistTracks(ctx, playlistID, uris)
}

// SortPlaylistByArtist sorts playlist tracks alphabetically by artist
func (m *Manager) SortPlaylistByArtist(ctx context.Context, playlistID spotify.ID) error {
	tracks, err := m.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return fmt.Errorf("playlist is empty")
	}

	// Sort by first artist (case insensitive)
	sort.Slice(tracks, func(i, j int) bool {
		artistI := ""
		artistJ := ""
		if len(tracks[i].Artists) > 0 {
			artistI = strings.ToLower(tracks[i].Artists[0])
		}
		if len(tracks[j].Artists) > 0 {
			artistJ = strings.ToLower(tracks[j].Artists[0])
		}
		return artistI < artistJ
	})

	// Extract URIs
	uris := make([]spotify.URI, len(tracks))
	for i, track := range tracks {
		uris[i] = track.URI
	}

	// Replace playlist with sorted tracks
	return m.replacePlaylistTracks(ctx, playlistID, uris)
}

// ReversePlaylist reverses the order of tracks in a playlist
func (m *Manager) ReversePlaylist(ctx context.Context, playlistID spotify.ID) error {
	tracks, err := m.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return err
	}

	if len(tracks) == 0 {
		return fmt.Errorf("playlist is empty")
	}

	// Reverse the tracks
	uris := make([]spotify.URI, len(tracks))
	for i, track := range tracks {
		uris[len(tracks)-1-i] = track.URI
	}

	// Replace playlist with reversed tracks
	return m.replacePlaylistTracks(ctx, playlistID, uris)
}

// RemoveOldTracks removes tracks older than specified days
func (m *Manager) RemoveOldTracks(ctx context.Context, playlistID spotify.ID, days int) (int, error) {
	tracks, err := m.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return 0, err
	}

	cutoffDate := time.Now().AddDate(0, 0, -days)
	var tracksToKeep []spotify.URI

	for _, track := range tracks {
		if track.AddedAt.IsZero() || track.AddedAt.After(cutoffDate) {
			tracksToKeep = append(tracksToKeep, track.URI)
		}
	}

	removedCount := len(tracks) - len(tracksToKeep)
	if removedCount == 0 {
		return 0, nil
	}

	// Replace playlist with tracks to keep
	err = m.replacePlaylistTracks(ctx, playlistID, tracksToKeep)
	return removedCount, err
}

// RemoveTracksByArtist removes all tracks by a specific artist
func (m *Manager) RemoveTracksByArtist(ctx context.Context, playlistID spotify.ID, artistName string) (int, error) {
	tracks, err := m.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return 0, err
	}

	var tracksToKeep []spotify.URI
	artistLower := strings.ToLower(artistName)

	for _, track := range tracks {
		keepTrack := true
		for _, artist := range track.Artists {
			if strings.Contains(strings.ToLower(artist), artistLower) {
				keepTrack = false
				break
			}
		}
		if keepTrack {
			tracksToKeep = append(tracksToKeep, track.URI)
		}
	}

	removedCount := len(tracks) - len(tracksToKeep)
	if removedCount == 0 {
		return 0, nil
	}

	// Replace playlist with tracks to keep
	err = m.replacePlaylistTracks(ctx, playlistID, tracksToKeep)
	return removedCount, err
}

// replacePlaylistTracks replaces all tracks in a playlist with new ones
func (m *Manager) replacePlaylistTracks(ctx context.Context, playlistID spotify.ID, uris []spotify.URI) error {
	if len(uris) == 0 {
		// Clear playlist
		return m.client.ReplacePlaylistTracks(ctx, playlistID)
	}

	// Spotify API limit is 100 tracks per request
	const batchSize = 100

	// Replace first batch
	firstBatch := uris
	if len(uris) > batchSize {
		firstBatch = uris[:batchSize]
	}

	// Convert URIs to IDs for ReplacePlaylistTracks
	var firstBatchIDs []spotify.ID
	for _, uri := range firstBatch {
		firstBatchIDs = append(firstBatchIDs, spotify.ID(strings.TrimPrefix(string(uri), "spotify:track:")))
	}

	if err := m.client.ReplacePlaylistTracks(ctx, playlistID, firstBatchIDs...); err != nil {
		return fmt.Errorf("failed to replace playlist tracks: %w", err)
	}

	// Add remaining batches
	for i := batchSize; i < len(uris); i += batchSize {
		end := i + batchSize
		if end > len(uris) {
			end = len(uris)
		}
		batch := uris[i:end]

		// Convert URIs to IDs for AddTracksToPlaylist
		var batchIDs []spotify.ID
		for _, uri := range batch {
			batchIDs = append(batchIDs, spotify.ID(strings.TrimPrefix(string(uri), "spotify:track:")))
		}

		_, err := m.client.AddTracksToPlaylist(ctx, playlistID, batchIDs...)
		if err != nil {
			return fmt.Errorf("failed to add tracks to playlist: %w", err)
		}
	}

	return nil
}

// GetUniqueArtists returns all unique artists in a playlist
func (m *Manager) GetUniqueArtists(ctx context.Context, playlistID spotify.ID) ([]string, error) {
	tracks, err := m.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	artistSet := make(map[string]bool)
	for _, track := range tracks {
		for _, artist := range track.Artists {
			artistSet[artist] = true
		}
	}

	var artists []string
	for artist := range artistSet {
		artists = append(artists, artist)
	}

	sort.Strings(artists)
	return artists, nil
}

// CreatePlaylist creates a new playlist
func (m *Manager) CreatePlaylist(ctx context.Context, name, description string, public bool) (spotify.ID, error) {
	user, err := m.client.CurrentUser(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	playlist, err := m.client.CreatePlaylistForUser(ctx, user.ID, name, description, public, false)
	if err != nil {
		return "", fmt.Errorf("failed to create playlist: %w", err)
	}

	return playlist.ID, nil
}

// FindPlaylistByName finds a playlist by name (case insensitive)
func (m *Manager) FindPlaylistByName(ctx context.Context, name string) (spotify.ID, error) {
	limit := 50
	offset := 0
	nameLower := strings.ToLower(name)

	for {
		playlists, err := m.client.CurrentUsersPlaylists(ctx, spotify.Limit(limit), spotify.Offset(offset))
		if err != nil {
			return "", fmt.Errorf("failed to get playlists: %w", err)
		}

		for _, playlist := range playlists.Playlists {
			if strings.ToLower(playlist.Name) == nameLower {
				return playlist.ID, nil
			}
		}

		if len(playlists.Playlists) < limit {
			break
		}
		offset += limit
	}

	return "", fmt.Errorf("playlist not found")
}

// CreateFreshPlaylist creates a playlist with tracks added within the last N days
func (m *Manager) CreateFreshPlaylist(ctx context.Context, sourcePlaylistID spotify.ID, name string, days int, overwrite bool) (int, error) {
	tracks, err := m.GetPlaylistTracks(ctx, sourcePlaylistID)
	if err != nil {
		return 0, err
	}

	// Filter tracks by date
	cutoffDate := time.Now().AddDate(0, 0, -days)
	var freshTracks []spotify.URI

	for _, track := range tracks {
		if !track.AddedAt.IsZero() && track.AddedAt.After(cutoffDate) {
			freshTracks = append(freshTracks, track.URI)
		}
	}

	if len(freshTracks) == 0 {
		return 0, fmt.Errorf("no tracks found within the last %d days", days)
	}

	// Check if playlist exists
	var playlistID spotify.ID
	if existingID, err := m.FindPlaylistByName(ctx, name); err == nil {
		if !overwrite {
			return 0, fmt.Errorf("playlist '%s' already exists", name)
		}
		playlistID = existingID
		// Clear existing playlist
		if err := m.client.ReplacePlaylistTracks(ctx, playlistID); err != nil {
			return 0, fmt.Errorf("failed to clear existing playlist: %w", err)
		}
	} else {
		// Create new playlist
		description := fmt.Sprintf("Fresh tracks from the last %d days", days)
		newID, err := m.CreatePlaylist(ctx, name, description, false)
		if err != nil {
			return 0, err
		}
		playlistID = newID
	}

	// Add tracks to playlist
	if err := m.replacePlaylistTracks(ctx, playlistID, freshTracks); err != nil {
		return 0, fmt.Errorf("failed to add tracks to playlist: %w", err)
	}

	return len(freshTracks), nil
}

// CreateChunkPlaylists creates multiple playlists with random chunks of tracks
func (m *Manager) CreateChunkPlaylists(ctx context.Context, sourcePlaylistID spotify.ID, baseName string, chunkSize int, overwrite bool) (int, error) {
	tracks, err := m.GetPlaylistTracks(ctx, sourcePlaylistID)
	if err != nil {
		return 0, err
	}

	if len(tracks) == 0 {
		return 0, fmt.Errorf("source playlist is empty")
	}

	// Extract URIs and shuffle
	uris := make([]spotify.URI, len(tracks))
	for i, track := range tracks {
		uris[i] = track.URI
	}
	rand.Shuffle(len(uris), func(i, j int) {
		uris[i], uris[j] = uris[j], uris[i]
	})

	// Calculate number of chunks
	totalChunks := (len(uris) + chunkSize - 1) / chunkSize
	createdCount := 0

	for chunkNum := 0; chunkNum < totalChunks; chunkNum++ {
		start := chunkNum * chunkSize
		end := start + chunkSize
		if end > len(uris) {
			end = len(uris)
		}
		chunkTracks := uris[start:end]

		// Create playlist name
		chunkName := fmt.Sprintf("%s-%02d", baseName, chunkNum)

		// Check if playlist exists
		var playlistID spotify.ID
		if existingID, err := m.FindPlaylistByName(ctx, chunkName); err == nil {
			if !overwrite {
				continue // Skip existing playlists if not overwriting
			}
			playlistID = existingID
			// Clear existing playlist
			if err := m.client.ReplacePlaylistTracks(ctx, playlistID); err != nil {
				continue
			}
		} else {
			// Create new playlist
			description := fmt.Sprintf("Chunk %d of %d from %s", chunkNum+1, totalChunks, baseName)
			newID, err := m.CreatePlaylist(ctx, chunkName, description, false)
			if err != nil {
				continue
			}
			playlistID = newID
		}

		// Add tracks to playlist
		if err := m.replacePlaylistTracks(ctx, playlistID, chunkTracks); err != nil {
			continue
		}

		createdCount++
	}

	return createdCount, nil
}

// GetPlaylistGenres gets all genres in a playlist with track counts
func (m *Manager) GetPlaylistGenres(ctx context.Context, playlistID spotify.ID) (map[string]int, error) {
	tracks, err := m.GetPlaylistTracks(ctx, playlistID)
	if err != nil {
		return nil, err
	}

	if len(tracks) == 0 {
		return nil, fmt.Errorf("playlist is empty")
	}

	// Get track IDs
	var trackIDs []spotify.ID
	for _, track := range tracks {
		trackIDs = append(trackIDs, track.ID)
	}

	// Get track genres by looking up artists
	trackGenres, err := m.getTrackGenres(ctx, trackIDs)
	if err != nil {
		return nil, err
	}

	// Count genre occurrences
	genreCounts := make(map[string]int)
	for _, genres := range trackGenres {
		for _, genre := range genres {
			genreCounts[genre]++
		}
	}

	return genreCounts, nil
}

// CreateGenrePlaylist creates a playlist with tracks from a specific genre
func (m *Manager) CreateGenrePlaylist(ctx context.Context, sourcePlaylistID spotify.ID, name, targetGenre string, overwrite bool) (int, error) {
	tracks, err := m.GetPlaylistTracks(ctx, sourcePlaylistID)
	if err != nil {
		return 0, err
	}

	if len(tracks) == 0 {
		return 0, fmt.Errorf("source playlist is empty")
	}

	// Get track IDs
	var trackIDs []spotify.ID
	for _, track := range tracks {
		trackIDs = append(trackIDs, track.ID)
	}

	// Get track genres
	trackGenres, err := m.getTrackGenres(ctx, trackIDs)
	if err != nil {
		return 0, err
	}

	// Filter tracks by genre (case-insensitive partial match)
	var genreTracks []spotify.URI
	targetGenreLower := strings.ToLower(targetGenre)

	for _, track := range tracks {
		if genres, exists := trackGenres[track.ID]; exists {
			for _, genre := range genres {
				if strings.Contains(strings.ToLower(genre), targetGenreLower) {
					genreTracks = append(genreTracks, track.URI)
					break
				}
			}
		}
	}

	if len(genreTracks) == 0 {
		return 0, fmt.Errorf("no tracks found for genre '%s'", targetGenre)
	}

	// Check if playlist exists
	var playlistID spotify.ID
	if existingID, err := m.FindPlaylistByName(ctx, name); err == nil {
		if !overwrite {
			return 0, fmt.Errorf("playlist '%s' already exists", name)
		}
		playlistID = existingID
		// Clear existing playlist
		if err := m.client.ReplacePlaylistTracks(ctx, playlistID); err != nil {
			return 0, fmt.Errorf("failed to clear existing playlist: %w", err)
		}
	} else {
		// Create new playlist
		description := fmt.Sprintf("Tracks with genre: %s", targetGenre)
		newID, err := m.CreatePlaylist(ctx, name, description, false)
		if err != nil {
			return 0, err
		}
		playlistID = newID
	}

	// Add tracks to playlist
	if err := m.replacePlaylistTracks(ctx, playlistID, genreTracks); err != nil {
		return 0, fmt.Errorf("failed to add tracks to playlist: %w", err)
	}

	return len(genreTracks), nil
}

// getTrackGenres gets genres for tracks by looking up their artists
func (m *Manager) getTrackGenres(ctx context.Context, trackIDs []spotify.ID) (map[spotify.ID][]string, error) {
	trackGenres := make(map[spotify.ID][]string)
	batchSize := 50

	// Get unique artist IDs from tracks
	artistIDs := make(map[spotify.ID]bool)
	trackArtists := make(map[spotify.ID][]spotify.ID)

	// Process tracks in batches
	for i := 0; i < len(trackIDs); i += batchSize {
		end := i + batchSize
		if end > len(trackIDs) {
			end = len(trackIDs)
		}
		batch := trackIDs[i:end]

		tracks, err := m.client.GetTracks(ctx, batch)
		if err != nil {
			return nil, fmt.Errorf("failed to get tracks: %w", err)
		}

		for j, track := range tracks {
			if track == nil {
				continue
			}
			trackID := batch[j]
			var artists []spotify.ID
			for _, artist := range track.Artists {
				artists = append(artists, artist.ID)
				artistIDs[artist.ID] = true
			}
			trackArtists[trackID] = artists
		}
	}

	// Get genres for all unique artists
	artistGenres := make(map[spotify.ID][]string)
	var artistIDsList []spotify.ID
	for id := range artistIDs {
		artistIDsList = append(artistIDsList, id)
	}

	for i := 0; i < len(artistIDsList); i += batchSize {
		end := i + batchSize
		if end > len(artistIDsList) {
			end = len(artistIDsList)
		}
		batch := artistIDsList[i:end]

		artists, err := m.client.GetArtists(ctx, batch...)
		if err != nil {
			return nil, fmt.Errorf("failed to get artists: %w", err)
		}

		for j, artist := range artists {
			if artist != nil {
				artistGenres[batch[j]] = artist.Genres
			}
		}
	}

	// Map track IDs to genres
	for trackID, artists := range trackArtists {
		var genres []string
		genreSet := make(map[string]bool)
		for _, artistID := range artists {
			if artistGenreList, exists := artistGenres[artistID]; exists {
				for _, genre := range artistGenreList {
					if !genreSet[genre] {
						genres = append(genres, genre)
						genreSet[genre] = true
					}
				}
			}
		}
		trackGenres[trackID] = genres
	}

	return trackGenres, nil
}
