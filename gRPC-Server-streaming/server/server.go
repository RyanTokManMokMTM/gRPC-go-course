package main

import (
	"fmt"
	pb "github.com/ryantokmanmokmtm/gRPC-Server-streaming/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	ADDR = "localhost:50001"
)

type Service struct {
	pb.UnimplementedGreetingServer
}

func (s *Service) GreetingMany(in *pb.GreetingRequest, stream pb.Greeting_GreetingManyServer) error {
	log.Printf("Server streaming greeting with %v", in)

	for i := 0; i < 10; i++ {
		stream.Send(&pb.GreetingResponse{
			Result: fmt.Sprintf("Message '%v' recevied ,time %v", in.Message, i+1),
		})
	}
	return nil
}

func (s *Service) CalculatePrimeNum(in *pb.Number, stream pb.Greeting_CalculatePrimeNumServer) error {
	log.Printf("Calculating a prime number of %d", in.Num)
	k := int32(2)
	N := in.Num

	for {
		if N < 2 {
			break
		}
		if N%k == 0 {
			stream.Send(&pb.PrimeResponse{
				Prime: k,
			})
			N /= k
		} else {
			k += 1
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
	pb.RegisterGreetingServer(grpcServer, &Service{})
	log.Printf("RPC Server is listening on %v", ADDR)
	log.Fatalln(grpcServer.Serve(tcp))
}
