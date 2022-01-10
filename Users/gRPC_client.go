package main

import (
	"SampleTraceScylla/pb"
	"fmt"

	"google.golang.org/grpc"
)

func ConnectGrpc() pb.MessageMethodsClient {
	conn, err := grpc.Dial("host.docker.internal:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Cannot connect to gRPC server")
		panic(err)
	}

	msg := pb.NewMessageMethodsClient(conn)
	return msg
}
