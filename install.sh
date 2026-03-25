#!/usr/bin/env bash
set -euo pipefail

REPO="orion-rep/gmmit"
BIN_NAME="gmmit"
INSTALL_DIR="/usr/local/bin"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "${OS}" in
  darwin) OS="darwin" ;;
  linux)  OS="linux" ;;
  *)
    echo "[ERROR] Unsupported operating system: ${OS}"
    exit 1
    ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "${ARCH}" in
  x86_64)          ARCH="amd64" ;;
  aarch64 | arm64) ARCH="arm64" ;;
  i386 | i686)     ARCH="386" ;;
  *)
    echo "[ERROR] Unsupported architecture: ${ARCH}"
    exit 1
    ;;
esac

# Resolve latest release tag
echo "[INFO] Fetching latest release..."
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": *"\(.*\)".*/\1/')

if [[ -z "${LATEST}" ]]; then
  echo "[ERROR] Could not determine latest release. Check your internet connection or https://github.com/${REPO}/releases"
  exit 1
fi

echo "[INFO] Latest version: ${LATEST}"

ASSET="${OS}-${ARCH}.tgz"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST}/${ASSET}"

# Download and install
TMP_DIR=$(mktemp -d)
trap 'rm -rf "${TMP_DIR}"' EXIT

echo "[INFO] Downloading ${ASSET}..."
curl -fsSL "${DOWNLOAD_URL}" -o "${TMP_DIR}/${ASSET}"

echo "[INFO] Extracting..."
tar -xzf "${TMP_DIR}/${ASSET}" -C "${TMP_DIR}"

echo "[INFO] Installing to ${INSTALL_DIR}..."
if [[ -w "${INSTALL_DIR}" ]]; then
  install -m 755 "${TMP_DIR}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
else
  sudo install -m 755 "${TMP_DIR}/${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"
fi

echo "[INFO] gmmit ${LATEST} installed successfully!"
echo "[INFO] Run 'gmmit' to get started."
