# Installation

GoMetrum publishes Linux AMD64 and ARM64 binaries and Debian packages in
[GitHub Releases](https://github.com/Miklakapi/gometrum/releases). Binary
releases start with the first new version after release automation was added;
earlier tags remain source-only. Go is not required on the target device.

## Recommended: one-command APT installation

Debian, Ubuntu and Raspberry Pi OS 64-bit:

```bash
curl -fsSL https://miklakapi.github.io/gometrum/install.sh | sudo bash
```

The [installer](../install.sh) adds the signing key and APT repository, installs
GoMetrum and removes its temporary files, including on failure. It preserves
existing configuration. Running it again installs any available update; regular
updates also work through APT. For an existing manual installation, follow
[migration](#migrating-a-manual-installation) first.

The installer does not configure MQTT or start an unconfigured service. Continue
with [first-time configuration](#first-time-configuration) below.

## Manual APT setup

Add the signed repository once:

```bash
sudo install -d -m 0755 /etc/apt/keyrings
curl -fsSL https://miklakapi.github.io/gometrum/gometrum.asc -o /tmp/gometrum.asc
sudo install -m 0644 /tmp/gometrum.asc /etc/apt/keyrings/gometrum.asc
rm /tmp/gometrum.asc
echo 'deb [arch=amd64,arm64 signed-by=/etc/apt/keyrings/gometrum.asc] https://miklakapi.github.io/gometrum stable main' | sudo tee /etc/apt/sources.list.d/gometrum.list
sudo apt update
sudo apt install gometrum
```

The APT repository becomes available after its first successful release
and Pages deployment. Packages target Debian 12+ and Ubuntu 22.04+ on AMD64
or ARM64. NVIDIA sensors also require the host's NVIDIA driver/NVML library.

## First-time configuration

On first installation, edit the configuration:

```bash
sudoedit /etc/gometrum.yaml
sudo gometrum --config /etc/gometrum.yaml --validate
sudo systemctl enable --now gometrum
```

Set the MQTT broker, optional credentials, a unique MQTT `client_id` and
`agent.device_id`. The initial configuration explicitly enables CPU usage,
memory usage and uptime. See [configuration](configuration.md) for other
sensors, NVIDIA GPU metrics and buttons. The service runs as root, matching
host metrics and command execution; enable only buttons you intend to expose.
It is not enabled or started automatically on first installation.

### Updates

```bash
sudo apt update
sudo apt upgrade
```

APT installs newer published versions and the package restarts an already
running service when the configuration validates. Your `/etc/gometrum.yaml`
is a Debian configuration file and is preserved during updates. Installing a
package does not configure background automatic updates; that is a separate
system policy.

## Install a downloaded .deb

Download the `.deb` matching `dpkg --print-architecture` from a specific
[release](https://github.com/Miklakapi/gometrum/releases), then run, for example:

```bash
sudo apt install ./gometrum_1.5.0_arm64.deb
```

Use the actual version and architecture from the release. Without adding the
APT repository, future updates require downloading and installing a new `.deb`.

## Standalone binary

For manual installation, choose AMD64 (`x86_64`) or ARM64 (`aarch64`):

```bash
curl -fLO https://github.com/Miklakapi/gometrum/releases/latest/download/gometrum-linux-amd64
curl -fLO https://github.com/Miklakapi/gometrum/releases/latest/download/SHA256SUMS
sha256sum --check --ignore-missing SHA256SUMS
```

Continue only after the downloaded binary reports `OK`:

```bash
sudo install -m 0755 gometrum-linux-amd64 /usr/local/bin/gometrum
rm gometrum-linux-amd64 SHA256SUMS
```

Use `gometrum-linux-arm64` instead on Raspberry Pi OS 64-bit. Specific versions
use `/releases/download/vX.Y.Z/` instead of `/releases/latest/download/`.
Standalone binaries use the system C library; they are not fully static.
See [running](running.md) for manual configuration and systemd setup.
For a running manual installation, stop the service before replacing its
binary and start it afterwards. Manual downloads do not register an APT package.

## Migrating a manual installation

The package uses `/usr/bin/gometrum` and provides its systemd unit under
`/usr/lib/systemd/system`. An older file under `/usr/local/bin` or a unit in
`/etc/systemd/system` could override these, so package installation refuses
this situation rather than changing your files automatically.

For the standard manual installation described in previous documentation:

```bash
sudo systemctl stop gometrum
sudo systemctl disable gometrum
sudo mv /etc/systemd/system/gometrum.service /etc/systemd/system/gometrum.service.manual-backup
sudo mv /usr/local/bin/gometrum /usr/local/bin/gometrum.manual-backup
sudo systemctl daemon-reload
sudo apt install gometrum
```

Keep the existing `/etc/gometrum.yaml` if dpkg asks about a configuration file.
Validate it and enable the service using the first-install commands above.
Only after checking the service and Home Assistant metrics, delete the two
`.manual-backup` files. If your service or executable uses different paths,
adapt this migration to the existing unit's `ExecStart`.

## Building from source

Developers may still build with Go and a C compiler (CGO is required):

```bash
go install github.com/Miklakapi/gometrum/cmd/gometrum@latest
```

This installs to `$(go env GOPATH)/bin`, or `GOBIN` when configured. It does
not install the Debian package, configure systemd or set up APT updates.
