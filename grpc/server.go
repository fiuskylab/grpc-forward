package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/fiuskylab/grpc-forward/proto"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type MessageServer struct {
	proto.UnimplementedMessageServer
	messageChan chan *proto.MessageResponse
}

func (m *MessageServer) Subscribe(_ *proto.Empty, stream proto.Message_SubscribeServer) error {
	for {
		msg, ok := <-m.messageChan
		if ok {
			if err := stream.Send(msg); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("chan not ok")
		}
	}
}

func main() {
	fmt.Println("Hallo")

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", os.Getenv("GRPC_PORT")))
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	m := &MessageServer{
		messageChan: make(chan *proto.MessageResponse, 1),
	}

	proto.RegisterMessageServer(grpcServer, m)

	go func() {
		for i := 0; i < 300; i++ {
			m.messageChan <- &proto.MessageResponse{
				Id:  uint32(i),
				Msg: fmt.Sprintf("Message %d", i),
			}
			time.Sleep(time.Second)
		}
	}()

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
