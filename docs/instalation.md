# Installation

This document describes how to install GoMetrum on a Linux system.

## Requirements

- Linux host
- Go ≥ 1.25
- MQTT broker
- Home Assistant with MQTT integration enabled

## Recommended: install using `go install`

This method:

- installs a single static binary,
- does not require root,
- is easy to update.

```bash
go install github.com/Miklakapi/gometrum/cmd/gometrum@latest
```

Make sure your Go binary directory is in `$PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

> Note: this command affects only the current shell session.
> To make the change permanent, add it to your shell configuration file
> (e.g. ~/.bashrc, ~/.zshrc).

Verify installation:

```bash
gometrum --help
```

## Optional: system-wide installation

For production systems, you may want the binary available system-wide.

```bash
sudo cp "$(go env GOPATH)/bin/gometrum" /usr/local/bin/gometrum
sudo chmod +x /usr/local/bin/gometrum
```

Verify:

```bash
which gometrum
sudo gometrum --help
```
