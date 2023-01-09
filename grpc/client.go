package main

import (
	"context"
	"fmt"
	"os"

	"github.com/fiuskylab/grpc-forward/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT")), opts...)
	if err != nil {
		panic(err)
	}

	client := proto.NewMessageClient(conn)

	empty := &proto.Empty{}

	stream, err := client.Subscribe(context.Background(), empty)
	if err != nil {
		panic(err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			panic(err)
		}
		fmt.Println("ID: %d | Message: %s", msg.Id, msg.Msg)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			panic(err)
		}
	}()
}
