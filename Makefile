all:
	go build -o myapp ./cmd/xkcd/main.go

test:
	go test ./...

clean:
	rm myapp