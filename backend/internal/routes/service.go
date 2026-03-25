package routes

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/tmcnulty387/LaundryStatus/internal/config"
	repo "github.com/tmcnulty387/LaundryStatus/internal/repository/sqlc"
	"github.com/tmcnulty387/LaundryStatus/internal/sms"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEndTimeInPast = errors.New("reservation end time is in the past")

type Service interface {
	GetMachines(ctx context.Context, room RoomSlug) (machinesResponse, error)
	CreateReservation(ctx context.Context, args reservationParams) error
	SetMachineAvailable(ctx context.Context, room RoomSlug, machineID int32) error
	SetMachineOutOfOrder(ctx context.Context, room RoomSlug, machineID int32) error
}

type svc struct {
	repo *repo.Queries
	pool *pgxpool.Pool
	cfg  *config.Config
}

func NewService(repo *repo.Queries, pool *pgxpool.Pool, cfg *config.Config) Service {
	return &svc{repo: repo, pool: pool, cfg: cfg}
}

func (s *svc) GetMachines(ctx context.Context, room RoomSlug) (machinesResponse, error) {
	rows, err := s.repo.GetMachines(ctx, string(room))
	if err != nil {
		return machinesResponse{}, err
	}

	resp := machinesResponse{
		Washers: make([]machine, 0, len(rows)),
		Dryers:  make([]machine, 0, len(rows)),
	}

	for _, row := range rows {
		m := machine{
			ID:            row.ID,
			IsOperational: row.IsOperational,
			IsWasher:      row.IsWasher,
			EndAt:         getTime(row.EndAt),
		}
		if m.IsWasher {
			resp.Washers = append(resp.Washers, m)
		} else {
			resp.Dryers = append(resp.Dryers, m)
		}
	}

	return resp, nil
}

func (s *svc) CreateReservation(ctx context.Context, args reservationParams) error {
	if args.EndAt.Before(time.Now()) {
		log.Printf("Attempted to create reservation with end time in the past: %v", args.EndAt)
		return ErrEndTimeInPast
	}
	setMachineAvailableParams := repo.SetMachineAvailableParams{
		RoomSlug: string(args.RoomSlug),
		ID:       args.MachineID,
	}
	var phoneNumber pgtype.Text
	if args.PhoneNumber != nil {
		phoneNumber = pgtype.Text{String: *args.PhoneNumber, Valid: true}
	} else {
		phoneNumber = pgtype.Text{Valid: false}
	}
	reserveMachineParams := repo.ReserveMachineParams{
		RoomSlug:    string(args.RoomSlug),
		MachineID:   args.MachineID,
		EndAt:       pgtype.Timestamptz{Time: args.EndAt, Valid: true},
		PhoneNumber: phoneNumber,
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)
	if err := qtx.SetMachineAvailable(ctx, setMachineAvailableParams); err != nil {
		return err
	}
	if err := qtx.ReserveMachine(ctx, reserveMachineParams); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	// New context to prevent cancellation when request is closed
	go func() {
		if err := s.awaitReservationEnd(context.Background(), args); err != nil {
			log.Printf("Failed to process reservation end for room %s, machine %d: %v", args.RoomSlug, args.MachineID, err)
		}
	}()
	return nil
}

func (s *svc) SetMachineAvailable(ctx context.Context, room RoomSlug, machineID int32) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)
	if err := qtx.ClearReservation(ctx, repo.ClearReservationParams{
		RoomSlug:  string(room),
		MachineID: machineID,
	}); err != nil {
		return err
	}
	if err := qtx.SetMachineAvailable(ctx, repo.SetMachineAvailableParams{
		RoomSlug: string(room),
		ID:       machineID,
	}); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *svc) SetMachineOutOfOrder(ctx context.Context, room RoomSlug, machineID int32) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)
	if err := qtx.ClearReservation(ctx, repo.ClearReservationParams{
		RoomSlug:  string(room),
		MachineID: machineID,
	}); err != nil {
		return err
	}
	if err := qtx.SetMachineOutOfOrder(ctx, repo.SetMachineOutOfOrderParams{
		RoomSlug: string(room),
		ID:       machineID,
	}); err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *svc) awaitReservationEnd(ctx context.Context, args reservationParams) error {
	delay := time.Until(args.EndAt)
	if delay > 0 {
		timer := time.NewTimer(delay)
		log.Printf("[Starting timer] room: %s, id: %d, end at: %s", args.RoomSlug, args.MachineID, args.EndAt)
		<-timer.C
		log.Printf("Finished timer: %s %d %s", args.RoomSlug, args.MachineID, args.EndAt)
	}

	exists, err := s.repo.ReservationExists(ctx, repo.ReservationExistsParams{
		RoomSlug:  string(args.RoomSlug),
		MachineID: args.MachineID,
	})
	if err != nil {
		return err
	}
	if !exists {
		log.Printf("Reservation for room %s, machine %d no longer exists", args.RoomSlug, args.MachineID)
		return nil
	}
	if err := s.repo.ClearReservation(ctx, repo.ClearReservationParams{
		RoomSlug:  string(args.RoomSlug),
		MachineID: args.MachineID,
	}); err != nil {
		return err
	}
	if args.PhoneNumber != nil {
		log.Printf("Sending SMS to %s for reservation end", *args.PhoneNumber)
		sms.SendSms(s.cfg, *args.PhoneNumber) // checks if sms is enabled internally
	}
	log.Printf("Reservation successfully ended: %s %d %s", args.RoomSlug, args.MachineID, args.EndAt)
	return nil
}
