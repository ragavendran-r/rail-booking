package internal

import (
	"context"
	"fmt"
	"log"

	bk "rail-booking/protogen/booking"

	"github.com/google/uuid"
)

// RailBookingService should implement the BookingServiceServer interface generated from grpc.
// UnimplementedBookingServiceServer must be embedded to have forward compatible implementations.

type RailBookingService struct {
	db *DB
	st *AllocateSeat
	bk.UnimplementedBookingServiceServer
}

// Create a new rail booking service
func NewRailBookingService(db *DB, st *AllocateSeat) RailBookingService {
	return RailBookingService{db: db, st: st}
}

// RailBooking adds a new ticket to the DB collection. Returns an error on duplicate ids
func (rbs *RailBookingService) RailBooking(_ context.Context, bookingreq *bk.BookingRequest) (*bk.Booking, error) {
	log.Println("Service Entry: booking a ticket and storing in db")

	seats, err := rbs.st.allocateSeat()
	if err != nil {
		return nil, err
	}
	log.Println("Allocated seats : ", seats)
	booking := &bk.Booking{}
	booking.Id = uuid.New().String()
	booking.From = bookingreq.GetFrom()
	booking.To = bookingreq.GetTo()
	booking.User = bookingreq.GetUser()
	booking.Price = bookingreq.GetPrice()
	booking.Seat = seats.GetSeat()
	booking.Section = seats.GetSection()
	booking, err1 := rbs.db.RailBooking(booking)
	if err1 != nil {
		return nil, err1
	}
	log.Println("Service Exit: After booking request in db", booking)

	return booking, nil
}

// GetBookingByUser returns a booking for the user
func (rbs *RailBookingService) GetBookingByUser(_ context.Context, bookingreq *bk.GetBookingByUserRequest) (*bk.Booking, error) {
	log.Println("Service Entry: Get booking details for a user")
	booking, err := rbs.db.GetBookingByUser(bookingreq)
	if err != nil {
		return nil, err
	}
	log.Println("Service Exit: booking a ticket and storing in db")
	return booking, nil
}

// GetAllBookings returns list of all bookings
func (rbs *RailBookingService) GetAllBookings(_ context.Context, emp *bk.Empty) (*bk.BookingList, error) {
	log.Println("Service Entry: Get All Rail bookings ")
	bookingList, err := rbs.db.GetAllBookings()
	if err != nil {
		return nil, err
	}
	log.Println("Service Exit: Get All Rail bookings ")
	return bookingList, nil
}

// cancelBooking for the user
func (rbs *RailBookingService) CancelBooking(_ context.Context, cancelBookingRequest *bk.CancelBookingRequest) (*bk.CancelBookingResponse, error) {
	log.Println("Service Entry: Cancel Rail booking for a user")
	cancelbooking, err := rbs.db.CancelBooking(cancelBookingRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Service Exit: Cancel Rail booking for a user")
	return cancelbooking, nil
}

// ModifySeatByUser modifies the seat for the user
func (rbs *RailBookingService) ModifySeatByUser(_ context.Context, seatModReq *bk.SeatModificationRequest) (*bk.SeatModificationResponse, error) {
	log.Println("Service Entry: Modify Rail booking for a user")
	modbookingresp, err := rbs.db.ModifySeatByUser(seatModReq)
	if err != nil {
		return nil, err
	}
	log.Println("Service Exit: Modify Rail booking for a user")
	return modbookingresp, nil
}

// GetSectionBookingCount returns count of section bookings
func (rbs *RailBookingService) GetSectionBookingCount(_ context.Context, emp *bk.Empty) (*bk.BookingCountResponse, error) {
	log.Println("Service Entry: Get count of section bookings ")
	bkCountResp := &bk.BookingCountResponse{}
	bkCountResp.Message = fmt.Sprintf("Section A Booking Count: %d Section B Booking Count : %d ", rbs.db.sectACounter.Load(), rbs.db.sectBCounter.Load())
	log.Printf("Service Exit: \n Section A Booking Count: %d and Section B Booking Count : %d ", rbs.db.sectACounter.Load(), rbs.db.sectBCounter.Load())
	return bkCountResp, nil
}
