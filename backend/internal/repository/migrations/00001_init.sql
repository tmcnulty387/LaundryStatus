-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS machine (
    room_slug TEXT NOT NULL,
    id INT NOT NULL,
    -- number in room
    is_operational BOOLEAN NOT NULL,
    is_washer BOOLEAN NOT NULL,
    PRIMARY KEY (room_slug, id)
);
CREATE TABLE IF NOT EXISTS reservation (
    room_slug TEXT NOT NULL,
    machine_id INT NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,
    phone_number TEXT,
    PRIMARY KEY (room_slug, machine_id),
    -- max 1 reservation per machine
    FOREIGN KEY (room_slug, machine_id) REFERENCES machine (room_slug, id) ON UPDATE CASCADE ON DELETE CASCADE
);
-- Seed initial machines
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'sol-heumann',
    gs,
    TRUE,
    TRUE
FROM generate_series(1, 22) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'sol-heumann',
    gs,
    TRUE,
    FALSE
FROM generate_series(23, 42) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'gibson',
    gs,
    TRUE,
    TRUE
FROM generate_series(1, 10) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'gibson',
    gs,
    TRUE,
    FALSE
FROM generate_series(11, 22) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'peterson',
    gs,
    TRUE,
    TRUE
FROM generate_series(1, 20) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'peterson',
    gs,
    TRUE,
    FALSE
FROM generate_series(21, 40) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'kate-gleason',
    gs,
    TRUE,
    TRUE
FROM generate_series(1, 22) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'kate-gleason',
    gs,
    TRUE,
    FALSE
FROM generate_series(23, 44) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'residence-hall-a',
    gs,
    TRUE,
    TRUE
FROM generate_series(1, 6) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'residence-hall-a',
    gs,
    TRUE,
    FALSE
FROM generate_series(7, 12) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'residence-hall-b',
    gs,
    TRUE,
    TRUE
FROM generate_series(1, 6) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'residence-hall-b',
    gs,
    TRUE,
    FALSE
FROM generate_series(7, 12) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'residence-hall-c',
    gs,
    TRUE,
    TRUE
FROM generate_series(1, 6) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
INSERT INTO machine (room_slug, id, is_operational, is_washer)
SELECT 'residence-hall-c',
    gs,
    TRUE,
    FALSE
FROM generate_series(7, 12) AS gs ON CONFLICT (room_slug, id) DO NOTHING;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reservation;
DROP TABLE IF EXISTS machine;
-- +goose StatementEnd