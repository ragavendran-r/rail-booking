package internal

import (
	"fmt"
	"log"
	"sync/atomic"

	bk "rail-booking/protogen/booking"
	st "rail-booking/protogen/seats"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DB struct {
	collection   []*bk.Booking
	sectACounter atomic.Int32
	sectBCounter atomic.Int32
}

const MaxSeatsPerSection = 72
const BKNGS_NOT_FOUND string = "bookings not found"

// NewDB creates a new array to mimic the behaviour of a in-memory database
func NewDB() *DB {
	return &DB{
		collection: make([]*bk.Booking, 0),
	}
}

// RailBooking adds a new ticket to the DB collection. Returns an error on duplicate booking for a user
// Duplicate is avoided to help with easy modification and cancelation of bookings for a user
func (d *DB) RailBooking(booking *bk.Booking) (*bk.Booking, error) {
	log.Printf("DB Entry : booking a ticket and storing in slice db")
	for _, b := range d.collection {
		if b.GetUser().Email == booking.GetUser().Email {
			log.Printf("Duplicate found in db")
			return nil, fmt.Errorf("duplicate request for user: %s", booking.GetUser().Email)
		}
	}
	d.collection = append(d.collection, booking)
	log.Println("After booking request in db", d)
	if booking.Section == "A" {
		d.sectACounter.Add(1)
	} else {
		d.sectBCounter.Add(1)
	}
	log.Println("DB Exit : booking a ticket and storing in slice db")
	return booking, nil
}

// GetBookingByUser returns a booking for the user
func (d *DB) GetBookingByUser(bookingreq *bk.GetBookingByUserRequest) (*bk.Booking, error) {
	log.Println("DB Entry : Get booking ticket for a user")
	for _, b := range d.collection {
		if b.User.Email == bookingreq.User.Email {
			return b, nil
		}
	}
	log.Println("DB Exit : Get booking ticket for a user")
	return nil, fmt.Errorf(BKNGS_NOT_FOUND)
}

// GetAllBookings returns list of all bookings
func (d *DB) GetAllBookings() (*bk.BookingList, error) {
	log.Println("DB Entry : Get All booking tickets")
	bookingList := &bk.BookingList{}
	if len(d.collection) > 0 {
		bookingList.Bookings = d.collection
	} else {
		log.Println("DB Entry : Get All booking tickets, no bookings")
		return nil, fmt.Errorf(BKNGS_NOT_FOUND)
	}
	log.Println("DB Exit : Get All booking tickets")
	return bookingList, nil
}

// GetSectionBookings returns list of all bookings for a section
func (d *DB) GetSectionBookings(req *bk.GetSectionBookingsRequest) (*bk.BookingList, error) {
	log.Println("DB Entry : Get All booking tickets for a section")
	bookingList := &bk.BookingList{}
	for _, b := range d.collection {
		if b.Section == req.GetSection() {
			bookingList.Bookings = append(bookingList.Bookings, b)
		}
	}
	log.Println("DB Exit : Get All booking tickets for a section")
	return bookingList, nil
}

// cancelBooking for the user
func (d *DB) CancelBooking(cancelBookingRequest *bk.CancelBookingRequest) (*bk.CancelBookingResponse, error) {
	log.Println("DB Entry : Cancel booking tickets for a user")
	for i, b := range d.collection {
		if b.User.Email == cancelBookingRequest.User.Email {
			d.collection[i] = d.collection[len(d.collection)-1]
			d.collection = d.collection[:len(d.collection)-1]
			if b.Section == "A" {
				d.sectACounter.Add(-1)
			} else {
				d.sectBCounter.Add(-1)
			}
			log.Println("DB Exit : Canceling booking in db")
			return &bk.CancelBookingResponse{Message: "Cancelation success"}, nil
		}

	}
	log.Println("DB Exit : Cancel booking tickets for a user")
	//return nil, fmt.Errorf("booking not found")
	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("booking not found: %s\n", cancelBookingRequest.User.Email))
}

// ModifySeatByUser modifies the seat for the user
func (d *DB) ModifySeatByUser(seatmodRequest *bk.SeatModificationRequest) (*bk.SeatModificationResponse, error) {
	log.Println("DB Entry : Modify  booking ticket for a user")
	for i, b := range d.collection {

		if b.User.Email == seatmodRequest.User.Email {
			log.Printf("Updating seat in db existing sec: %s, incoming sec: %s", b.Section, seatmodRequest.Section)
			if b.Section != seatmodRequest.Section {
				log.Printf("Inside section compare existing sec: %s, incoming sec: %s", b.Section, seatmodRequest.Section)
				if b.Section == "A" {
					log.Println("Modify booking section counter change for A")
					d.sectACounter.Add(-1)
					d.sectBCounter.Add(1)
				} else {
					log.Println("Modify booking section counter change for B")
					d.sectBCounter.Add(-1)
					d.sectACounter.Add(1)
				}
			}
			d.collection[i].Section = seatmodRequest.Section
			d.collection[i].Seat = seatmodRequest.Seat
			log.Printf("After Updating seat in db existing sec: %s, incoming sec: %s", b.Section, seatmodRequest.Section)
			log.Println("DB Entry : Modify  booking ticket for a user")
			return &bk.SeatModificationResponse{Message: "Modification success"}, nil
		}
	}
	log.Println("DB Exit : Modify  booking ticket for a user")
	return nil, fmt.Errorf("booking not found")
}

// CheckBookingBySeatAndSection checks a booking for the seat and section
func (d *DB) CheckBookingBySeatAndSection(bookingreq *st.Seats) (bool, error) {
	log.Println("DB Entry : CheckBooking By Seat And Section ")
	if d.sectACounter.Load() > MaxSeatsPerSection && d.sectBCounter.Load() > MaxSeatsPerSection {
		log.Println("DB Exit : CheckBooking By Seat And Section, seats full")
		return true, fmt.Errorf("seats full")
	}
	for _, b := range d.collection {
		log.Printf("Search request in db")
		if b.Seat == bookingreq.Seat && b.Section == bookingreq.Section {
			return false, nil
		}
	}
	log.Println("DB Exit : CheckBooking By Seat And Section, no bookings ")
	return true, fmt.Errorf(BKNGS_NOT_FOUND)
}
