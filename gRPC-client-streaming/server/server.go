package main

import (
	"fmt"
	pb "github.com/ryantokmanmokmtm/gRPC-client-streaming/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

const (
	ADDR = "localhost:50001"
)

type Service struct {
	pb.UnimplementedGreetingServer
}

func (s *Service) KeepGreeting(stream pb.Greeting_KeepGreetingServer) error {
	log.Printf("Long Greeting Service calling")
	var res string
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("Finished and wait to close")
			return stream.SendAndClose(&pb.GreetingResponse{
				Result: res,
			})
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Receving message : %v", data)
		res += fmt.Sprintf("Recevied message: %v \n", data)
	}

	return nil
}

func (s *Service) AvgCalculator(stream pb.Greeting_AvgCalculatorServer) error {
	log.Printf("AvgCalculator Service Calling...")
	var res float32
	var counter int = 0
	for {
		num, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Finished and wait to close...")
			return stream.SendAndClose(&pb.CalculatorResponse{
				Result: res / float32(counter),
			})
		}
		if err != nil {
			log.Fatalln(err)
		}
		res += float32(num.Number)
		counter++
	}
	return nil
}
func main() {
	tcp, err := net.Listen("tcp", ADDR)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreetingServer(grpcServer, &Service{})
	log.Printf("RPC Server is listen on %v", ADDR)
	log.Fatalln(grpcServer.Serve(tcp))
}
