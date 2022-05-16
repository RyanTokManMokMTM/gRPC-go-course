package main

import (
	"context"
	pb "github.com/ryantokmanmokmtm/gRPC-Server-streaming/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

const (
	ADDR = "localhost:50001"
)

func greetingMany(client pb.GreetingClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stream, err := client.GreetingMany(ctx, &pb.GreetingRequest{Message: message})
	if err != nil {
		log.Fatalln(err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("response message : %v", msg)
	}
}

func getPrime(client pb.GreetingClient, num int32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := client.CalculatePrimeNum(ctx, &pb.Number{
		Num: num,
	})
	if err != nil {
		log.Fatalln(err)
	}
	for {
		primeNum, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("%v", primeNum)
	}
}
func main() {
	client, err := grpc.Dial(ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	stub := pb.NewGreetingClient(client)
	greetingMany(stub, "Server again")
	getPrime(stub, 120)
}
