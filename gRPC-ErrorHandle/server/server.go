package main

import (
	"context"
	"fmt"
	pb "github.com/ryantokmanmokmtm/gRPC-ErrorHandle/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"math"
	"net"
	"time"
)

const (
	ADDR = "localhost:50001"
	TLS  = false
)

type Service struct {
	pb.UnimplementedMathServiceServer
}

func (s *Service) Sqrt(ctx context.Context, in *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	log.Println("Sqrt Service calling...")

	//sending rpc error
	if in.SqrtNum < 1 {
		return nil, status.Errorf(codes.InvalidArgument,
			fmt.Sprintf("Recevied a negative number %v", in.SqrtNum))
	}

	return &pb.SqrtResponse{
		Result: float32(math.Sqrt(float64(in.SqrtNum))),
	}, nil
}
func (s *Service) DeadlineHandle(ctx context.Context, in *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	log.Println("Sqrt Service calling...")

	//simulate network delay
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			return nil,
				status.Errorf(codes.Canceled, fmt.Sprintf("Client canceled the request of {%v}", in.Name))
		}

		//network delay by 1 second
		time.Sleep(time.Second * 1)
	}

	return &pb.GreetingResponse{Result: fmt.Sprintf("Hello,%v", in.Name)}, nil
}

func main() {
	tcp, err := net.Listen("tcp", ADDR)
	if err != nil {
		log.Fatalln(err)
	}
	var opts []grpc.ServerOption
	//using TLS
	if TLS {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"

		creds, err := credentials.NewClientTLSFromFile(certFile, keyFile)
		if err != nil {
			log.Fatalln(err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMathServiceServer(grpcServer, &Service{})
	//do a reflection and allow evans to use
	reflection.Register(grpcServer)
	log.Printf("RPC Server is listening on %v", ADDR)
	log.Fatalln(grpcServer.Serve(tcp))
}
