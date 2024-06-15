package internal

import (
	"fmt"
	"log"
	"math/rand"
	st "rail-booking/protogen/seats"
)

type AllocateSeat struct {
	db *DB
}

// NewAllocateSeat creates a new instance
func NewAllocateSeat(db *DB) *AllocateSeat {
	return &AllocateSeat{db: db}
}

// For random allocation of a Section for a booking
var sectionRunes = []rune("AB")

// Find the available random seat and random section for a booking comparing the existing bookings in db (slice)
func (as *AllocateSeat) allocateSeat() (*st.Seats, error) {
	log.Println("Service Entry : Seat allocation call")

	availSeat := &st.Seats{}

	boolcheck := false
	var err error
	for !boolcheck {
		availSeat.Section = string(sectionRunes[rand.Intn(len(sectionRunes))])
		availSeat.Seat = int32(rand.Intn(MaxSeatsPerSection))
		boolcheck, err = as.db.CheckBookingBySeatAndSection(availSeat)
		if boolcheck && err != nil && err.Error() == "seats full" {
			log.Println("Service Exit : Seat allocation call, seats full")
			return nil, fmt.Errorf("seats full, try later")
		}
	}
	log.Println("Service Exit : Seat allocation call")
	return availSeat, nil
}
