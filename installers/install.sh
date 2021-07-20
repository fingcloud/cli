#!/bin/sh
# Copyright 2019 the Deno authors. All rights reserved. MIT license.
# TODO(everyone): Keep this script simple and easily auditable.

set -e


if [ "$OS" = "Windows_NT" ]; then
  echo "Error: Windows platform detected"
  echo "check out documentation for installation on Windows" 1>&2
  echo 'https://docs.fing.ir/cli/installation' 1>&2
  exit 1
fi

case $(uname -s) in 
  "Darwin") os="darwin"  ;;
  "Linux") os="linux" ;;
esac

case $(uname -m) in 
"x86_64") arch="amd64" ;;
"arm64") arch="arm64" ;;
*) arch="amd64" ;;
esac

echo "==> Downloading for $os $arch"

ext="tar.gz"
if [ $# -eq 0 ]; then
  fing_uri="https://github.com/fingcloud/cli/releases/latest/download/fing-${os}-${arch}.${ext}"
else
  fing_uri="https://github.com/fingcloud/cli/releases/download/${1}/fing-${os}-${arch}.${ext}"
fi

fing_install="${FING_INSTALL:-$HOME/.fing}"
bin_dir="$fing_install/bin"
exe="$bin_dir/fing"

if [ ! -d "$bin_dir" ]; then
  mkdir -p "$bin_dir"
fi

curl --fail --location --progress-bar --output "$exe.$ext" "$fing_uri"
cd "$bin_dir"
tar xzf "$exe.$ext"
chmod +x "$exe"
rm "$exe.$ext"

echo "Finc CLI was installed successfully to $exe"
if command -v fing >/dev/null; then
  echo "Run 'fing --help' to get started"
else
  case $SHELL in
  /bin/zsh) shell_profile=".zshrc" ;;
  *) shell_profile=".bash_profile" ;;
  esac
  echo "Manually add the directory to your \$HOME/$shell_profile (or similar)"
  echo "  export FING_INSTALL=\"$fing_install\""
  echo "  export PATH=\"\$FING_INSTALL/bin:\$PATH\""
  echo "Run '$exe --help' to get started"
fi