#!/usr/bin/env bash
# install.sh — Install the latest gow CLI

set -e

REPO="mechneerd/gow"
INSTALL_DIR="/usr/local/bin"

echo "Installing gow from $REPO..."

OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map to goreleaser naming conventions
if [ "$ARCH" = "x86_64" ]; then
  ARCH="x86_64"
elif [ "$ARCH" = "aarch64" ]; then
  ARCH="arm64"
fi

# Capitalize OS for goreleaser (linux -> Linux, darwin -> Darwin)
OS_TITLE="$(echo "${OS:0:1}" | tr '[:lower:]' '[:upper:]')${OS:1}"

# goreleaser uses .tar.gz for Linux/macOS
EXT="tar.gz"
if [ "$OS" = "windows" ]; then
  EXT="zip"
fi

BINARY="gow_${OS_TITLE}_${ARCH}.${EXT}"
URL="https://github.com/${REPO}/releases/latest/download/${BINARY}"

echo "Downloading ${BINARY}..."

curl -sSfL "$URL" -o /tmp/gow-archive

if [ "$EXT" = "tar.gz" ]; then
  tar -xzf /tmp/gow-archive -C /tmp
elif [ "$EXT" = "zip" ]; then
  unzip -o /tmp/gow-archive -d /tmp
fi

chmod +x /tmp/gow

sudo mv /tmp/gow "$INSTALL_DIR/gow"

rm -f /tmp/gow-archive /tmp/gow

echo "gow installed successfully!"
gow --version
