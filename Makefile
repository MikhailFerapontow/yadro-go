all:
	go build -o myapp ./cmd/xkcd/main.go

install:
	go mod download

test:
	go test ./...

clean:
	rm myapp