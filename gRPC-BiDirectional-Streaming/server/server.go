package main

import (
	"fmt"
	pb "github.com/ryantokmanmokmtm/gRPC-BiDirectional-Streaming/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

const (
	ADDR = "localhost:8080"
)

type Service struct {
	pb.UnimplementedGreetingServiceServer
}

func (s *Service) GreetEveryOne(stream pb.GreetingService_GreetEveryOneServer) error {
	log.Println("GreetEveryOne Service Calling...")
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Server recevied : %v", data)
		err = stream.Send(&pb.GreetResponse{
			Result: fmt.Sprintf("Hello,%v !", data.Name),
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
	return nil
}

func (s *Service) MaxCalculator(stream pb.GreetingService_MaxCalculatorServer) error {
	log.Println("MaxCalculator Service calling...")
	var max int32 = 0
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Server Received : %v\n", data)
		if data.Num > max {
			max = data.Num
			fmt.Printf("Server Sending data %v\n", max)
			err = stream.Send(&pb.MaxCalculatorResponse{
				Result: max,
			})
			if err != nil {
				log.Fatalln(err)
			}
		}

	}

	return nil
}

func main() {
	tcp, err := net.Listen("tcp", ADDR)
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreetingServiceServer(grpcServer, &Service{})
	log.Printf("RPC Server is listening on %v", ADDR)
	grpcServer.Serve(tcp)
}
