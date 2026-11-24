#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."
: "${APT_SIGNING_KEY:?Add the APT_SIGNING_KEY repository secret before releasing.}"

export GNUPGHOME="$(mktemp -d)"
chmod 700 "$GNUPGHOME"
trap 'gpgconf --kill gpg-agent || true; rm -rf "$GNUPGHOME"' EXIT
printf '%s\n' "$APT_SIGNING_KEY" | gpg --batch --import
unset APT_SIGNING_KEY
key="$(gpg --batch --with-colons --list-secret-keys | awk -F: '$1 == "fpr" { print $10; exit }')"
test -n "$key"
mkdir -p _site/pool/main/g/gometrum
cp dist/gometrum_*.deb _site/pool/main/g/gometrum/
gpg --batch --armor --export-options export-minimal --export "$key" > _site/gometrum.asc
cd _site
for arch in amd64 arm64; do
  directory="dists/stable/main/binary-$arch"
  mkdir -p "$directory"
  dpkg-scanpackages --arch "$arch" pool /dev/null > "$directory/Packages"
  gzip -n -k -f "$directory/Packages"
done
apt-ftparchive \
  -o APT::FTPArchive::Release::Origin=GoMetrum \
  -o APT::FTPArchive::Release::Label=GoMetrum \
  -o APT::FTPArchive::Release::Suite=stable \
  -o APT::FTPArchive::Release::Codename=stable \
  -o 'APT::FTPArchive::Release::Architectures=amd64 arm64' \
  -o APT::FTPArchive::Release::Components=main \
  release dists/stable > dists/stable/Release
gpg --batch --yes --local-user "$key" --digest-algo SHA256 \
  --clearsign --output dists/stable/InRelease dists/stable/Release
gpg --batch --yes --local-user "$key" --digest-algo SHA256 --armor \
  --detach-sign --output dists/stable/Release.gpg dists/stable/Release
printf 'GoMetrum APT repository. See https://github.com/Miklakapi/gometrum\n' > index.txt
