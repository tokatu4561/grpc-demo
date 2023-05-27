package main

import (
	"context"
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
	callListFiles(client)
}

func callListFiles(client pb.FileServiceClient) {
	// call ListFiles
	res, err := client.ListFile(context.Background(), &pb.ListFilesRequest{})
	if err != nil {
		log.Fatalf("failed to ListFiles: %v¥n", err)
	}

	log.Printf("filenames: %v¥n", res.GetFilenames())
}