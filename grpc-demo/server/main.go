package main

import (
	"github.com/tokatu4561/grpc-demo/grpc-demo/pb"
)

type server struct {
	pb.UnimplementedFileServiceServer
}