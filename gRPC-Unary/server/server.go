package main

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/ryantokmanmokmtm/gRPC-Unary/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	ADDR = "localhost:50001"
)

type greetingService struct {
	pb.UnimplementedGreetingServer
}

func (s *greetingService) SayHello(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	name := req.YourName
	if name == "" {
		return nil, errors.New("required a name")
	}

	return &pb.GreetingResponse{
		Result: fmt.Sprintf("Hello,%v", name),
	}, nil
}

func (s *greetingService) SumUp(ctx context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {
	sumA := req.SumA
	sumB := req.SumB

	return &pb.SumResponse{Result: sumA + sumB}, nil
}

func main() {
	//create tcp connection
	tcp, err := net.Listen("tcp", ADDR)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreetingServer(grpcServer, &greetingService{})
	log.Printf("RPC Server is listening on %v", ADDR)
	log.Println(grpcServer.Serve(tcp))
}
