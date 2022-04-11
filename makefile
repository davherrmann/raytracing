sinclude .env

dev:
	find . -name "*.go" -o -name "*.html" | entr -rc go run ./... -port ${PORT}

dev-race:
	go run -race ./... -port ${PORT}