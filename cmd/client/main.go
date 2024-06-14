package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	bk "rail-booking/protogen/booking"
	ur "rail-booking/protogen/user"
	"runtime/pprof"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	//name = flag.String("name", defaultName, "Name to greet")
	TEST_USER1 = &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}
	TEST_USER2 = &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragavenr@gmail.com"}
	TEST_USER3 = &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragaven@gmail.com"}
)

const BKNGS_NOT_FOUND string = "Bookings not found"
const BKNGS_NOT_SUCCESS string = "Bookings not successful"

func main() {
	flag.Parse()
	// Start profiling
	f, err := os.Create("railbooking.prof")
	if err != nil {

		fmt.Println(err)
		return

	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	//Enable or disable TLS
	tls := false // change that to true if needed
	opts := []grpc.DialOption{}
	if tls {
		certFile := "ssl/ca.crt"
		creds, err := credentials.NewClientTLSFromFile(certFile, "")

		if err != nil {
			log.Fatalf("Error loading CA trust certificate: %v\n", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
		// Set up a connection to the server.

	}
	//Comment if TLS is true
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//Uncomment if TLS is true
	//conn, err: = grpc.NewClient(*addr, opts...)
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := bk.NewBookingServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//Rail Booking Serice calls
	railBookingClientCall(ctx, c)

	//Get Rail Booking by User Serice calls
	GetBookingReqClientCall(ctx, c)

	//Get All Rail Bookings Serice calls
	getAllBookingsClientCall(ctx, c)

	//Get Section Rail Bookings Serice calls
	getSectionBookingsClientCall(ctx, c)

	//Modify Rail Booking by User Serice calls
	modifyRailBookingClientCall(ctx, c)

	//Cancel Rail Booking by User Serice calls
	cancelBookingClientCall(ctx, c)

}

func cancelBookingClientCall(ctx context.Context, c bk.BookingServiceClient) {
	cancelBookingsRequest := &bk.CancelBookingRequest{User: TEST_USER1}
	log.Println("Cancel Rail Booking by User call check")
	cancelResp, err := c.CancelBooking(ctx, cancelBookingsRequest)
	if err != nil {
		//gRPC Error handling
		e, ok := status.FromError(err)
		if ok {
			log.Printf("Error message from Server : %s", e.Message())
			log.Printf("Error Code from Server : %s", e.Code())
		} else {
			log.Fatalf("A non gRPC error: %v\n", e)
		}

		log.Printf("Could not Cancel the booking: %v", err)
	}
	log.Printf("Successful Cancelation of Booking : %s", cancelResp)
	log.Println("")

	//Get All Rail Bookings Serice calls - After Cancelation
	getAllBookingsClientCall(ctx, c)

	//Try to Cancel wrong booking
	cancelBookingsRequest = &bk.CancelBookingRequest{User: TEST_USER2}
	log.Println("Cancel Rail Booking by User call check")
	cancelResp, err = c.CancelBooking(ctx, cancelBookingsRequest)
	if err != nil {
		//gRPC Error handling
		e, ok := status.FromError(err)
		if ok {
			log.Printf("Error message from Server : %s", e.Message())
			log.Printf("Error Code from Server : %s", e.Code())
		} else {
			log.Fatalf("A non gRPC error: %v\n", e)
		}
		log.Printf("Could not Cancel the booking: %v", err)
	} else {
		log.Printf("Successful Cancelation of Booking : %s", cancelResp)
	}
	log.Println("")
}

func modifyRailBookingClientCall(ctx context.Context, c bk.BookingServiceClient) {
	modifyBookingsRequest := &bk.SeatModificationRequest{Section: "A", Seat: 40, User: TEST_USER1}
	log.Println("Modify Rail Booking by User call check")
	modifyResp, err := c.ModifySeatByUser(ctx, modifyBookingsRequest)
	if err != nil {
		log.Printf("Could not Modify the booking: %v", err)
	}
	log.Printf("Successful Modification of Booking : %s", modifyResp)
	log.Println("")

	//Get Rail Booking by User Serice calls after Modification

	bookingRequestByUser := &bk.GetBookingByUserRequest{User: TEST_USER1}
	log.Println("Get Rail Booking by User call check after Modification")
	bookingResp, err1 := c.GetBookingByUser(ctx, bookingRequestByUser)
	if err != nil {
		log.Printf("BKNGS_NOT_FOUND: %v", err1)
	}
	log.Printf("Successful Fetch of Booking for user after modification: %s", bookingResp)
	log.Println("")

	//Try to Modify wrong booking and error is thrown
	modifyBookingsRequest = &bk.SeatModificationRequest{Section: "A", Seat: 40, User: TEST_USER2}
	log.Println("Modify Rail Booking by wrong User call check and error thrown")
	modifyResp, err = c.ModifySeatByUser(ctx, modifyBookingsRequest)
	if err != nil {
		log.Printf("Could not Modify the booking: %v", err)
	} else {
		log.Printf("Successful Modification of Booking : %s", modifyResp)
	}
	log.Println("")
}

func getAllBookingsClientCall(ctx context.Context, c bk.BookingServiceClient) {
	getAllBookingsRequest := &bk.Empty{}
	log.Println("Get All Rail Bookings call check")
	r, err := c.GetAllBookings(ctx, getAllBookingsRequest)
	if err != nil {
		log.Printf("BKNGS_NOT_FOUND: %v", err)
	}
	log.Printf("Successful Fetch of All Bookings : %s", r)
	log.Println("")
}

func getSectionBookingsClientCall(ctx context.Context, c bk.BookingServiceClient) {
	getSectionBookingsRequest := &bk.GetSectionBookingsRequest{Section: "A"}
	log.Println("Get Section Rail Bookings call check")
	r, err := c.GetSectionBookings(ctx, getSectionBookingsRequest)
	if err != nil {
		log.Printf("BKNGS_NOT_FOUND: %v", err)
	}
	log.Printf("Successful Fetch of All Bookings for a section : %s", r)
	log.Println("")
}

func GetBookingReqClientCall(ctx context.Context, c bk.BookingServiceClient) {
	bookingRequestByUser := &bk.GetBookingByUserRequest{User: TEST_USER1}
	log.Println("Get Rail Booking by User call check")
	bookingResp, err := c.GetBookingByUser(ctx, bookingRequestByUser)
	if err != nil {
		log.Printf("BKNGS_NOT_FOUND: %v", err)
	}
	log.Printf("Successful Fetch of Booking for user: %s", bookingResp)
	log.Println("")

	//Try to Get booking request of wrong user and error is thrown
	bookingRequestByUser = &bk.GetBookingByUserRequest{User: TEST_USER2}
	log.Println("Get Rail Booking by User call check for wrong user and error is thrown")
	bookingResp, err = c.GetBookingByUser(ctx, bookingRequestByUser)
	if err != nil {
		log.Printf("BKNGS_NOT_FOUND: %v", err)
	} else {
		log.Printf("Successful Fetch of Booking for user: %s", bookingResp)
	}
	log.Println("")
}

func railBookingClientCall(ctx context.Context, c bk.BookingServiceClient) {
	//Rail Booking Serice calls

	bookingRequest := &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: TEST_USER1}
	log.Println("Rail Booking call check")
	r, err := c.RailBooking(ctx, bookingRequest)
	if err != nil {
		log.Printf("BKNGS_NOT_SUCCESS: %v", err)
	}
	log.Printf("Successful Booking: %s", r)
	log.Println("")

	//2nd Rail Booking Serice calls

	bookingRequest = &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: TEST_USER3}
	log.Println("2nd Rail Booking call check")
	r, err = c.RailBooking(ctx, bookingRequest)
	if err != nil {
		log.Printf("BKNGS_NOT_SUCCESS: %v", err)
	} else {
		log.Printf("2nd Successful Booking: %s", r)
	}
	log.Println("")

	//Duplicate Rail Booking Serice calls
	bookingRequest = &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: TEST_USER3}
	log.Println("Duplicate Rail Booking call check")
	r, err = c.RailBooking(ctx, bookingRequest)
	if err != nil {
		log.Printf("BKNGS_NOT_SUCCESS: %v", err)
	}
	log.Printf("Successful Duplicate Booking check: %s", r)
	log.Println("")

}
