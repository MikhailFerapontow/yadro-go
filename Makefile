all:
	go build -o xkcd ./cmd/xkcd/main.go

server:
	go build -o xkcd-server ./cmd/xkcd/main.go

install:
	go mod tidy && go mod download

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm xkcd
	rm database.json
	rm index.json