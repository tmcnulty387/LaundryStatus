package routes

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/tmcnulty387/LaundryStatus/backend/internal/json"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

func NewHandler(s Service) *handler {
	return &handler{
		service: s,
	}
}

// GET /api/rooms/{room_slug}/machines
func (h *handler) GetMachines(w http.ResponseWriter, r *http.Request) {
	room, err := getRoomSlug(w, r)
	if err != nil {
		return
	}
	machines, err := h.service.GetMachines(r.Context(), room)
	if err != nil {
		log.Printf("Failed to get machines: %v", err)
		http.Error(w, "Failed to get machines", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, machines)
}

// PUT /api/rooms/{room_slug}/machines/{machine_id}/reserve
func (h *handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	room, err := getRoomSlug(w, r)
	if err != nil {
		return
	}
	machineId, err := getMachineID(w, r)
	if err != nil {
		return
	}
	var body createReservationBody
	if err := json.Read(r, &body); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isWasher, err := h.service.IsWasher(r.Context(), room, machineId)
	if err != nil {
		log.Printf("Failed to check if machine is washer: %v", err)
		http.Error(w, "Failed to check machine type", http.StatusInternalServerError)
		return
	}

	if err := h.service.CreateReservation(r.Context(), reservationParams{
		RoomSlug:    room,
		MachineID:   machineId,
		EndAt:       epochToTime(body.EndAt),
		PhoneNumber: body.PhoneNumber,
		IsWasher:    isWasher,
	}); err != nil {
		if errors.Is(err, ErrEndTimeInPast) {
			http.Error(w, "Reservation end time must be in the future", http.StatusBadRequest)
			return
		}
		log.Printf("Failed to create reservation: %v", err)
		http.Error(w, "Failed to create reservation", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, nil)
}

// POST /api/rooms/{room_slug}/machines/{machine_id}/available
func (h *handler) SetMachineAvailable(w http.ResponseWriter, r *http.Request) {
	room, err := getRoomSlug(w, r)
	if err != nil {
		return
	}
	machineId, err := getMachineID(w, r)
	if err != nil {
		return
	}
	if err := h.service.SetMachineAvailable(r.Context(), room, machineId); err != nil {
		log.Printf("Failed to set machine available: %v", err)
		http.Error(w, "Failed to set machine available", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, nil)
}

// POST /api/rooms/{room_slug}/machines/{machine_id}/out-of-order
func (h *handler) SetMachineOutOfOrder(w http.ResponseWriter, r *http.Request) {
	room, err := getRoomSlug(w, r)
	if err != nil {
		return
	}
	machineId, err := getMachineID(w, r)
	if err != nil {
		return
	}
	if err := h.service.SetMachineOutOfOrder(r.Context(), room, machineId); err != nil {
		log.Printf("Failed to set machine out of order: %v", err)
		http.Error(w, "Failed to set machine out of order", http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusOK, nil)
}

func getRoomSlug(w http.ResponseWriter, r *http.Request) (roomSlug, error) {
	room := roomSlug(chi.URLParam(r, "room_slug"))
	if !room.isValid() {
		http.Error(w, "Invalid room slug", http.StatusBadRequest)
		return "", errors.New("Invalid room slug")
	}
	return room, nil
}

func getMachineID(w http.ResponseWriter, r *http.Request) (int32, error) {
	machineId, err := strconv.ParseInt(chi.URLParam(r, "machine_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid machine id", http.StatusBadRequest)
		return 0, err
	}
	return int32(machineId), nil
}
