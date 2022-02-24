package main

import (
	"log"
	"net"

	pb "gRPC-gin/proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	// Embed the unimplemented server
	pb.UnimplementedByeServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayGoodbye(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Goodbye " + in.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50040")
	if err != nil {
		log.Fatalf("failed to listen on port 50040: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterByeServer(srv, &server{})
	log.Printf("server listening at %v", lis.Addr())
	// Register reflection service on gRPC server.
	reflection.Register(srv)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
