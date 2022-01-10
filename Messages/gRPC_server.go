package main

import (
	"Messages/pb"
	"context"
	"net"

	"github.com/gocql/gocql"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCServer struct {
	db *gocql.Session
	pb.UnimplementedMessageMethodsServer
}

func NewGrpcServer(session *gocql.Session) *gRPCServer {
	return &gRPCServer{db: session}
}

func (g *gRPCServer) Send(ctx context.Context, in *pb.Message) (*empty.Empty, error) {
	uuid, _ := gocql.RandomUUID()
	user_uuid, err := gocql.ParseUUID(in.User)
	if err != nil {
		return nil, nil
	}
	err = g.db.Query(`INSERT INTO messages (id, content, user_id) VALUES (?, ?, ?)`,
		uuid,
		in.Content,
		user_uuid,
	).Exec()

	return &empty.Empty{}, err
}

func StartGrpc(session *gocql.Session) {
	lstn, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic(err)
	}

	serverMethods := NewGrpcServer(session)
	s := grpc.NewServer()

	reflection.Register(s)
	pb.RegisterMessageMethodsServer(s, serverMethods)

	s.Serve(lstn)
}
