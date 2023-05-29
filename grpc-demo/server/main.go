package main

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/tokatu4561/grpc-demo/grpc-demo/pb"
	"google.golang.org/grpc"

	// grpc middleware
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	// auth
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

type server struct {
	pb.UnimplementedFileServiceServer
}

// リクエストとレスポンスの前後でログを挟む interceptor
func myLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{},error){
		// リクエストデータ
		log.Printf("request data: %v", req)
		resp, err := handler(ctx, req)
		if err != nil { 
			return nil, err
		}
		// レスポンスデータ
		log.Printf("response data: %v", resp)

		return resp, nil
	}
}

func authorize(ctx context.Context) (context.Context, error) {
	// 認証処理
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	// test でなければ認証エラー
	if token != "test" {
		return nil, errors.New("invalid token")
	}

	return ctx, nil
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
	s := grpc.NewServer(grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(
			// ログ出力
			myLogging(),
			// 認証
			grpc_auth.UnaryServerInterceptor(authorize),
		),
	))

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