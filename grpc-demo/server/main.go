package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"

	"github.com/tokatu4561/grpc-demo/grpc-demo/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

// implement ListFiles method
func (s *server) ListFile(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFIlesResponse, error) {
	log.Printf("ListFiles called")

	fileDirPath := "storage"

	paths, err := ioutil.ReadDir(fileDirPath)
	if err != nil {
		return nil, err
	}

	var filenames []string
	for _, path := range paths {
		// ファイルであれば
		if !path.IsDir() {
			filenames = append(filenames, path.Name())
		}
	}

	// ファイル名一覧を返す
	return &pb.ListFIlesResponse{
		Filenames: filenames,
	}, nil
}

func main() {
	// create server
	s := grpc.NewServer()

	// register service
	pb.RegisterFileServiceServer(s, &server{})

	// listen port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v¥n", err)
	}

	// start server
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v¥n", err)
	}
}