# LaundryStatus

A website made to help RIT students living in residence halls keep track of
their laundry. This website provides a real-time display of each laundry room,
showing which machines are out of order, in use, and available. Students will
be able to reserve machines using this website and enter their phone number to
recieve a SMS notification when their laundry is ready.

## Build Instructions
Backend:
```bash
cd backend/
# Add env variables
cp .env.example .env
# Start db
podman compose up db
# View db from CLI
PGPASSWORD=postgres psql -h localhost -U postgres -d laundry
# Initialize db
source .env
goose up
# Clear/drop db
goose down
# Create new migration
goose -dir backend/internal/adapters/migrations create {migration_name} sql
# Stop db
podman compose down db
# Start backend (air)
air --build.cmd "go build -o ./tmp/api ./cmd" --build.bin "./tmp/api"
# Start backend (podman)
podman compose up backend
```

frontend:
```
cd frontend/
npm i
npm run dev
```

## Services

| Service  | URL                    | Description           |
|----------|------------------------|-----------------------|
| Frontend | http://localhost:3000  | React+Vite web app    |
| Backend  | http://localhost:8080  | Go RESTful API        |
| Database | http://localhost:5432  | PostgreSQL 16         |


