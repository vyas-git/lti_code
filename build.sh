GOOS=windows GOARCH=amd64 go build -o bin/spotify-win-amd64.exe
GOOS=windows GOARCH=386  go build -o bin/spotify-win-amd32.exe  

GOOS=darwin GOARCH=amd64 go build -o bin/spotify-macos-amd64
GOOS=darwin GOARCH=386 go build -o bin/spotify-macos-amd32


GOOS=linux GOARCH=amd64 go build -o bin/spotify-linux-amd64
GOOS=linux GOARCH=386 go build -o bin/spotify-linux-amd32

