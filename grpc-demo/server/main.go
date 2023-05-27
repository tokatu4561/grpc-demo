package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

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

func (s *server) Download(req *pb.DownloadFileRequest, stream pb.FileService_DownloadServer ) error{
	filename := req.GetFilename()
	filePath := "storage/" + filename 

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	// ファイルを読み込んでstreamに書き込む 5バイトずつ
	buf := make([]byte, 5)
	for {
		n, err := file.Read(buf)
		// ファイルの終端に到達したら終了
		if n == 0 || err == io.EOF{
			break
		}

		if err != nil {
			return err
		}
		
		// streamに書き込む
		if err := stream.Send(&pb.DownloadFileResponse{
			Content: buf[:n],
		}); err != nil {
			return err
		}

		// すぐ終了してしまうので、一時的にスリープ
		time.Sleep(1 * time.Second)
	}

	return nil
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