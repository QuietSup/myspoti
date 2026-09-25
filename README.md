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

4. Install the CLI and put Go's bin dir on your `PATH`:

```bash
go install .
```

Add this to `~/.zshrc` (or `~/.bashrc`):

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

Then reload your shell (`source ~/.zshrc`).

5. Log in (from anywhere):

```bash
myspoti auth
```

Tokens are stored in `~/.config/myspoti/tokens.json`.  
Playback control needs Spotify Premium and an active device (open Spotify somewhere first).  
Re-run `myspoti auth` after pulling new scopes (e.g. library like/unlike, recently played).

## Usage

```bash
myspoti auth    # log in
myspoti now     # currently playing track
myspoti queue   # recent, current, and up next
myspoti play    # resume
myspoti pause   # pause
myspoti next    # skip forward
myspoti prev    # skip back
myspoti like    # save current track
myspoti unlike  # remove current track from library
myspoti help
```

While developing in the repo you can also use `go run . <command>` or `make build` → `./myspoti`.  
Re-run `go install .` after code changes to refresh the global binary.

## Development

```bash
make fmt
make lint       # requires golangci-lint on PATH
make test
make build      # → ./myspoti
```

Install the linter:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```
