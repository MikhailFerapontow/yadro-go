all:
	go build -o xkcd ./cmd/xkcd/main.go

install:
	go mod tidy && go mod download

test:
	go test ./...

bench:
	go

clean:
	rm xkcd
	rm database.json
	rm index.json