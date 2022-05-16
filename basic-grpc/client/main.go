package main

import (
	"context"
	"flag"
	pb "github.com/ryantokmanmokmtm/basic-grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"math/rand"
	"time"
)

var (
	serverAddr = flag.String("server_addr", "localhost:50051", "")
)

func getFeature(client pb.RouteGuideClient, point *pb.Point) {
	//LOG INFO
	log.Printf("Getting Feature for point (%d,%d)", point.Latitude, point.Longitude)

	//Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() //cancel the request

	//Calling the simple RPC service -Get feature
	feature, err := client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("%v.GetFeature().Error : %v", client, err)
	}
	log.Println(feature)
}

func listRecord(client pb.RouteGuideClient, rect *pb.Rectangle) {
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("%v.ListFeatures.Error: %v", client, err)
	}

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v.ListFeatures.Error: %v", client, err)
		}
		log.Printf("Feature: name: %q, point:(%d,%d)", data.Name, data.GetLocation().GetLatitude(), data.GetLocation().GetLongitude())
	}
}

func recordRoute(client pb.RouteGuideClient) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(rand.Int31n(100) + 2)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var points []*pb.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}

	stream, err := client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("%v,RecordRoute(_).Error: %v", client, err)
	}
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Fatalf("%v,Send(%v).Error: %v", stream, point, err)
		}
	}
	//close and stop receiving
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v,CloseAndRecv().Error: %v", stream, err)
	}
	log.Printf("Route summary: %v", reply)
}

func routeChat(client pb.RouteGuideClient) {
	exampleNote := []*pb.RouteNote{
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "First Message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Second Message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Third Message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 4}, Message: "Fourth Message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 5}, Message: "Fifth Message"},
		{Location: &pb.Point{Latitude: 0, Longitude: 6}, Message: "Sixth Message"},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	channel := make(chan struct{})
	stream, err := client.RouteChat(ctx)
	if err != nil {
		log.Fatalf("%v.RouteChat().Error: %v", client, err)
	}
	//coz it's bi-streaming
	//so, it can read a message from the server and write a message to the server

	//create a goroutine for receiving data from server
	go func() {
		for {
			data, err := stream.Recv()
			if err == io.EOF {
				//close the channel
				close(channel)
				return
			}

			if err != nil {
				log.Fatalf("Failed to recevie a note: err %v", err)
			}
			log.Printf("Got message %s at point(%d,%d)", data.Message, data.GetLocation().GetLatitude(), data.GetLocation().GetLongitude())
		}
	}()

	//sending message to server
	for _, note := range exampleNote {
		log.Println(note)
		if err := stream.Send(note); err != nil {
			log.Fatalf("Failed to send a note: error %v", err)
		}
		time.Sleep(time.Second * 1)
	}
	stream.CloseSend()
	//blocking until closed
	<-channel
}

//randomPoint - generating random point by a seed
func randomPoint(r *rand.Rand) *pb.Point {
	latitude := (r.Int31n(180) - 90) * 1e7
	longitude := (r.Int31n(360) - 180) * 1e7
	return &pb.Point{
		Latitude:  latitude,
		Longitude: longitude,
	}
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	//create a channel which uses to communicate with server
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close() //close the connection

	//using stub to calling the service method
	client := pb.NewRouteGuideClient(conn)

	//getFeature(client, &pb.Point{
	//	Latitude:  409146138,
	//	Longitude: -746188906,
	//})

	//location not exists
	//getFeature(client, &pb.Point{
	//	Latitude:  0,
	//	Longitude: 0,
	//})

	//rect := &pb.Rectangle{
	//	Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
	//	Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
	//}
	//listRecord(client, rect)
	//
	//recordRoute(client)
	//
	routeChat(client)
}
