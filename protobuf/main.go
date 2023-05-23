package main

import (
	"fmt"
	"log"

	"github.com/tokatu4561/grpc-demo/protobuf/pb"
	"google.golang.org/protobuf/proto"
)


func main () {
	person := pb.Person{
		Id: 1,
		Name: "tokatu",
		Email: "",
		Phone: pb.PhoneType_HOME,
		Friends: []string{"tokatu", "tokatu2"},
		Languages: map[string]string{
			"ja": "日本語",
		},
	}

	// シリアライズ
	binData, err := proto.Marshal(&person)
	if err != nil {
		log.Fatalln("failed to marshal person:", err)
	}

	// デシリアライズ
	var person2 pb.Person
	if err := proto.Unmarshal(binData, &person2); err != nil {
		log.Fatalln("failed to unmarshal person:", err)
	}

	fmt.Println(person2)
}