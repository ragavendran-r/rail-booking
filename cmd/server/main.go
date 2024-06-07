package main

import (
	"log"
	"net"
	"rail-booking/internal"
	bk "rail-booking/protogen/booking"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

	healthServer := health.NewServer()

	go func() {
		for {
			status := healthpb.HealthCheckResponse_SERVING
			// Check if user Service is valid
			if time.Now().Second()%2 == 0 {
				status = healthpb.HealthCheckResponse_NOT_SERVING
			}

			healthServer.SetServingStatus(bk.BookingService_ServiceDesc.ServiceName, status)
			healthServer.SetServingStatus("", status)

			time.Sleep(1 * time.Second)
		}
	}()

	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(bk.BookingService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)

	// register the rail booking service with the grpc server
	bk.RegisterBookingServiceServer(server, &railBookingService)
	healthpb.RegisterHealthServer(server, healthServer)
	// start listening to requests
	log.Printf("server listening at %v", listener.Addr())
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
