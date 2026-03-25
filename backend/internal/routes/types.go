package routes

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type RoomSlug string

const (
	ResidenceHallA RoomSlug = "residence-hall-a"
	ResidenceHallB RoomSlug = "residence-hall-b"
	ResidenceHallC RoomSlug = "residence-hall-c"
	Gleason        RoomSlug = "kate-gleason"
	Gibson         RoomSlug = "gibson"
	Peterson       RoomSlug = "peterson"
	SolHeumann     RoomSlug = "sol-heumann"
)

var validRoomSlugs = map[RoomSlug]struct{}{
	ResidenceHallA: {},
	ResidenceHallB: {},
	ResidenceHallC: {},
	Gleason:        {},
	Gibson:         {},
	Peterson:       {},
	SolHeumann:     {},
}

func (r RoomSlug) valid() bool {
	_, ok := validRoomSlugs[r]
	return ok
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
	epoch := timeToEpoch(value.Time)
	return &epoch
}

func epochToTime(value int64) time.Time {
	return time.Unix(value, 0)
}

func timeToEpoch(value time.Time) int64 {
	return value.Unix()
}

type createReservationBody struct {
	EndAt       int64   `json:"end_at"`
	PhoneNumber *string `json:"phone_number"`
}

type reservationParams struct {
	RoomSlug    RoomSlug
	MachineID   int32
	EndAt       time.Time
	PhoneNumber *string
}
