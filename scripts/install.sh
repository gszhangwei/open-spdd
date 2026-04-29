#!/usr/bin/env bash
#
# install.sh — One-shot installer for openspdd via `go install`.
#
# Usage:
#   ./scripts/install.sh                     # installs @latest
#   ./scripts/install.sh v1.2.3              # installs a specific tag
#   VERSION=v1.2.3 ./scripts/install.sh      # alternative form
#
# After a successful install, this script verifies that the resulting
# binary directory is on $PATH and prints shell-specific instructions
# if it is not.

set -euo pipefail

MODULE="github.com/gszhangwei/open-spdd"
PACKAGE="${MODULE}/cmd/openspdd"
VERSION="${1:-${VERSION:-latest}}"

bold()  { printf '\033[1m%s\033[0m\n' "$*"; }
green() { printf '\033[32m%s\033[0m\n' "$*"; }
yellow(){ printf '\033[33m%s\033[0m\n' "$*"; }
red()   { printf '\033[31m%s\033[0m\n' "$*" >&2; }

if ! command -v go >/dev/null 2>&1; then
  red "Error: 'go' is not installed or not on PATH."
  red "Install Go first: https://go.dev/dl/"
  exit 1
fi

GOBIN="$(go env GOBIN)"
if [ -z "$GOBIN" ]; then
  GOBIN="$(go env GOPATH)/bin"
fi

bold "Installing ${PACKAGE}@${VERSION} ..."
go install "${PACKAGE}@${VERSION}"
green "✓ Installed: ${GOBIN}/openspdd"
echo

case ":${PATH}:" in
  *":${GOBIN}:"*)
    green "✓ ${GOBIN} is already on your PATH. You can run:"
    echo
    echo "    openspdd --help"
    echo
    exit 0
    ;;
esac

yellow "⚠ ${GOBIN} is NOT on your PATH."
echo "  The 'openspdd' command will not be found until you add it."
echo

SHELL_NAME="$(basename "${SHELL:-}")"
case "$SHELL_NAME" in
  zsh)   RC_FILE="$HOME/.zshrc"  ; LINE="export PATH=\"${GOBIN}:\$PATH\""    ;;
  bash)  RC_FILE="$HOME/.bashrc" ; LINE="export PATH=\"${GOBIN}:\$PATH\""    ;;
  fish)  RC_FILE="$HOME/.config/fish/config.fish" ; LINE="set -gx PATH ${GOBIN} \$PATH" ;;
  *)     RC_FILE="$HOME/.profile"; LINE="export PATH=\"${GOBIN}:\$PATH\""    ;;
esac

bold "Run these commands to fix it:"
echo
echo "    echo '${LINE}' >> ${RC_FILE}"
echo "    source ${RC_FILE}"
echo
echo "Then verify with:"
echo
echo "    openspdd --help"
echo
