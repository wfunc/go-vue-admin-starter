#!/bin/sh
set -e

if [ "$(id -u)" = "0" ]; then
    mkdir -p /app/data
    chown -R app:app /app/data 2>/dev/null || true
    exec su-exec app "$0" "$@"
fi

if [ "${1#-}" != "$1" ]; then
    set -- /app/server "$@"
fi

exec "$@"
