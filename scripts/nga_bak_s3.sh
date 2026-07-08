#!/usr/bin/env bash

set -euo pipefail

BACKUP_FILE=$(/path/to/nga_backup_create.sh)

NAME=$(basename "$BACKUP_FILE")

aws s3 cp \
    "$BACKUP_FILE" \
    "s3://YOUR-BUCKET-NAME/nga/$NAME" \
    --endpoint-url https://YOUR-URL


rm -f "$BACKUP_FILE"