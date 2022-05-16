package main

import (
	"context"
	pb "github.com/ryantokmanmokmtm/gRPC-client-streaming/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	ADDR = "localhost:50001"
)

func KeepGreeting(client pb.GreetingClient) {
	log.Printf("Keep Greeting....")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := client.KeepGreeting(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	datas := []string{
		"Jackson",
		"Tom",
		"Amy",
		"Karry",
	}

	//Keep sending message to server
	for _, data := range datas {
		log.Printf("Sending %v", data)
		if err := stream.Send(&pb.GreetingRequest{
			Name: data,
		}); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Second * 1)
	}

	//waiting for response and close
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	print(res.Result)
}

func AvgCalculator(client pb.GreetingClient) {
	log.Printf("AvgCalculator Service calling...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := client.AvgCalculator(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	nums := []int{
		1, 2, 3, 4,
	}

	for _, num := range nums {
		if err := stream.Send(&pb.CalculatorRequest{
			Number: int32(num),
		}); err != nil {
			log.Fatalln(err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Avg number is %v", res.Result)

}

func main() {
	client, err := grpc.Dial(ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	stub := pb.NewGreetingClient(client)
	//KeepGreeting(stub)
	AvgCalculator(stub)
}
