#!/usr/bin/env bash
set -euo pipefail

if [ "$(id -u)" -ne 0 ]; then
  echo 'Run this installer with sudo bash.' >&2
  exit 1
fi
for command in apt-get dpkg curl; do
  command -v "$command" >/dev/null || { echo "Required command: $command" >&2; exit 1; }
done
arch="$(dpkg --print-architecture)"
case "$arch" in
  amd64|arm64) ;;
  *) echo "Unsupported architecture: $arch (supported: amd64, arm64)." >&2; exit 1 ;;
esac

repository='https://miklakapi.github.io/gometrum'
temporary="$(mktemp -d)"
trap 'rm -rf "$temporary"' EXIT

curl -fsSL "$repository/gometrum.asc" -o "$temporary/gometrum.asc"
printf 'deb [arch=%s signed-by=/etc/apt/keyrings/gometrum.asc] %s stable main\n' \
  "$arch" "$repository" > "$temporary/gometrum.list"
install -d -m 0755 /etc/apt/keyrings
install -m 0644 "$temporary/gometrum.asc" /etc/apt/keyrings/gometrum.asc
install -m 0644 "$temporary/gometrum.list" /etc/apt/sources.list.d/gometrum.list

apt-get update
apt-get -o Dpkg::Options::=--force-confdef -o Dpkg::Options::=--force-confold \
  install -y gometrum </dev/null

printf '\nGoMetrum is installed. Updates: sudo apt update && sudo apt upgrade\n'
printf 'First setup: https://github.com/Miklakapi/gometrum/blob/master/docs/installation.md\n'
