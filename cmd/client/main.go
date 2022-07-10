package main

import (
	"context"
	ping_pb "easycoding/api/ping"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := "localhost:10001"
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := ping_pb.NewPingSvcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.Ping(ctx, &ping_pb.PingRequest{Req: "ping"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r)
}
