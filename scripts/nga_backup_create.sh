#!/usr/bin/env bash

set -euo pipefail

DB_FILE="/path/to/nga.db"

DATE=$(date +"%Y-%m-%d_%H-%M")
BACKUP_FILE="/tmp/nga.${DATE}.db"

sqlite3 "$DB_FILE" ".backup '$BACKUP_FILE'"

echo "$BACKUP_FILE"