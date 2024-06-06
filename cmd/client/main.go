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

	//Rail Booking Serice calls

	bookingRequest := &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Rail Booking call check")
	r, err := c.RailBooking(ctx, bookingRequest)
	if err != nil {
		log.Fatalf("Could not complete booking: %v", err)
	}
	log.Printf("Successful Booking: %s", r)
	log.Println("")
	//2nd Rail Booking Serice calls
	bookingRequest = &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragaven@gmail.com"}}
	log.Println("2nd Rail Booking call check")
	r, err = c.RailBooking(ctx, bookingRequest)
	if err != nil {
		log.Fatalf("Could not complete booking: %v", err)
	}
	log.Printf("2nd Successful Booking: %s", r)
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

	//Get Rail Booking by User Serice calls

	bookingRequestByUser := &bk.GetBookingByUserRequest{User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Get Rail Booking by User call check")
	r, err = c.GetBookingByUser(ctx, bookingRequestByUser)
	if err != nil {
		log.Fatalf("Could not get the booking: %v", err)
	}
	log.Printf("Successful Fetch of Booking for user: %s", r)
	log.Println("")

	//Get All Rail Bookings Serice calls

	getAllBookingsRequest := &bk.Empty{}
	log.Println("Get All Rail Bookings call check")
	r1, err1 := c.GetAllBookings(ctx, getAllBookingsRequest)
	if err1 != nil {
		log.Fatalf("Could not get the booking: %v", err1)
	}
	log.Printf("Successful Fetch of All Bookings : %s", r1)
	log.Println("")

	//Modify Rail Booking by User Serice calls

	modifyBookingsRequest := &bk.SeatModificationRequest{Section: "A", Seat: 40, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Modify Rail Booking by User call check")
	r2, err2 := c.ModifySeatByUser(ctx, modifyBookingsRequest)
	if err2 != nil {
		log.Fatalf("Could not Modify the booking: %v", err2)
	}
	log.Printf("Successful Modification of Booking : %s", r2)
	log.Println("")

	//Get Rail Booking by User Serice calls after Modification

	bookingRequestByUser = &bk.GetBookingByUserRequest{User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Get Rail Booking by User call check after Modification")
	r, err = c.GetBookingByUser(ctx, bookingRequestByUser)
	if err != nil {
		log.Fatalf("Could not get the booking: %v", err)
	}
	log.Printf("Successful Fetch of Booking for user after modification: %s", r)
	log.Println("")

	//Cancel Rail Booking by User Serice calls

	cancelBookingsRequest := &bk.CancelBookingRequest{User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Cancel Rail Booking by User call check")
	r3, err3 := c.CancelBooking(ctx, cancelBookingsRequest)
	if err3 != nil {
		log.Fatalf("Could not Cancel the booking: %v", err3)
	}
	log.Printf("Successful Cancelation of Booking : %s", r3)
	log.Println("")

	//Get All Rail Bookings Serice calls - After Cancelation

	getAllBookingsRequest = &bk.Empty{}
	log.Println("Get All Rail Bookings call check - After Cancelation")
	r1, err1 = c.GetAllBookings(ctx, getAllBookingsRequest)
	if err1 != nil {
		log.Fatalf("Could not get the booking: %v", err1)
	}
	log.Printf("Successful Fetch of All Bookings after cancelation : %s", r1)
	log.Println("")

}
