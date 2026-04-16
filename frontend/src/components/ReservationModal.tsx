import { useState, type SubmitEvent } from "react";
import "react-phone-number-input/style.css";
import PhoneInput from "react-phone-number-input";
import type { Machine } from "../utils";
import { useCookies } from "react-cookie";

interface CookieValues {
  phoneNumber?: string;
}

type ReservationModalProps = {
  machine: Machine;
  minutes: number;
  roomSlug: string;
  onSuccess: () => void;
  onClose: () => void;
};

// PUT /api/rooms/{room_slug}/machines/{machine_id}/reserve
// {
//     "end_at": 1700000000000,
//     "phone_number": "+15551234567" // optional
// }
function handleSubmit(
  e: SubmitEvent<HTMLFormElement>,
  minutes: number,
  roomSlug: string,
  machineId: number,
  phoneNumber: string | undefined,
  onSuccess: () => void,
  onClose: () => void,
) {
  e.preventDefault();

  const apiUrl = import.meta.env.VITE_API_URL ?? "";
  const endAt = Date.now() + minutes * 60000;

  fetch(`${apiUrl}/api/rooms/${roomSlug}/machines/${machineId}/reserve`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      end_at: endAt,
      phone_number: phoneNumber,
    }),
  })
    .then((response) => {
      if (!response.ok) throw new Error("Failed to reserve machine");
      onSuccess();
      onClose();
    })
    .catch((error) => {
      console.error(error);
      alert("Error reserving machine");
    });
}

function ReservationModal({
  machine,
  minutes,
  roomSlug,
  onSuccess,
  onClose,
}: ReservationModalProps) {
  const [cookies, setCookies, removeCookie] = useCookies<
    "phoneNumber",
    CookieValues
  >(["phoneNumber"]);
  const [smsConsent, setSmsConsent] = useState(false);

  const canSubmit = (cookies.phoneNumber ?? "").trim() === "" || smsConsent;

  return (
    <div
      className="modal-backdrop"
      role="dialog"
      aria-modal="true"
      aria-label="Reservation confirmation"
      onClick={onClose}
    >
      <div className="modal-card" onClick={(event) => event.stopPropagation()}>
        <h2 className="confirmation-title">
          Reserving {minutes} minutes on{" "}
          {machine.is_washer ? "Washer" : "Dryer"} #{machine.id}
        </h2>

        <form
          className="confirmation-form"
          onSubmit={(e) =>
            handleSubmit(
              e,
              minutes,
              roomSlug,
              machine.id,
              cookies.phoneNumber,
              onSuccess,
              onClose,
            )
          }
        >
          <div className="phone-input-wrapper">
            <label htmlFor="phone-number" className="phone-label">
              Phone Number (optional)
            </label>
            <PhoneInput
              id="phone-number"
              defaultCountry="US"
              placeholder="Enter phone number"
              autoComplete="tel"
              value={cookies.phoneNumber}
              onChange={(value) =>
                value
                  ? setCookies("phoneNumber", value, { path: "/" })
                  : removeCookie("phoneNumber")
              }
            />
            <div className="sms-consent-container">
              <input
                type="checkbox"
                id="sms-consent"
                className="sms-consent-checkbox"
                checked={smsConsent}
                onChange={(e) => setSmsConsent(e.target.checked)}
              />
              <label htmlFor="sms-consent" className="sms-consent-label">
                I agree to receive a one-time SMS text message from
                LaundryStatus when my laundry is complete. Message and data
                rates may apply.
              </label>
            </div>
          </div>

          <div className="confirmation-buttons">
            <button
              type="button"
              className="confirmation-button cancel-button"
              onClick={onClose}
            >
              Cancel
            </button>
            <button
              disabled={!canSubmit}
              type="submit"
              className="confirmation-button submit-button"
            >
              Submit
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default ReservationModal;
