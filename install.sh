#!/bin/sh
set -eu

# one-liner: curl -sSL https://raw.githubusercontent.com/hyperpuncher/yd-dl/main/install.sh | sh

os=$(uname -s)
arch=$(uname -m)

case "$os-$arch" in
	Linux-x86_64)  binary=yd-dl-linux-amd64 ;;
	Linux-aarch64) binary=yd-dl-linux-arm64 ;;
	Darwin-arm64)  binary=yd-dl-darwin-arm64 ;;
	Darwin-x86_64) binary=yd-dl-darwin-amd64 ;;
	*)             echo "unsupported: $os $arch"; exit 1 ;;
esac

url="https://github.com/hyperpuncher/yd-dl/releases/latest/download/$binary"
dest="$HOME/.local/bin/yd-dl"

mkdir -p "$(dirname "$dest")"

echo "→ $url"
curl -sSL "$url" -o "$dest"
chmod +x "$dest"
echo "→ yd-dl installed to $dest"
