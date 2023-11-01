# HC-dl

A performant CLI downloader for hentai-cosplays.com.

## Usage

```sh
hc-dl -u https://hentai-cosplays.com/...

# optionally set how may concurrent download perform
# Example: perform 4 concurrent download
# Default: perform [CPU count] concurrent download
hc-dl -u https://hentai-cosplays.com/... -c 4
```

## Build
```sh
go build cmd/main.go -o hc-dl
```