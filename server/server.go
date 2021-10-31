package main

import (
	"context"
	"fielddaygrpc/fieldday"
	pb "fielddaygrpc/fieldday"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	port = ":50051"
)

type server struct {
	visitors []*pb.Visitor
	fieldday.UnimplementedFieldDayServer
	// https://github.com/grpc/grpc-go/issues/3794
}

// AddVisitor adds a visitor to the database
func (s *server) AddVisitor(
	ctx context.Context, in *pb.Visitor) (*fieldday.TotalNum, error) {
	s.visitors = append(s.visitors, in)
	out := new(fieldday.TotalNum)
	out.Number = int32(len(s.visitors))
	return out, nil

}

// GetTotal returns the total number of visitors
func (s *server) GetTotal(ctx context.Context, _ *emptypb.Empty) (*fieldday.TotalNum, error) {
	out := new(fieldday.TotalNum)
	out.Number = int32(len(s.visitors))
	return out, nil
}

// ListVisitors sends a stream of registered visitors
func (s *server) ListVisitors(_ *emptypb.Empty, stream pb.FieldDay_ListVisitorsServer) error {
	for _, visitor := range s.visitors {
		log.Print(visitor)
		err := stream.Send(visitor)
		if err != nil {
			return fmt.Errorf("error sending message to stream: %v", err)
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFieldDayServer(s, &server{})
	// Register reflection service on gRPC server
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
