package main

import (
	"context"
	"fmt"
	pb "github.com/ryantokmanmokmtm/gRPC-BiDirectional-Streaming/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

const (
	ADDR = "localhost:8080"
)

func Greeting(client pb.GreetingServiceClient) {
	log.Println("Greeting Service calling...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := client.GreetEveryOne(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	req := []*pb.GreetRequest{
		&pb.GreetRequest{Name: "Jackson"},
		&pb.GreetRequest{Name: "Tom"},
		&pb.GreetRequest{Name: "Amy"},
		&pb.GreetRequest{Name: "Alex"},
	}

	ch := make(chan struct{})

	go func() {
		for _, data := range req {
			log.Printf("Sending %v", data)
			if err := stream.Send(data); err != nil {
				log.Fatalln(err)
			}

			time.Sleep(time.Second * 1)
		}
		err := stream.CloseSend()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("Recevied: %v", data.Result)
		}

		close(ch) //close the channel
	}()

	<-ch //wait for closing
}

func GetMaxValue(client pb.GreetingServiceClient) {
	log.Println("GetMaxValue Service calling...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	stream, err := client.MaxCalculator(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	reqs := []*pb.MaxCalculatorRequest{
		&pb.MaxCalculatorRequest{Num: 1},
		&pb.MaxCalculatorRequest{Num: 5},
		&pb.MaxCalculatorRequest{Num: 3},
		&pb.MaxCalculatorRequest{Num: 6},
		&pb.MaxCalculatorRequest{Num: 20},
	}
	ch := make(chan struct{})
	go func() {
		for _, req := range reqs {
			log.Printf("Sending data %v\n", req.Num)
			if err := stream.Send(req); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second * 1)
		}

		err := stream.CloseSend()
		if err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("Recevied data %d\n", data.Result)
		}
		close(ch)
	}()
	<-ch
}

func main() {
	client, err := grpc.Dial(ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	stub := pb.NewGreetingServiceClient(client)
	//Greeting(stub)
	GetMaxValue(stub)
}
