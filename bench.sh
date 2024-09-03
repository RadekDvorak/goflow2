#!/bin/bash
set -euo pipefail

# english errors
export LC_ALL=C

FORMAT=${1:-json}
NAME=${2:-cpu.pprof}

FILE="$PWD/$NAME"

./dist/goflow2 -cpuprofile "${FILE}" -format "${FORMAT}" &
CHILD_PID=$!
( sleep 10s && kill -s INT $CHILD_PID || true ) &
socat -u FILE:raw_netflow_UDP_DATA.log UDP-SENDTO:127.0.0.1:2055

wait
go tool pprof -http=:8080 "${FILE}" &
CHILD_PID=$!
trap 'kill -s TERM $CHILD_PID' EXIT

wait