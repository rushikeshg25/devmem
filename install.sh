#!/bin/sh
# devmem installer — downloads the latest release binary and puts it on your PATH.
#
#   curl -fsSL https://raw.githubusercontent.com/rushikeshg25/devmem/master/install.sh | sh
#
# Environment overrides:
#   DEVMEM_VERSION   version to install (default: latest release)
#   DEVMEM_BIN_DIR   install directory (default: $HOME/.local/bin)
set -eu

REPO="rushikeshg25/devmem"
BIN="devmem"
BIN_DIR="${DEVMEM_BIN_DIR:-$HOME/.local/bin}"

err() { echo "install: $*" >&2; exit 1; }

# --- detect platform ---------------------------------------------------------
os="$(uname -s)"
arch="$(uname -m)"
case "$os" in
  Linux)  os="linux" ;;
  Darwin) os="darwin" ;;
  *) err "unsupported OS: $os (use the Windows zip from the releases page)" ;;
esac
case "$arch" in
  x86_64|amd64) arch="amd64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) err "unsupported architecture: $arch" ;;
esac

# --- resolve version ---------------------------------------------------------
version="${DEVMEM_VERSION:-}"
if [ -z "$version" ]; then
  version="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name":' | head -n1 | cut -d'"' -f4)"
  [ -n "$version" ] || err "could not determine latest version; set DEVMEM_VERSION"
fi
num="${version#v}"

# --- download & unpack -------------------------------------------------------
asset="${BIN}_${num}_${os}_${arch}.tar.gz"
url="https://github.com/$REPO/releases/download/$version/$asset"
tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

echo "Downloading $asset ..."
curl -fsSL "$url" -o "$tmp/$asset" || err "download failed: $url"
tar -xzf "$tmp/$asset" -C "$tmp" || err "failed to extract $asset"

# --- install -----------------------------------------------------------------
mkdir -p "$BIN_DIR"
install -m 0755 "$tmp/$BIN" "$BIN_DIR/$BIN" 2>/dev/null \
  || { cp "$tmp/$BIN" "$BIN_DIR/$BIN" && chmod 0755 "$BIN_DIR/$BIN"; }
echo "Installed $BIN $version to $BIN_DIR/$BIN"

# --- PATH hint ---------------------------------------------------------------
case ":$PATH:" in
  *":$BIN_DIR:"*) ;;
  *)
    echo
    echo "$BIN_DIR is not on your PATH. Add this to your shell profile:"
    echo "  export PATH=\"$BIN_DIR:\$PATH\""
    ;;
esac
