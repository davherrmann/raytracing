dev:
	find . -name "*.go" -o -name "*.html" | entr -rc go run ./...