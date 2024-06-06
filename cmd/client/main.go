package main

import (
	"context"
	"flag"
	"log"
	bk "rail-booking/protogen/booking"
	ur "rail-booking/protogen/user"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := bk.NewBookingServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//Rail Booking Serice call

	bookingRequest := &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Rail Booking call check")
	r, err := c.RailBooking(ctx, bookingRequest)
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}
	log.Printf("Successful Booking: %s", r)
}
