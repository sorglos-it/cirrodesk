#!/bin/sh
# Test service for contract V3: echoes stdin to stdout, reports who asked on
# stderr and exits with the code given in CD_PARAMS, e.g. {"exit":2}.
set -eu
printf '%s %s %s %s\n' "$CD_APP" "$CD_SERVICE" "$CD_ACTION" "$CD_USER" >&2
cat
code=$(printf '%s' "${CD_PARAMS:-}" | sed -n 's/.*"exit":\([0-9][0-9]*\).*/\1/p')
exit "${code:-0}"
