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

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	//name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	//Modify Rail Booking by User Serice calls
	modifyRailBookingClientCall(ctx, c)

	//Cancel Rail Booking by User Serice calls
	cancelBookingClientCall(ctx, c)

}

func cancelBookingClientCall(ctx context.Context, c bk.BookingServiceClient) {
	cancelBookingsRequest := &bk.CancelBookingRequest{User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Cancel Rail Booking by User call check")
	cancelResp, err := c.CancelBooking(ctx, cancelBookingsRequest)
	if err != nil {
		log.Printf("Could not Cancel the booking: %v", err)
	}
	log.Printf("Successful Cancelation of Booking : %s", cancelResp)
	log.Println("")

	//Get All Rail Bookings Serice calls - After Cancelation
	getAllBookingsClientCall(ctx, c)

	//Try to Cancel wrong booking
	cancelBookingsRequest = &bk.CancelBookingRequest{User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragavenr@gmail.com"}}
	log.Println("Cancel Rail Booking by User call check")
	cancelResp, err = c.CancelBooking(ctx, cancelBookingsRequest)
	if err != nil {
		log.Printf("Could not Cancel the booking: %v", err)
	} else {
		log.Printf("Successful Cancelation of Booking : %s", cancelResp)
	}
	log.Println("")
}

func modifyRailBookingClientCall(ctx context.Context, c bk.BookingServiceClient) {
	modifyBookingsRequest := &bk.SeatModificationRequest{Section: "A", Seat: 40, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Modify Rail Booking by User call check")
	modifyResp, err := c.ModifySeatByUser(ctx, modifyBookingsRequest)
	if err != nil {
		log.Printf("Could not Modify the booking: %v", err)
	}
	log.Printf("Successful Modification of Booking : %s", modifyResp)
	log.Println("")

	//Get Rail Booking by User Serice calls after Modification

	bookingRequestByUser := &bk.GetBookingByUserRequest{User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Get Rail Booking by User call check after Modification")
	bookingResp, err1 := c.GetBookingByUser(ctx, bookingRequestByUser)
	if err != nil {
		log.Printf("Could not get the booking: %v", err1)
	}
	log.Printf("Successful Fetch of Booking for user after modification: %s", bookingResp)
	log.Println("")

	//Try to Modify wrong booking and error is thrown
	modifyBookingsRequest = &bk.SeatModificationRequest{Section: "A", Seat: 40, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragavenr@gmail.com"}}
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
		log.Printf("Could not get the booking: %v", err)
	}
	log.Printf("Successful Fetch of All Bookings : %s", r)
	log.Println("")
}

func GetBookingReqClientCall(ctx context.Context, c bk.BookingServiceClient) {
	bookingRequestByUser := &bk.GetBookingByUserRequest{User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Get Rail Booking by User call check")
	bookingResp, err := c.GetBookingByUser(ctx, bookingRequestByUser)
	if err != nil {
		log.Printf("Could not get the booking: %v", err)
	}
	log.Printf("Successful Fetch of Booking for user: %s", bookingResp)
	log.Println("")

	//Try to Get booking request of wrong user and error is thrown
	bookingRequestByUser = &bk.GetBookingByUserRequest{User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragavenr@gmail.com"}}
	log.Println("Get Rail Booking by User call check for wrong user and error is thrown")
	bookingResp, err = c.GetBookingByUser(ctx, bookingRequestByUser)
	if err != nil {
		log.Printf("Could not get the booking: %v", err)
	} else {
		log.Printf("Successful Fetch of Booking for user: %s", bookingResp)
	}
	log.Println("")
}

func railBookingClientCall(ctx context.Context, c bk.BookingServiceClient) {
	//Rail Booking Serice calls

	bookingRequest := &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Rail Booking call check")
	r, err := c.RailBooking(ctx, bookingRequest)
	if err != nil {
		log.Printf("Could not complete booking: %v", err)
	}
	log.Printf("Successful Booking: %s", r)
	log.Println("")

	//2nd Rail Booking Serice calls

	bookingRequest = &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragaven@gmail.com"}}
	log.Println("2nd Rail Booking call check")
	r, err = c.RailBooking(ctx, bookingRequest)
	if err != nil {
		log.Printf("Could not complete booking: %v", err)
	} else {
		log.Printf("2nd Successful Booking: %s", r)
	}
	log.Println("")

	//Duplicate Rail Booking Serice calls
	bookingRequest = &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragaven@gmail.com"}}
	log.Println("Duplicate Rail Booking call check")
	r, err = c.RailBooking(ctx, bookingRequest)
	if err != nil {
		log.Printf("Could not complete booking: %v", err)
	}
	log.Printf("Successful Duplicate Booking check: %s", r)
	log.Println("")

}
