server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go

echo-client:
	go run cmd/echo-client/main.go

build:
	go build -o bin/server.exe cmd/server/main.go

clean:
	rm -rf bin/

.PHONY: server client echo-client build clean