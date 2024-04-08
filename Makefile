all:
	go build -o xkcd ./cmd/xkcd/main.go

install:
	go mod download

test:
	go test ./...

clean:
	rm xkcd
	rm database.json