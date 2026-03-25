-- name: GetMachines :many
SELECT m.id,
    m.is_operational,
    m.is_washer,
    r.end_at
FROM machine m
    LEFT JOIN reservation r ON r.room_slug = m.room_slug
    AND r.machine_id = m.id
WHERE m.room_slug = $1
ORDER BY m.id;

-- name: ReserveMachine :exec
INSERT INTO reservation (room_slug, machine_id, end_at, phone_number)
VALUES ($1, $2, $3, $4) ON CONFLICT (room_slug, machine_id) DO
UPDATE
SET end_at = EXCLUDED.end_at,
    phone_number = EXCLUDED.phone_number;

-- name: SetMachineAvailable :exec
UPDATE machine
SET is_operational = TRUE
WHERE room_slug = $1
    AND id = $2;

-- name: SetMachineOutOfOrder :exec
UPDATE machine
SET is_operational = FALSE
WHERE room_slug = $1
    AND id = $2;

-- name: ClearReservation :exec
DELETE FROM reservation
WHERE room_slug = $1
    AND machine_id = $2;

-- name: ReservationExists :one
SELECT EXISTS (
        SELECT 1
        FROM reservation
        WHERE room_slug = $1
            AND machine_id = $2
    ) AS exists;
