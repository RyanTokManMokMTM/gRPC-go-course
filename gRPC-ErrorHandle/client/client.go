package main

import (
	"context"
	pb "github.com/ryantokmanmokmtm/gRPC-ErrorHandle/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

const (
	ADDR = "localhost:50001"
	TLS  = false
)

func doSqrt(client pb.MathServiceClient, n int32) {
	log.Printf("Doing Sqrt of %v", n)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	res, err := client.Sqrt(ctx, &pb.SqrtRequest{SqrtNum: n})
	if err != nil {
		//getting RPC Error
		e, ok := status.FromError(err)
		if ok {
			//it is an RPC error if ok is true
			log.Printf("Error Message from server: %v \n", e.Message())
			log.Printf("Error code from server: %v \n", e.Code())

			//what the error is?
			//is it what we expected?
			if e.Code() == codes.InvalidArgument {
				log.Printf("We proberly sent a negative number!")
				return
			}

		} else {
			//another error
			log.Fatalln(err)
		}
	}
	log.Printf("Sqrt Result : %v", res.Result)
}

func doDeadline(client pb.MathServiceClient, timeout time.Duration) {
	log.Printf("Doing deadline service...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := client.DeadlineHandle(ctx, &pb.GreetingRequest{Name: "jackson"})
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			if e.Code() == codes.DeadlineExceeded {
				log.Println("Deadline Exceeded...")
				return
			}
		} else {
			log.Fatalln(err)
		}
	}
	log.Printf("Recevied from server : %v", res.Result)
}

func main() {
	var opts []grpc.DialOption
	if TLS {
		certFile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")
		if err != nil {
			log.Fatalln(err)
		}

		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	//client, err := grpc.Dial(ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	client, err := grpc.Dial(ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//client, err := grpc.Dial(ADDR, opts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	stub := pb.NewMathServiceClient(client)
	//doSqrt(stub, -5)
	doSqrt(stub, 5)
	//doDeadline(stub, time.Second*5)
	//doDeadline(stub, time.Second*1)
}
