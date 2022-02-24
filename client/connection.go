package main

import (
	"crypto/tls"
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	portHello    = flag.String("server hello port", ":50050", "Server address (port)")
	portBye      = flag.String("server bye port", ":50040", "Server address (port)")
	server       = flag.String("server-host", "localhost", "Host name to which server IP should resolve")
	insecureFlag = flag.Bool("insecure", true, "Skip SSL validation? [false]")
	skipVerify   = flag.Bool("skip-verify", false, "Skip server hostname verification in SSL validation [false]")
)

func init() {
	flag.Parse()
}

// Connection creates a new gRPC connection to the server.
// host should be of the form domain:port, e.g., example.com:443
func Connection(host, port *string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if *host != "" {
		opts = append(opts, grpc.WithAuthority(*host))
	} else {
		log.Fatal("Server address is empty")
	}
	if *insecureFlag {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		cred := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: *skipVerify,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	return grpc.Dial(*host+*port, opts...)
}
