build:
	go build -o bin/haste ./main.go

test:
	go test -v ./...
