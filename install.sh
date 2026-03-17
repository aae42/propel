#!/bin/sh
set -e

# propel installer script
# Usage: curl -fsSL https://raw.githubusercontent.com/aae42/propel/main/install.sh | bash

GITHUB_DOWNLOAD="https://github.com/aae42/propel/releases/latest/download"
BINARY_NAME="propel"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "Linux" ;;
        Darwin*) echo "Darwin" ;;
        *)       echo "unknown" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)  echo "x86_64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)             echo "unknown" ;;
    esac
}

main() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$OS" = "unknown" ]; then
        echo "Error: Unsupported operating system"
        exit 1
    fi

    if [ "$ARCH" = "unknown" ]; then
        echo "Error: Unsupported architecture"
        exit 1
    fi

    echo "Detected: ${OS} ${ARCH}"

    # Construct download URL (uses GitHub's redirect to latest release)
    FILENAME="${BINARY_NAME}_${OS}_${ARCH}"
    DOWNLOAD_URL="${GITHUB_DOWNLOAD}/${FILENAME}"

    echo "Downloading ${DOWNLOAD_URL}..."

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT

    # Download binary
    curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${BINARY_NAME}"

    # Make executable
    chmod +x "${TMP_DIR}/${BINARY_NAME}"

    # Install
    echo "Installing to ${INSTALL_DIR}/${BINARY_NAME}..."
    if [ -w "$INSTALL_DIR" ]; then
        mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        echo "Need sudo to install to ${INSTALL_DIR}"
        sudo mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    fi

    echo ""
    echo "Successfully installed propel to ${INSTALL_DIR}/${BINARY_NAME}"
    echo ""
    echo "Run 'propel version' to verify the installation"
}

main
