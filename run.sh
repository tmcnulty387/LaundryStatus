#!/bin/bash

cleanup() {
    podman compose down
    kill 0
}

trap cleanup SIGINT

podman compose up -d db

air &
npm --prefix frontend run dev &

wait
