generate:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/helloworld.proto

server:
	go run grpc/server.go
client:
	go run gin/main.go