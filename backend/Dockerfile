# syntax=docker/dockerfile:1

FROM golang:1.26 AS build
WORKDIR /app

COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /app/backend
RUN go mod download

COPY backend/. .

# RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /laundry-backend ./cmd

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=build /laundry-backend /laundry-backend

EXPOSE 8080
ENTRYPOINT ["/laundry-backend"]
