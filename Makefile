server:
	go run cmd/server/main.go

client:
	go run cmd/client/main.go

echo-client:
	go run cmd/echo-client/main.go


.PHONY: server client echo-client