package main

/*
TODO for RPC Server
1. Implementing our services interface that defined in protobuf
2. Running a gRPC Server to listen for request from client
*/

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	pb "github.com/ryantokmanmokmtm/basic-grpc/proto"
	"google.golang.org/grpc"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"sync"
	"time"
)

var (
	//tls        = flag.Bool("flag", false, "")
	//certFile   = flag.String("cert_file", "", "")
	//keyFile    = flag.String("key_file", "", "")
	//jsonDBFile = flag.String("json_db_file", "", "")
	port = flag.Int("port", 50051, "")
)

type routeGuideService struct {
	pb.UnimplementedRouteGuideServer //must embedded
	savedFeatures                    []*pb.Feature

	mutx       sync.Mutex
	routeNotes map[string][]*pb.RouteNote
}

func (s *routeGuideService) GetFeature(ctx context.Context, in *pb.Point) (*pb.Feature, error) {
	//getting feature from our fake datas
	for _, feature := range s.savedFeatures {
		if proto.Equal(in, feature) {
			return feature, nil
		}
	}
	//feature was not found

	return &pb.Feature{Location: in}, nil
}

func (s *routeGuideService) ListFeatures(in *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.savedFeatures {
		//need a function to help for calculating a range of feature
		if inRange(feature.Location, in) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	return nil
}

// RecordRoute
// getting a stream of points and response
func (s *routeGuideService) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	var pointCounter, featureCount, distance int32
	var lastPoint *pb.Point
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			//ended streaming
			endTime := time.Now()

			//after client have sent all the point
			//server return a summary
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCounter,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return nil
		}
		//server received a point
		pointCounter++

		//before adding feature count, it needs to check the feature exits in our lists
		for _, feature := range s.savedFeatures {
			//find a matched location from feature
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}

		if lastPoint != nil {
			//calculating the distance
			distance += calculateDistance(lastPoint, point)
		}

		lastPoint = point
	}

}

//RouteChat receives a stream of location pairs and response with a stream of all previous messages
//at each of those locations
func (s *routeGuideService) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
	for {
		//receives
		in, err := stream.Recv()
		if err == io.EOF {
			return stream.Send(&pb.RouteNote{})
		}
		if err != nil {
			log.Printf("err:%v", err)
			return err
		}
		key := serialize(in.Location)
		log.Printf("received:%v", in)
		//adding the new notes to our datasets
		s.mutx.Lock()
		s.routeNotes[key] = append(s.routeNotes[key], in)

		res := make([]*pb.RouteNote, len(s.routeNotes))
		copy(res, s.routeNotes[key]) //copy all data from datasets
		s.mutx.Unlock()

		for _, note := range res {
			//log.Printf("%v", note)
			if err := stream.Send(note); err != nil {
				log.Printf("%v", err)
				return err
			}
		}
	}
}

//LoadFeature from JSON file
func (s *routeGuideService) LoadFeature(filePath string) {
	//filePath
	var data []byte
	if filePath != "" {
		var err error
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalln(err)
		}

		//read data from default???
	}

	if err := json.Unmarshal(data, &s.savedFeatures); err != nil {
		log.Fatalln(err)
	}
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

func calculateDistance(p1, p2 *pb.Point) int32 {
	const Factor = 1e7
	const R = float64(6371000) //earth r in metres

	lat1 := toRadians(float64(p1.Latitude) / Factor)
	lat2 := toRadians(float64(p2.Latitude) / Factor)
	lng1 := toRadians(float64(p1.Longitude) / Factor)
	lng2 := toRadians(float64(p2.Longitude) / Factor)

	distanceLat := lat2 - lat1
	distanceLng := lng2 - lng1

	a := math.Sin(distanceLat/2)*math.Sin(distanceLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(distanceLng/2)*math.Sin(distanceLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)

}

func serialize(p *pb.Point) string {
	return fmt.Sprintf("%d %d", p.Latitude, p.Longitude)
}

//inRange whether the point in the rect
func inRange(p *pb.Point, rect *pb.Rectangle) bool {
	l := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	r := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	t := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	b := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(p.Longitude) >= l && float64(p.Longitude) <= r && float64(p.Latitude) >= t && float64(p.Latitude) <= b {
		return true
	}
	return false
}

func newServer() *routeGuideService {
	s := &routeGuideService{
		routeNotes: map[string][]*pb.RouteNote{},
	}
	s.LoadFeature("../route_guide_db.json")
	return s
}

func main() {
	flag.Parse()

	//open tcp connection
	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		log.Fatalln(err)
	}

	server := grpc.NewServer()
	//register gRPC server
	pb.RegisterRouteGuideServer(server, newServer())
	log.Println("RPC Server is listening...")
	log.Println(server.Serve(listen))

}
