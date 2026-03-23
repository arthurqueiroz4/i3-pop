#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="${HOME}/.local/bin"
SYSTEMD_USER_DIR="${HOME}/.config/systemd/user"
SERVICE_NAME="i3-pop.service"

mkdir -p "${BIN_DIR}" "${SYSTEMD_USER_DIR}"

echo "Building i3-pop binary..."
go build -o "${BIN_DIR}/i3-pop" "${ROOT_DIR}"

echo "Installing user service..."
install -m 0644 "${ROOT_DIR}/systemd/${SERVICE_NAME}" "${SYSTEMD_USER_DIR}/${SERVICE_NAME}"

echo "Reloading and starting systemd user service..."
systemctl --user daemon-reload
systemctl --user enable --now "${SERVICE_NAME}"

echo
echo "Daemon installed and running."
echo "Status: systemctl --user status ${SERVICE_NAME}"
echo "Logs:   journalctl --user -u ${SERVICE_NAME} -f"
