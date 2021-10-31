package main

import (
	"context"
	pb "fielddaygrpc/fieldday"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFieldDayClient(conn)

	// Add the visitor and print the server's response
	first := "Pavel"
	last := "Anni"
	callsign := "AC4PA"
	youth := false
	nfarl := true
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.AddVisitor(ctx, &pb.Visitor{First: first, Last: last, Callsign: callsign, Youth: youth, NFARL: nfarl})
	if err != nil {
		log.Fatalf("could not add visitor: %v", err)
	}
	log.Printf("Visitor %s %s was added successfully. Total number: %d", first, last, int(r.Number))
	r, err = c.GetTotal(ctx, &emptypb.Empty{})
	if err != nil {
		log.Fatalf("could not get the number of visitors: %v", err)
	}
	log.Printf("Total number of Field Day visitors: %d", int(r.Number))

	listStream, _ := c.ListVisitors(ctx, &emptypb.Empty{})
	for {
		visitor, err := listStream.Recv()
		if err == io.EOF {
			log.Print("EOF")
			break
		}
		if err == nil {
			log.Println(visitor.First, visitor.Last, visitor.Callsign)
		}
	}
}
