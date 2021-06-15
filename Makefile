server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go

echo-client:
	go run cmd/echo-client/main.go

certs:
	mkdir -p certs/
	openssl genrsa -out certs/server.key 2048
	openssl req -new -x509 -sha256 -key certs/server.key -out certs/server.crt -days 3650

build:
	go build -o bin/server.exe cmd/server/main.go

clean:
	rm -rf bin/

.PHONY: server client echo-client certs build clean