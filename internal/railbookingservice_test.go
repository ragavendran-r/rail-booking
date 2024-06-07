package internal

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	bk "rail-booking/protogen/booking"
	ur "rail-booking/protogen/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func newServer(t *testing.T, register func(srv *grpc.Server)) *grpc.ClientConn {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	register(srv)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}

	return conn
}

func TestRailBookingService_RailBooking(t *testing.T) {

	// create a rail booking service instance with a reference to the db and seat allocation service
	db := NewDB()
	st := NewAllocateSeat(db)
	railBookingService := NewRailBookingService(db, st)
	conn := newServer(t, func(srv *grpc.Server) {
		bk.RegisterBookingServiceServer(srv, &railBookingService)
	})

	client := bk.NewBookingServiceClient(conn)
	res, err := client.RailBooking(context.Background(), &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}})
	if err != nil {
		t.Fatalf("client.RailBooking %v", err)
	}

	if res.GetSection() == "" && res.GetSeat() == 0 {
		t.Fatalf("Unexpected values %v", res)
	}
}

func TestRailBookingService_GetBookingByUser(t *testing.T) {

	// create a rail booking service instance with a reference to the db and seat allocation service
	db := NewDB()
	st := NewAllocateSeat(db)
	railBookingService := NewRailBookingService(db, st)
	conn := newServer(t, func(srv *grpc.Server) {
		bk.RegisterBookingServiceServer(srv, &railBookingService)
	})

	client := bk.NewBookingServiceClient(conn)
	res, err := client.RailBooking(context.Background(), &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}})
	if err != nil {
		t.Fatalf("client.RailBooking %v %v", err, res)
	}
	bookingRequestByUser := &bk.GetBookingByUserRequest{User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Get Rail Booking by User call check for wrong user and error is thrown")
	bookingResp, err1 := client.GetBookingByUser(context.Background(), bookingRequestByUser)
	if err1 != nil {
		t.Fatalf("client.GetBookingByUser %v", err1)
	}
	if bookingResp.GetSection() == "" && bookingResp.GetSeat() == 0 {
		t.Fatalf("Unexpected values %v", bookingResp)
	}
}

func TestRailBookingService_GetAllBookings(t *testing.T) {
	log.Println("Testing GetAllBookings Service")
	// create a rail booking service instance with a reference to the db and seat allocation service
	db := NewDB()
	st := NewAllocateSeat(db)
	railBookingService := NewRailBookingService(db, st)
	conn := newServer(t, func(srv *grpc.Server) {
		bk.RegisterBookingServiceServer(srv, &railBookingService)
	})

	client := bk.NewBookingServiceClient(conn)
	res, err := client.RailBooking(context.Background(), &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}})
	if err != nil {
		t.Fatalf("client.RailBooking %v %v", err, res)
	}
	getAllBookingsRequest := &bk.Empty{}
	log.Println("Get All Rail Bookings")
	bookingsResp, err1 := client.GetAllBookings(context.Background(), getAllBookingsRequest)
	if err1 != nil {
		t.Fatalf("client.GetBookingByUser %v", err1)
	}
	if bookingsResp.GetBookings()[0].GetSection() == "" && bookingsResp.GetBookings()[0].GetSeat() == 0 {
		t.Fatalf("Unexpected values %v", bookingsResp)
	}
}

func TestRailBookingService_GetSectionBookings(t *testing.T) {
	log.Println("Testing GetAllBookings Service")
	// create a rail booking service instance with a reference to the db and seat allocation service
	db := NewDB()
	st := NewAllocateSeat(db)
	railBookingService := NewRailBookingService(db, st)
	conn := newServer(t, func(srv *grpc.Server) {
		bk.RegisterBookingServiceServer(srv, &railBookingService)
	})

	client := bk.NewBookingServiceClient(conn)
	res, err := client.RailBooking(context.Background(), &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}})
	if err != nil {
		t.Fatalf("client.RailBooking %v %v", err, res)
	}
	getSectionBookingsRequest := &bk.GetSectionBookingsRequest{Section: "A"}
	log.Println("Get All Rail Bookings")
	bookingsResp, err1 := client.GetSectionBookings(context.Background(), getSectionBookingsRequest)
	if err1 != nil {
		t.Fatalf("client.GetBookingByUser %v", err1)
	}
	if len(bookingsResp.GetBookings()) > 0 && bookingsResp.GetBookings()[0].GetSection() == "" && bookingsResp.GetBookings()[0].GetSeat() == 0 {
		t.Fatalf("Unexpected values %v", bookingsResp)
	}
}

func TestRailBookingService_ModifySeatByUser(t *testing.T) {

	// create a rail booking service instance with a reference to the db and seat allocation service
	db := NewDB()
	st := NewAllocateSeat(db)
	railBookingService := NewRailBookingService(db, st)
	conn := newServer(t, func(srv *grpc.Server) {
		bk.RegisterBookingServiceServer(srv, &railBookingService)
	})

	client := bk.NewBookingServiceClient(conn)
	res, err := client.RailBooking(context.Background(), &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}})
	if err != nil {
		t.Fatalf("client.RailBooking %v %v", err, res)
	}
	modifyBookingsRequest := &bk.SeatModificationRequest{Section: "A", Seat: 40, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Get Rail Booking by User call check for wrong user and error is thrown")
	modResp, err1 := client.ModifySeatByUser(context.Background(), modifyBookingsRequest)
	if err1 != nil {
		t.Fatalf("client.ModifySeatByUser %v", err1)
	}
	if modResp.GetMessage() != "Modification success" {
		t.Fatalf("Unexpected values %v", modResp)
	}
}

func TestRailBookingService_CancelBooking(t *testing.T) {

	// create a rail booking service instance with a reference to the db and seat allocation service
	db := NewDB()
	st := NewAllocateSeat(db)
	railBookingService := NewRailBookingService(db, st)
	conn := newServer(t, func(srv *grpc.Server) {
		bk.RegisterBookingServiceServer(srv, &railBookingService)
	})

	client := bk.NewBookingServiceClient(conn)
	res, err := client.RailBooking(context.Background(), &bk.BookingRequest{From: "Chennai", To: "Delhi", Price: 20, User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}})
	if err != nil {
		t.Fatalf("client.RailBooking %v %v", err, res)
	}
	cancelBookingsRequest := &bk.CancelBookingRequest{User: &ur.User{FirstName: "Ragav", LastName: "Ram", Email: "ragav@gmail.com"}}
	log.Println("Get Rail Booking by User call check for wrong user and error is thrown")
	cancelResp, err1 := client.CancelBooking(context.Background(), cancelBookingsRequest)
	if err1 != nil {
		t.Fatalf("client.ModifySeatByUser %v", err1)
	}
	if cancelResp.GetMessage() != "Cancelation success" {
		t.Fatalf("Unexpected values %v", cancelResp)
	}
}
