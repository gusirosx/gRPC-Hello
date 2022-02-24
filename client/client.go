package main

import (
	"fmt"
	"log"
	"net/http"

	pb "gRPC-gin/proto"

	"github.com/gin-gonic/gin"
)

func main() {

	// Set up a http server.
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		fmt.Fprintln(ctx.Writer, "Up and running...")
	})
	router.GET("/hello/:name", Hello)
	router.GET("/bye/:name", Bye)
	// Run http server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func Hello(ctx *gin.Context) {
	// Set up a connection to the server hello
	conn, err := Connection(server, portHello)
	if err != nil {
		log.Printf("failed to dial server %s: %v", *server, err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	name := ctx.Param("name")
	// Contact the server and print out its response.
	req := &pb.HelloRequest{Name: name}
	res, err := client.SayHello(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprint(res.Message),
	})
}

func Bye(ctx *gin.Context) {
	// Set up a connection to the server hello
	conn, err := Connection(server, portBye)
	if err != nil {
		log.Printf("failed to dial server %s: %v", *server, err)
	}
	defer conn.Close()
	client := pb.NewByeClient(conn)

	name := ctx.Param("name")
	// Contact the server and print out its response.
	req := &pb.HelloRequest{Name: name}
	res, err := client.SayGoodbye(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": fmt.Sprint(res.Message),
	})
}
