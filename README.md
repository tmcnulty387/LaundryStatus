# LaundryStatus

A website made to help RIT students living in residence halls keep track of
their laundry. This website provides a real-time display of each laundry room,
showing which machines are out of order, in use, and available. Students will
be able to reserve machines using this website and enter their phone number to
recieve a SMS notification when their laundry is ready.

## Build Instructions
```bash
# Add env variables
cp .env.example .env
# Start db
podman compose up db
# Start backend
air
# Start frontend
npm --prefix frontend run dev
```

## Database
```bash
# Initialize db
source .env
goose up
# Clear/drop db
goose down
# Create new migration
goose -dir backend/internal/adapters/migrations create {migration_name} sql
# View db from CLI
PGPASSWORD=postgres psql -h localhost -U postgres -d laundry
```
