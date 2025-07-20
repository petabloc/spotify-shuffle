# Spotify Shuffle - Go Edition

A fast, cross-platform CLI tool for managing your Spotify playlists. Built in Go for optimal performance and single-binary distribution.

🎯 **New: Interactive Mode** - The easiest way to manage your playlists with a guided, menu-driven interface!

## 🚀 30-Second Quick Start

```bash
# 1. Download binary for your platform (see below)
# 2. Run interactive mode
./spotify-shuffle                    # Linux/macOS
spotify-shuffle.exe                  # Windows

# 3. Follow the guided setup:
#    → Enter Spotify app credentials (one-time)
#    → Browse and select your playlists
#    → Choose operations with guided menus
#    → Confirm changes with safety prompts
```

**That's it!** No command-line syntax to learn, no playlist IDs to copy. The interactive mode handles everything with a user-friendly interface.

## Features

- 🎯 **Interactive Mode** - Guided interface for all operations
- 🔀 **Shuffle playlist** - Randomize track order
- 🔤 **Sort playlist** - Sort by title or artist name  
- 🔄 **Reverse playlist** - Reverse current order
- 🗑️ **Remove tracks** - By age or artist name
- ➕ **Create playlists** - Fresh (recent tracks), Chunk (split large playlists), Genre-based
- 📋 **Playlist selection** - Browse your playlists or enter ID/URL
- ⚡ **Fast execution** - Compiled Go binary
- 🌍 **Cross-platform** - Windows, macOS, Linux
- 📦 **Single binary** - No runtime dependencies

## Quick Start

### Option 1: Interactive Mode (Recommended for Beginners)

1. Download the binary for your platform (see below)
2. Set up your Spotify app credentials
3. Run `./spotify-shuffle` and follow the guided interface!

### Option 2: Command Line Mode (For Automation)

### 1. Download Binary

