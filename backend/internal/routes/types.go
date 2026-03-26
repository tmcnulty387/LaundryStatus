package routes

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type roomSlug string

const (
	ResidenceHallA roomSlug = "residence-hall-a"
	ResidenceHallB roomSlug = "residence-hall-b"
	ResidenceHallC roomSlug = "residence-hall-c"
	KateGleason    roomSlug = "kate-gleason"
	Gibson         roomSlug = "gibson"
	Peterson       roomSlug = "peterson"
	SolHeumann     roomSlug = "sol-heumann"
)

var validRoomSlugs = map[roomSlug]struct{}{
	ResidenceHallA: {},
	ResidenceHallB: {},
	ResidenceHallC: {},
	KateGleason:    {},
	Gibson:         {},
	Peterson:       {},
	SolHeumann:     {},
}

func (r roomSlug) isValid() bool {
	_, ok := validRoomSlugs[r]
	return ok
}

func (r roomSlug) ToName() string {
	switch r {
	case ResidenceHallA:
		return "Residence Hall A"
	case ResidenceHallB:
		return "Residence Hall B"
	case ResidenceHallC:
		return "Residence Hall C"
	case KateGleason:
		return "Kate Gleason"
	case Gibson:
		return "Gibson"
	case Peterson:
		return "Peterson"
	case SolHeumann:
		return "Sol Heumann"
	default:
		return string(r)
	}
}

type machine struct {
	ID            int32  `json:"id"`
	IsOperational bool   `json:"is_operational"`
	IsWasher      bool   `json:"is_washer"`
	EndAt         *int64 `json:"end_at"`
}

type machinesResponse struct {
	Washers []machine `json:"washers"`
	Dryers  []machine `json:"dryers"`
}

func getTime(value pgtype.Timestamptz) *int64 {
	// If null or in the past, reservation does not exist.
	if !value.Valid || value.Time.Before(time.Now()) {
		return nil
	}
	epoch := value.Time.UnixMilli() // JS Date.now() returns milliseconds
	return &epoch
}

func epochToTime(value int64) time.Time {
	return time.UnixMilli(value)
}

type createReservationBody struct {
	EndAt       int64   `json:"end_at"`
	PhoneNumber *string `json:"phone_number"`
}

type reservationParams struct {
	RoomSlug    roomSlug
	MachineID   int32
	EndAt       time.Time
	PhoneNumber *string
	IsWasher    bool
}
