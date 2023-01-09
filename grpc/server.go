package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/fiuskylab/grpc-forward/proto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type MessageServer struct {
	proto.UnimplementedMessageServer
	messageChan    chan *proto.MessageResponse
	streams        map[string]proto.Message_SubscribeServer
	streamMessages map[string]chan *proto.MessageResponse
}

func (m *MessageServer) Subscribe(_ *proto.Empty, stream proto.Message_SubscribeServer) error {
	var err error = nil

	id := uuid.NewString()

	m.streams[id] = stream
	m.streamMessages[id] = make(chan *proto.MessageResponse)

	for err == nil {
		msg := <-m.streamMessages[id]
		if err = stream.Send(msg); err != nil {
			fmt.Println("AAAAA")
			break
		}
	}
	fmt.Println("BBBBBB")
	delete(m.streams, id)
	delete(m.streamMessages, id)
	return err
}

func (m *MessageServer) StartForward() {
	for {
		msg, ok := <-m.messageChan
		if !ok {
			panic("stream not ok")
		}
		for _, v := range m.streamMessages {
			v <- msg
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
		messageChan:    make(chan *proto.MessageResponse, 1),
		streams:        make(map[string]proto.Message_SubscribeServer),
		streamMessages: make(map[string]chan *proto.MessageResponse),
	}

	go m.StartForward()

	proto.RegisterMessageServer(grpcServer, m)

	go func() {
		for i := 0; i < 300; i++ {
			fmt.Println("List of all IDs: ", getKeys(m.streams))
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

func getKeys[K comparable, V any](m map[K]V) []K {
	keys := []K{}
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