Download the appropriate binary for your platform from the [releases page](https://github.com/petabloc/spotify-shuffle/releases):

#### Binaries:
- **Windows**: `spotify-shuffle-windows-amd64.exe` (Works on Windows 10/11)
- **macOS Intel**: `spotify-shuffle-macos-amd64` (Intel Macs)
- **macOS Apple Silicon**: `spotify-shuffle-macos-arm64` (M1/M2/M3 Macs)
- **Linux x64**: `spotify-shuffle-linux-amd64` (Most Linux distributions)
- **Linux ARM**: `spotify-shuffle-linux-arm64` (Raspberry Pi, ARM servers)

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

**Option A: Interactive Setup (Easiest)**
The app will guide you through credential setup on first run and create `~/.spotify-shuffle.yaml`

**Option B: Manual Config File**
```yaml
spotify:
  client_id: "your_spotify_client_id"
  client_secret: "your_spotify_client_secret"
  redirect_uri: "http://127.0.0.1:8080/callback"
```

**Option C: Environment Variables**
```bash
export SPOTIFY_CLIENT_ID="your_client_id"
export SPOTIFY_CLIENT_SECRET="your_client_secret"
export SPOTIFY_REDIRECT_URI="http://127.0.0.1:8080/callback"
```

### 4. Run Interactive Mode

```bash
# Launch the guided interface
./spotify-shuffle

# Or explicitly
./spotify-shuffle interactive
```

## Usage

### Interactive Mode (Recommended)

The interactive mode provides a user-friendly, guided experience:

```bash
# Launch interactive mode
./spotify-shuffle interactive

# Or simply run without arguments
./spotify-shuffle
```

#### 🎯 Interactive Features:

**Playlist Selection:**
- 📋 **Browse your playlists** - Automatically loads and displays your Spotify playlists
- 🔗 **Manual entry** - Enter playlist ID or full Spotify URL
- 🔍 **Smart search** - Find playlists by name

**Operations Menu:**
- 🔀 **Shuffle** - Randomize track order with confirmation
- 🔤 **Sort** - Choose between title or artist sorting
- 🔄 **Reverse** - Reverse current playlist order
- 🗑️ **Remove tracks** - Multiple removal options:
  - By age: 90/180 days, 1/2/3 years, or custom
  - By artist: Interactive artist selection
- ➕ **Create playlists**:
  - **Fresh**: Recent tracks (30/90/180 days)
  - **Chunk**: Split large playlists (custom size)
  - **Genre**: Filter by music genre
- ℹ️ **Playlist info** - View statistics and top artists

**User Experience:**
- ✅ **Safety prompts** - Confirmation before any changes
- 📊 **Real-time feedback** - See results immediately
- 🔄 **Multi-session** - Work with multiple playlists
- 🎨 **Rich interface** - Emojis and clear formatting

#### Example Interactive Session:

```
🎵 Welcome to Spotify Shuffle Interactive Mode!
===============================================

📋 Select a playlist:
1. Enter playlist ID/URL manually
2. Choose from your playlists
3. Exit

Choose option (1-3): 2

🔍 Loading your playlists...

👤 Found 15 playlists for John Doe:
 1. My Liked Songs (1,234 tracks)
 2. Workout Mix (67 tracks)
 3. Chill Vibes (145 tracks)
 ...

Choose playlist (1-15): 2

🎵 Working with playlist: Workout Mix
==========================================

📋 Choose an operation:
1. 🔀 Shuffle tracks
2. 🔤 Sort tracks
3. 🔄 Reverse tracks
...
```

### Command Line Mode (For Automation & Scripts)

For automation, scripts, or when you prefer command-line arguments:

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

**Interactive Mode**: Automatically browses your playlists - no ID needed!

**Command Line Mode**: Get playlist ID from Spotify:
- Right-click playlist → Share → Copy link
- URL: `https://open.spotify.com/playlist/37i9dQZF1DXcBWIGoYBM5M`
- Use full URL or just the ID: `37i9dQZF1DXcBWIGoYBM5M`

**Supported Formats**:
- Full URL: `https://open.spotify.com/playlist/37i9dQZF1DXcBWIGoYBM5M`
- Spotify URI: `spotify:playlist:37i9dQZF1DXcBWIGoYBM5M`
- Direct ID: `37i9dQZF1DXcBWIGoYBM5M`

## Example Workflows

### 🎯 Beginner Workflow (Interactive Mode)

Perfect for first-time users or occasional playlist management:

```bash
# 1. Launch interactive mode
./spotify-shuffle

# The app guides you through:
# → Authentication (one-time setup)
# → Playlist selection (browse your playlists)
# → Operation selection (guided menus)
# → Confirmation prompts (safety checks)
# → Real-time feedback (see results)
```

**Example Session:**
1. **First Run**: Set up Spotify credentials
2. **Select Playlist**: Browse and choose "My Workout Mix"
3. **Choose Operation**: Shuffle tracks
4. **Confirm**: "This will shuffle 67 tracks. Continue? (y/N)"
5. **Result**: "✅ Playlist shuffled successfully!"
6. **Continue**: Work with another playlist or exit

### 🔧 Power User Workflow (Command Line)

Perfect for automation, scripts, and batch operations:

```bash
# Playlist management script
PLAYLIST="37i9dQZF1DXcBWIGoYBM5M"

# 1. Clean up old tracks
./spotify-shuffle remove --age --days 180 --playlist $PLAYLIST

# 2. Add fresh content
./spotify-shuffle create --type fresh --days 30 --name "Recent Hits" --playlist $PLAYLIST

# 3. Organize main playlist
./spotify-shuffle shuffle --playlist $PLAYLIST
```

### 🎵 Music Curation Workflow

Comprehensive playlist management using both modes:

```bash
# 1. Interactive discovery phase
./spotify-shuffle interactive
# → Browse playlists
# → Check playlist statistics
# → Identify playlists needing attention

# 2. Batch processing phase (command line)
for playlist in "Rock Classics" "Pop Hits" "Indie Discoveries"; do
  ./spotify-shuffle remove --age --days 365 --playlist "$playlist"
  ./spotify-shuffle shuffle --playlist "$playlist"
done

# 3. Create specialized playlists
./spotify-shuffle create --type genre --genre "electronic" --name "Electronic Vibes" --playlist "Main Mix"
./spotify-shuffle create --type chunk --size 200 --name "Road Trip" --playlist "Long Playlist"
```

### 🤖 Automation Workflow

Setting up automated playlist maintenance:

```bash
#!/bin/bash
# weekly-playlist-maintenance.sh

# Clean up workout playlist
./spotify-shuffle remove --age --days 90 --playlist "Workout Mix"

# Refresh discovery playlist
./spotify-shuffle create --type fresh --days 7 --name "This Week's Finds" --playlist "Discovery Weekly" --overwrite

# Shuffle main playlists
./spotify-shuffle shuffle --playlist "Daily Mix"
./spotify-shuffle shuffle --playlist "Liked Songs Sample"

# Create weekend playlist
./spotify-shuffle create --type chunk --size 100 --name "Weekend Vibes" --playlist "Chill Collection"
```

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

### First-Time Setup

On first run, the app will:

1. **Interactive Mode**: Guide you through credential setup
2. **All Modes**: Display an authorization URL
3. **Browser Authentication**: Automatically open your browser
4. **Success Confirmation**: Show a success page after authorization
5. **Save Credentials**: Store tokens for future use

### Technical Details

- **OAuth 2.0 Flow**: Secure authentication via Spotify
- **Local Server**: Uses port 8080 for OAuth callback
- **Token Storage**: Encrypted token storage for security
- **Auto-Refresh**: Automatically refreshes expired tokens

### Security Features

- ✅ **No password storage** - Uses OAuth tokens only
- ✅ **Local authentication** - No data sent to third parties
- ✅ **Secure redirect** - Uses IP-based localhost (127.0.0.1)
- ✅ **Minimal permissions** - Only requests necessary Spotify scopes

## Interactive vs Command Line Comparison

| Feature | Interactive Mode | Command Line Mode |
|---------|------------------|-------------------|
| **🎯 Best For** | Beginners, exploration, one-off tasks | Automation, scripts, power users |
| **📋 Playlist Selection** | Browse your playlists visually | Requires playlist ID/URL |
| **🛡️ Safety** | Confirmation prompts for all actions | Direct execution (be careful!) |
| **🎨 User Experience** | Rich, guided interface with emojis | Fast, scriptable commands |
| **📊 Feedback** | Real-time statistics and progress | Minimal output |
| **🔄 Workflow** | Multi-playlist sessions | Single operation per command |
| **📚 Learning Curve** | None - guided menus | Requires learning command syntax |

## Advantages over Python Version

| Feature | Go Version | Python Version |
|---------|------------|----------------|
| **Binary Size** | ~9MB | ~80MB (with PyInstaller) |
| **Startup Time** | Instant | 2-3 seconds |
| **Dependencies** | None | Python + venv |
| **Distribution** | Single file | Multiple files |
| **Performance** | Very fast | Moderate |
| **Memory Usage** | Low | Higher |
| **Interactive Mode** | ✅ Full-featured | ❌ Not available |
| **Cross-Platform** | ✅ All platforms | ✅ Limited packaging |

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

### Common Issues

**🔐 Authentication Problems:**
- **Redirect URI Error**: Ensure redirect URI is exactly `http://127.0.0.1:8080/callback`
- **Invalid Credentials**: Check Client ID/Secret in Spotify Developer Dashboard
- **Reset Authentication**: Delete `~/.spotify-shuffle.yaml` and `~/.spotify-shuffle-token.json`
- **Port Conflicts**: Ensure port 8080 is available

**📋 Playlist Issues:**
- **Permission Denied**: You can only modify playlists you own or follow
- **Playlist Not Found**: Check playlist ID/URL is correct and playlist is public/accessible
- **Missing Tracks**: Some tracks may not be available due to regional restrictions

**💻 System Issues:**
- **macOS Permission**: Run `chmod +x spotify-shuffle-macos-*`
- **Linux Permission**: Run `chmod +x spotify-shuffle-linux-*`
- **Windows**: Use `spotify-shuffle-windows-amd64.exe` (not the Linux binary!)
- **Windows Antivirus**: Add binary to antivirus exceptions if needed
- **Windows PowerShell**: Use `.\spotify-shuffle-windows-amd64.exe` if in PowerShell

**🎯 Interactive Mode Issues:**
- **No Playlists Shown**: Check Spotify account has playlists and proper authentication
- **Selection Not Working**: Use number keys (1-9) for menu selection
- **Browser Not Opening**: Copy the displayed URL manually into your browser

### Getting Help

- **Interactive Mode**: Built-in help and error messages
- **Command Help**: Use `--help` flag with any command
- **Verbose Output**: Use interactive mode for detailed feedback
- **GitHub Issues**: Report bugs at [github.com/petabloc/spotify-shuffle/issues](https://github.com/petabloc/spotify-shuffle/issues)

## License

MIT License - See LICENSE file for details.