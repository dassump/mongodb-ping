CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags "-s -w" -o mongodb-ping-linux-amd64         .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o mongodb-ping-windows-amd64.exe .
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags "-s -w" -o mongodb-ping-darwin-amd64       .
CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -ldflags "-s -w" -o mongodb-ping-darwin-arm64       .