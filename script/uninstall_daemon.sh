#!/usr/bin/env bash
set -euo pipefail

SERVICE_NAME="i3-pop.service"
SYSTEMD_USER_DIR="${HOME}/.config/systemd/user"
SERVICE_FILE="${SYSTEMD_USER_DIR}/${SERVICE_NAME}"
BIN_FILE="${HOME}/.local/bin/i3-pop"

echo "Stopping and disabling service (if active)..."
systemctl --user disable --now "${SERVICE_NAME}" >/dev/null 2>&1 || true

if [ -f "${SERVICE_FILE}" ]; then
  rm -f "${SERVICE_FILE}"
fi

systemctl --user daemon-reload

if [ -f "${BIN_FILE}" ]; then
  rm -f "${BIN_FILE}"
fi

echo "Daemon uninstalled."
