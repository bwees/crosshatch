#!/bin/bash
set -eu

GO2RTC_BIN="${GO2RTC_BIN:-/usr/local/bin/go2rtc}"
GO2RTC_CONFIG_DIR="${GO2RTC_CONFIG_DIR:-/app/server/data/go2rtc}"
SERVER_ENTRY="${SERVER_ENTRY:-/app/server/dist/src/main.js}"

mkdir -p "$GO2RTC_CONFIG_DIR"

(cd "$GO2RTC_CONFIG_DIR" && exec "$GO2RTC_BIN") &
GO2RTC_PID=$!

(cd /app/server && exec node "$SERVER_ENTRY") &
SERVER_PID=$!

cleanup() {
  kill -TERM "$GO2RTC_PID" "$SERVER_PID" 2>/dev/null || true
}
trap cleanup TERM INT

set +e
wait -n
EXIT_CODE=$?
set -e

cleanup
wait 2>/dev/null || true
exit "$EXIT_CODE"
