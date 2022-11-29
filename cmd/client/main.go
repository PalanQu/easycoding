package main

import (
	"context"
	"fmt"
	"log"
	"time"

	ping_pb "easycoding/api/ping"
	pkg_otel "easycoding/pkg/otel"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	address := "localhost:10001"
	tracer, otelShutdownFunc, err := pkg_otel.NewTracer()
	if err != nil {
		panic(err)
	}
	_, span := tracer.Start(context.Background(), "ping client")
	defer span.End()
	defer otelShutdownFunc()
	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
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
