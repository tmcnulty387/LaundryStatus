import {
  useEffect,
  useState,
  type Dispatch,
  type SetStateAction,
  type SubmitEvent,
} from "react";
import { getMachineStatusLabel, getMachineTimer, type Machine } from "../utils";
import ReservationModal from "./ReservationModal";

const apiUrl = import.meta.env.VITE_API_URL ?? "";

type MachineModalProps = {
  machine: Machine;
  roomSlug: string;
  onMachineUpdated: () => void;
  onClose: () => void;
};

function handleReservationSubmit(
  e: SubmitEvent<HTMLFormElement>,
  setShowReservation: Dispatch<SetStateAction<boolean>>,
) {
  e.preventDefault();
  const input = e.currentTarget.elements.namedItem(
    "reserve-minutes",
  ) as HTMLInputElement;
  const minutes = input.valueAsNumber;
  input.setCustomValidity("");

  if (Number.isNaN(minutes) || minutes < 1 || minutes > 180) {
    input.setCustomValidity(
      "Please enter a valid number of minutes between 1 and 180.",
    );
    input.reportValidity();
    return;
  }

  setShowReservation(true);
}

// POST /api/rooms/{room_slug}/machines/{machine_id}/available
function handleSetAvailable(
  roomSlug: string,
  machineId: number,
  onMachineUpdated: () => void,
  onClose: () => void,
) {
  fetch(`${apiUrl}/api/rooms/${roomSlug}/machines/${machineId}/available`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  })
    .then((response) => {
      if (!response.ok) throw new Error("Failed to mark as available");
      onMachineUpdated();
      onClose();
    })
    .catch((error) => {
      console.error(error);
      alert("Error marking machine as available");
    });
}

// POST /api/rooms/{room_slug}/machines/{machine_id}/out-of-order
function handleSetOutOfOrder(
  roomSlug: string,
  machineId: number,
  onMachineUpdated: () => void,
  onClose: () => void,
) {
  fetch(`${apiUrl}/api/rooms/${roomSlug}/machines/${machineId}/out-of-order`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
  })
    .then((response) => {
      if (!response.ok) throw new Error("Failed to mark as out of order");
      onMachineUpdated();
      onClose();
    })
    .catch((error) => {
      console.error(error);
      alert("Error marking machine as out of order");
    });
}

function MachineModal({
  machine,
  roomSlug,
  onMachineUpdated,
  onClose,
}: MachineModalProps) {
  const [nowSec, setNowSec] = useState(Math.floor(Date.now() / 1000));
  const status = getMachineStatusLabel(machine, nowSec);
  const timer = getMachineTimer(machine, nowSec);
  const [showReservation, setShowReservation] = useState(false);
  const [reserveMinutesInput, setReserveMinutesInput] = useState("30");
  const parsedReserveMinutes = Number(reserveMinutesInput);
  const isReserveMinutesValid =
    Number.isInteger(parsedReserveMinutes) &&
    parsedReserveMinutes >= 1 &&
    parsedReserveMinutes <= 180;
  const reserveMinutesError = isReserveMinutesValid
    ? null
    : "Enter a whole number of minutes between 1 and 180.";
  const isAlreadyAvailable = status === "Available";
  const isAlreadyOutOfOrder = status === "Out of Order";

  useEffect(() => {
    const intervalId = window.setInterval(() => {
      setNowSec(Math.floor(Date.now() / 1000));
    }, 1000);

    return () => {
      window.clearInterval(intervalId);
    };
  }, []);

  if (showReservation)
    return (
      <ReservationModal
        machine={machine}
        minutes={parsedReserveMinutes}
        roomSlug={roomSlug}
        onSuccess={onMachineUpdated}
        onClose={onClose}
      />
    );
  return (
    <div
      className="modal-backdrop"
      role="dialog"
      aria-modal="true"
      aria-label={`Machine ${machine.id} details`}
      onClick={onClose}
    >
      <div className="modal-card" onClick={(e) => e.stopPropagation()}>
        <div className="modal-header">
          <h2>
            {machine.is_washer ? `Washer` : `Dryer`} {machine.id}
          </h2>
          <div className="status-container">
            <p className="modal-status">{status}</p>
            {timer && <p className="modal-timer">{timer}</p>}
          </div>
        </div>

        <form
          className="reservation-form"
          onSubmit={(e) => handleReservationSubmit(e, setShowReservation)}
        >
          <div className="input-wrapper">
            <input
              id="reserve-minutes"
              name="reserve-minutes"
              type="number"
              min="1"
              max="180"
              step="1"
              required
              value={reserveMinutesInput}
              onChange={(e) => setReserveMinutesInput(e.currentTarget.value)}
              aria-invalid={!isReserveMinutesValid}
              aria-describedby={
                reserveMinutesError ? "reserve-minutes-error" : undefined
              }
              className="minutes-input"
            />
            <span className="minutes-suffix">Min</span>
          </div>
          {reserveMinutesError && (
            <p
              id="reserve-minutes-error"
              className="error-subtext"
              role="alert"
            >
              {reserveMinutesError}
            </p>
          )}
          <button
            type="submit"
            className="reserve-button"
            disabled={!isReserveMinutesValid}
          >
            Reserve
          </button>
        </form>

        <div className="status-buttons">
          <button
            type="button"
            className="status-button status-available"
            onClick={() =>
              handleSetAvailable(
                roomSlug,
                machine.id,
                onMachineUpdated,
                onClose,
              )
            }
            disabled={isAlreadyAvailable}
          >
            Available
          </button>
          <button
            type="button"
            className="status-button status-out-of-order"
            onClick={() =>
              handleSetOutOfOrder(
                roomSlug,
                machine.id,
                onMachineUpdated,
                onClose,
              )
            }
            disabled={isAlreadyOutOfOrder}
          >
            Out of Order
          </button>
        </div>
      </div>
    </div>
  );
}

export default MachineModal;
