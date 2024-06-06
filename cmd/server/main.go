package main

import (
	"log"
	"net"

	"rail-booking/internal"
	"rail-booking/protogen/booking"

	"google.golang.org/grpc"
)

func main() {
	const addr = "0.0.0.0:50051"

	// create a TCP listener on the specified port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server instance
	server := grpc.NewServer()

	// create a rail booking service instance with a reference to the db and seat allocation service
	db := internal.NewDB()
	st := internal.NewAllocateSeat(db)
	railBookingService := internal.NewRailBookingService(db, st)

	// register the rail booking service with the grpc server
	booking.RegisterBookingServiceServer(server, &railBookingService)

	// start listening to requests
	log.Printf("server listening at %v", listener.Addr())
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
