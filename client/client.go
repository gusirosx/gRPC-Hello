package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"

	pb "gRPC-gin/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	insc "google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("server", "localhost:50050", "Server address (host:port)")
	serverHost = flag.String("server-host", "localhost", "Host name to which server IP should resolve")
	insecure   = flag.Bool("insecure", true, "Skip SSL validation? [false]")
	skipVerify = flag.Bool("skip-verify", false, "Skip server hostname verification in SSL validation [false]")
)

func init() {
	flag.Parse()
}

// Connection creates a new gRPC connection to the server.
// host should be of the form domain:port, e.g., example.com:443
func Connection() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if *serverAddr == "" {
		log.Fatal("-server is empty")
	}
	if *serverHost != "" {
		opts = append(opts, grpc.WithAuthority(*serverHost))
	}
	if *insecure {
		opts = append(opts, grpc.WithTransportCredentials(insc.NewCredentials()))
	} else {
		cred := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: *skipVerify,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}
	return grpc.Dial(*serverAddr, opts...)
}

func main() {
	// Set up a connection to the server.
	conn, err := Connection()
	if err != nil {
		log.Printf("failed to dial server %s: %v", *serverAddr, err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// Set up a http server.
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		fmt.Fprintln(ctx.Writer, "Up and running...")
	})

	router.GET("/name/:name", func(c *gin.Context) {
		name := c.Param("name")

		// Contact the server and print out its response.
		req := &pb.HelloRequest{Name: name}
		res, err := client.SayHello(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(res.Message),
		})
	})
	// Run http server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
