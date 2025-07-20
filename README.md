# Spotify Shuffle - Go Edition

A fast, cross-platform CLI tool for managing your Spotify playlists. Built in Go for optimal performance and single-binary distribution.

## Features

- 🔀 **Shuffle playlist** - Randomize track order
- 🔤 **Sort playlist** - Sort by title or artist name  
- 🔄 **Reverse playlist** - Reverse current order
- 🗑️ **Remove tracks** - By age or artist name
- ➕ **Create playlists** - Fresh (recent tracks), Chunk (split large playlists), Genre-based
- ⚡ **Fast execution** - Compiled Go binary
- 🌍 **Cross-platform** - Windows, macOS, Linux
- 📦 **Single binary** - No runtime dependencies

## Quick Start

### 1. Download Binary

Download the appropriate binary for your platform from the [releases page](https://github.com/petabloc/spotify-shuffle/releases):

#### Binaries:
- **Windows**: `spotify-shuffle-windows-amd64.exe`
- **macOS Intel**: `spotify-shuffle-macos-amd64`
- **macOS Apple Silicon**: `spotify-shuffle-macos-arm64`
- **Linux**: `spotify-shuffle-linux-amd64`

#### Packages:
- **macOS**: `spotify-shuffle-macos-amd64.dmg` (drag & drop installer)
- **Windows**: `spotify-shuffle-windows-amd64.exe.msi` (Windows installer)
- **Debian/Ubuntu**: `spotify-shuffle-linux-amd64.deb` (APT package)

### 2. Setup Spotify App

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Add `http://127.0.0.1:8080/callback` to redirect URIs
4. Copy your Client ID and Client Secret

### 3. Configure Credentials

**Option A: Config File**
The app will create `~/.spotify-shuffle.yaml` on first run:

```yaml
spotify:
  client_id: "your_spotify_client_id"
  client_secret: "your_spotify_client_secret"
  redirect_uri: "http://127.0.0.1:8080/callback"
```

**Option B: Environment Variables**
```bash
export SPOTIFY_CLIENT_ID="your_client_id"
export SPOTIFY_CLIENT_SECRET="your_client_secret"
export SPOTIFY_REDIRECT_URI="http://127.0.0.1:8080/callback"
```

## Usage

### Basic Commands

```bash
# Shuffle a playlist
./spotify-shuffle shuffle --playlist 37i9dQZF1DXcBWIGoYBM5M

# Sort by title
./spotify-shuffle sort --by title --playlist 37i9dQZF1DXcBWIGoYBM5M

# Sort by artist
./spotify-shuffle sort --by artist --playlist 37i9dQZF1DXcBWIGoYBM5M

# Reverse order
./spotify-shuffle reverse --playlist 37i9dQZF1DXcBWIGoYBM5M

# Remove tracks older than 90 days
./spotify-shuffle remove --age --days 90 --playlist 37i9dQZF1DXcBWIGoYBM5M

# Remove tracks by artist (interactive)
./spotify-shuffle remove --artist --playlist 37i9dQZF1DXcBWIGoYBM5M

# Remove specific artist
./spotify-shuffle remove --artist --name "Artist Name" --playlist 37i9dQZF1DXcBWIGoYBM5M

# Create fresh playlist with tracks from last 90 days
./spotify-shuffle create --type fresh --days 90 --name "Recent Hits" --playlist 37i9dQZF1DXcBWIGoYBM5M

# Create chunk playlists (250 tracks each)
./spotify-shuffle create --type chunk --name "BigPlaylist" --size 250 --playlist 37i9dQZF1DXcBWIGoYBM5M

# Create genre playlist (interactive mode)
./spotify-shuffle create --type genre --interactive --playlist 37i9dQZF1DXcBWIGoYBM5M

# Create genre playlist (direct)
./spotify-shuffle create --type genre --genre "rock" --name "Rock Collection" --playlist 37i9dQZF1DXcBWIGoYBM5M
```

### Getting Playlist ID

From Spotify:
- Right-click playlist → Share → Copy link
- URL: `https://open.spotify.com/playlist/37i9dQZF1DXcBWIGoYBM5M`
- Use full URL or just the ID: `37i9dQZF1DXcBWIGoYBM5M`

## Building from Source

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- Git

### Build

```bash
# Clone repository
git clone https://github.com/petabloc/spotify-shuffle.git
cd spotify-shuffle

# Build for current platform
make build

# Build for all platforms
make build-all

# Install to system
make install

# Create packages
make package-all
```

### Development

```bash
# Install dependencies
make deps

# Build with race detection
make dev

# Run tests
make test

# Format and check code
make check

# See all targets
make help
```

## Authentication

On first run, the app will:

1. Display an authorization URL
2. Open your browser to authenticate
3. Show a success page after authorization
4. Save credentials for future use

The authentication flow uses a local server on port 8080 to handle the OAuth callback.

## Advantages over Python Version

| Feature | Go Version | Python Version |
|---------|------------|----------------|
| **Binary Size** | ~15MB | ~80MB (with PyInstaller) |
| **Startup Time** | Instant | 2-3 seconds |
| **Dependencies** | None | Python + venv |
| **Distribution** | Single file | Multiple files |
| **Performance** | Very fast | Moderate |
| **Memory Usage** | Low | Higher |

## Release & Distribution

The project includes automated GitHub Actions for:

### Continuous Integration
- ✅ **Testing** - Automated tests on all platforms
- ✅ **Code Quality** - Format and vet checks
- ✅ **Cross-compilation** - Build verification for all targets

### Release Automation
When a new tag is pushed (e.g., `v1.0.0`):
- ✅ **Binaries** - Built for all platforms automatically
- ✅ **macOS DMG** - Drag & drop installer package
- ✅ **Windows MSI** - Windows installer package
- ✅ **Debian DEB** - APT-compatible package
- ✅ **GitHub Release** - Automatic release with all assets

## Platform Support

- ✅ **Windows** (amd64)
- ✅ **macOS** (Intel + Apple Silicon)
- ✅ **Linux** (amd64 + arm64)
- ✅ **WSL** (Windows Subsystem for Linux)

## Troubleshooting

**Authentication Issues:**
- Ensure redirect URI is exactly: `http://127.0.0.1:8080/callback`
- Check that credentials are correctly configured
- Try deleting config file to reset authentication

**Permission Errors:**
- You can only modify playlists you own or follow
- Ensure your Spotify app has the correct scopes

**Binary Issues:**
- On macOS: `chmod +x spotify-shuffle-macos-*`
- On Linux: `chmod +x spotify-shuffle-linux-*`

## License

MIT License - See LICENSE file for details.