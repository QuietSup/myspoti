# myspoti

Control Spotify from the terminal.

## Setup

1. Create an app in the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard).
2. In the app settings, register redirect URI: `http://127.0.0.1:8080/callback`
3. Save your Client ID:

```bash
mkdir -p ~/.config/myspoti
echo '{"client_id":"YOUR_CLIENT_ID"}' > ~/.config/myspoti/config.json
```

Or export `SPOTIFY_CLIENT_ID` (overrides the config file).

1. Log in:

```bash
go run . auth
```

Tokens are stored in `~/.config/myspoti/tokens.json`.  
Playback control needs Spotify Premium and an active device (open Spotify somewhere first).

## Usage

```bash
go run . auth    # log in (re-run after pulling new scopes)
go run . now     # currently playing track
go run . play    # resume
go run . pause   # pause
go run . next    # skip forward
go run . prev    # skip back
go run . like    # save current track
go run . unlike  # remove current track from library
go run . help
```

Build a binary:

```bash
make build      # → ./myspoti
./myspoti now
```



## Development

```bash
make fmt
make lint       # requires golangci-lint on PATH
make test
```

Install the linter:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

