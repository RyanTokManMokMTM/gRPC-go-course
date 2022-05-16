package main

import (
	"context"
	pb "github.com/ryantokmanmokmtm/gRPC-Unary/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

const (
	ADDR = "localhost:50001"
)

func sayHello(name string, client pb.GreetingClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := client.SayHello(ctx, &pb.GreetingRequest{
		YourName: name,
	})

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("result %v", res.Result)
}

func getSum(valA, valB int32, client pb.GreetingClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := client.SumUp(ctx, &pb.SumRequest{
		SumA: valA,
		SumB: valB,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("sum up result:%v", res.GetResult())
}

func main() {
	//create a grpc client and disables transport security(NO TLS).
	client, err := grpc.Dial(ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	//create a client stub
	clientStub := pb.NewGreetingClient(client)
	sayHello("Jackson", clientStub) //calling sayHello service
	getSum(int32(3), int32(10), clientStub)
}
