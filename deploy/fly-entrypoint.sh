#!/bin/sh
set -e

if [ -n "$GCP_KEY_B64" ]; then
  mkdir -p /tmp/gcp
  echo "$GCP_KEY_B64" | base64 -d > /tmp/gcp/key.json
  export GOOGLE_APPLICATION_CREDENTIALS=/tmp/gcp/key.json
fi

exec "$@"
