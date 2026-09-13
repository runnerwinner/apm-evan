#!/bin/sh
set -eu

if [ "${ENABLE_FILEBEAT:-true}" = "true" ]; then
  filebeat -e -c /etc/filebeat/filebeat.yml -strict.perms=false &
fi

exec /app/app
