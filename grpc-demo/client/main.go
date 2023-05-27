package main

import (
	"context"
	"io"
	"log"

	"github.com/tokatu4561/grpc-demo/grpc-demo/pb"
	"google.golang.org/grpc"
)

func main() {
	// create connection WithInsecure() はTLSを使わない 本番で使うべきではない
	con, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v¥n", err)
	}
	defer con.Close()

	// create client
	client := pb.NewFileServiceClient(con)
	// call ListFiles
	callListFiles(client)
	// call Download
	callDownload(client)
}

func callListFiles(client pb.FileServiceClient) {
	// call ListFiles
	res, err := client.ListFile(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalf("failed to ListFiles: %v¥n", err)
	}

	log.Printf("filenames: %v¥n", res.GetFilenames())
}

func callDownload(client pb.FileServiceClient) {
	// call Download
	stream, err := client.Download(context.Background(), &pb.DownloadFileRequest{
		Filename: "test.txt",
	})
	if err != nil {
		log.Fatalf("failed to Download: %v¥n", err)
	}

	// ファイルを受信して標準出力に書き出す
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("failed to receive: %v¥n", err)
		}

		log.Printf("received: %v¥n", res.GetContent())
	}
}