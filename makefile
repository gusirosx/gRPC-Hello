generate:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/hello.proto

run_server:
	@echo "---- Running Server ----"
	@go run server/server.go

run_client:
	@echo "---- Running Client ----"
	@go run client/client.go