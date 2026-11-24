#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

tag="${1:?Usage: build-release.sh vX.Y.Z}"
[[ "$tag" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || { echo 'Use a stable vX.Y.Z tag.' >&2; exit 1; }
test "$(git rev-parse "refs/tags/$tag^{commit}")" = "$(git rev-parse HEAD)"
version="${tag#v}"
export CGO_ENABLED=1 GOOS=linux GOTOOLCHAIN=local
export SOURCE_DATE_EPOCH="$(git show -s --format=%ct HEAD)"
rm -rf dist
mkdir -p dist
go test -mod=readonly ./...

# Include redistribution notices alongside the binaries and inside each .deb.
cp LICENSE dist/LICENSE
{
  printf 'Go standard library\n\n'
  cat "$(go env GOROOT)/LICENSE"
  for arch in amd64 arm64; do
    compiler=gcc
    [ "$arch" != arm64 ] || compiler=aarch64-linux-gnu-gcc
    GOARCH="$arch" CC="$compiler" go list -mod=readonly -deps \
      -f '{{with .Module}}{{if not .Main}}{{.Path}}|{{.Version}}|{{.Dir}}{{end}}{{end}}' ./cmd/gometrum
  done | sort -u |
    while IFS='|' read -r module module_version directory; do
      [ -n "$module" ] || continue
      printf '\n\n===== %s %s =====\n' "$module" "$module_version"
      find "$directory" -maxdepth 1 -type f \
        \( -iname 'LICENSE*' -o -iname 'COPYING*' -o -iname 'NOTICE*' \) \
        -exec sh -c 'for file do printf "\n--- %s ---\n" "${file##*/}"; cat "$file"; done' sh {} +
    done
} > dist/THIRD_PARTY_NOTICES.txt

package_dir="$(mktemp -d)"
trap 'rm -rf "$package_dir"' EXIT
for arch in amd64 arm64; do
  compiler=gcc
  [ "$arch" != arm64 ] || compiler=aarch64-linux-gnu-gcc
  binary="dist/gometrum-linux-$arch"
  GOARCH="$arch" CC="$compiler" go build -mod=readonly -trimpath -buildvcs=true \
    -ldflags="-X github.com/Miklakapi/gometrum/internal/version.releaseVersion=$tag" \
    -o "$binary" ./cmd/gometrum
  if [ "$arch" = amd64 ]; then
    test "$("./$binary" --version)" = "$tag"
  fi
  root="$package_dir/$arch"
  mkdir -p "$root/DEBIAN" "$root/etc" "$root/usr/bin" "$root/usr/lib/systemd/system" \
    "$root/usr/share/doc/gometrum/examples"
  install -m 0755 "$binary" "$root/usr/bin/gometrum"
  install -m 0600 packaging/debian/gometrum.yaml "$root/etc/gometrum.yaml"
  printf '/etc/gometrum.yaml\n' > "$root/DEBIAN/conffiles"
  sed 's|/usr/local/bin/gometrum|/usr/bin/gometrum|' internal/service/example.service \
    > "$root/usr/lib/systemd/system/gometrum.service"
  install -m 0644 internal/config/example.yaml "$root/usr/share/doc/gometrum/examples/gometrum.yaml"
  cat LICENSE dist/THIRD_PARTY_NOTICES.txt > "$root/usr/share/doc/gometrum/copyright"
  install -m 0755 packaging/debian/preinst packaging/debian/postinst \
    packaging/debian/prerm packaging/debian/postrm "$root/DEBIAN/"
  libc_version="$(readelf --version-info "$binary" | grep -oE 'GLIBC_[0-9.]+' | sed 's/GLIBC_//' | sort -Vu | tail -n 1)"
  cat > "$root/DEBIAN/control" <<EOF
Package: gometrum
Version: $version
Architecture: $arch
Maintainer: Kacper Miklaszewski <Miklakapi@users.noreply.github.com>
Section: utils
Priority: optional
Depends: libc6 (>= $libc_version), init-system-helpers (>= 1.54)
Homepage: https://github.com/Miklakapi/gometrum
Description: System metrics for Home Assistant over MQTT
 Requires an MQTT broker and an explicitly configured YAML file.
EOF
  dpkg-deb --root-owner-group --build "$root" "dist/gometrum_${version}_${arch}.deb"
done
cd dist
sha256sum gometrum-linux-* gometrum_*.deb LICENSE THIRD_PARTY_NOTICES.txt > SHA256SUMS
