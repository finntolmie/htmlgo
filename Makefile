run:
	@go build -o bin/htmlgo && ./bin/htmlgo

test:
	@go test ./...