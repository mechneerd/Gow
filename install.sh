#!/usr/bin/env bash
# install.sh — Install the latest gow CLI

set -e

REPO="mechneerd/gow"
INSTALL_DIR="/usr/local/bin"

echo "Installing gow from $REPO..."

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$ARCH" = "x86_64" ]; then
  ARCH="amd64"
elif [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
fi

BINARY="gow-${OS}-${ARCH}"
URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"

curl -sSfL "$URL" -o /tmp/gow
chmod +x /tmp/gow

sudo mv /tmp/gow "$INSTALL_DIR/gow"

echo "✅ gow installed successfully!"
gow --version
