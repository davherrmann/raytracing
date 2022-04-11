sinclude .env

dev:
	find . -name "*.go" -o -name "*.html" | entr -rc go run cmd/rtgo/main.go -port ${PORT}

dev-race:
	go run -race cmd/rtgo/main.go -port ${PORT}