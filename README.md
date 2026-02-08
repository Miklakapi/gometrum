# GoMetrum

![license](https://img.shields.io/badge/license-MIT-blue)
![linux](https://img.shields.io/badge/os-Linux-green)
![language](https://img.shields.io/badge/language-Go_1.25.1-blue)
![version](https://img.shields.io/badge/version-1.0.2-success)
![status](https://img.shields.io/badge/status-development-blue)

A lightweight system metrics agent written in Go.

**GoMetrum collects host-level metrics and publishes them to Home Assistant using MQTT Discovery.**

The project focuses on **explicit configuration**, **deterministic behavior**, and **low runtime overhead**, while avoiding implicit defaults and hidden logic.

## Table of Contents

- [General info](#general-info)
- [Architecture](#architecture)
- [Technologies](#technologies)
- [Setup](#setup)
- [Features](#features)
- [Status](#status)

## General info

GoMetrum is a background agent designed to run on Linux hosts and continuously report system metrics to Home Assistant.

The agent:

- reads its configuration from a single YAML file,
- explicitly builds only the sensors defined in that configuration,
- periodically collects metrics from the system,
- publishes states via MQTT,
- registers and removes entities using Home Assistant MQTT Discovery.

### Explicit configuration model

Each sensor exists **only if it is explicitly defined in the configuration file**.

There are:

- no global enable/disable switches,
- no implicit or auto-generated sensors,
- no “magic defaults”.

As a result:

- the number of Home Assistant entities is always predictable,
- the YAML file fully describes the resulting system state,
- startup behavior is deterministic and easy to reason about.

The project was created as a simpler alternative to larger monitoring solutions, with a focus on

- low runtime overhead,
- predictable behavior,
- easy debugging and inspection.

## Architecture

GoMetrum follows a simple, modular architecture with clear separation of responsibilities.

High-level components:

- CLI layer

    Parses flags, handles one-shot modes (config generation, validation, purge).

- Configuration module

    Loads, normalizes, validates, and applies defaults to YAML configuration.

- Sensor registry

    Maps sensor identifiers to factories and default Home Assistant metadata.

- Agent
    - manages MQTT connection lifecycle,
    - publishes availability,
    - handles Home Assistant discovery,
    - schedules sensor collection by interval.

- MQTT client wrapper

    Thin abstraction over the Paho MQTT client with retained publishing and availability handling.

Sensors are grouped internally by refresh interval, allowing multiple sensors to share the same scheduler without additional goroutines.

## Technologies

Project is created with:

- Go 1.25.1
- MQTT (Eclipse Paho client)
- Home Assistant MQTT Discovery
- gopsutil – system metrics
- systemd (optional, for running as a service)

## Setup

### 1. Installation

Install using go install (recommended)

```bash
go install github.com/Miklakapi/gometrum/cmd/gometrum@latest
```

Make sure your Go binary directory is in `$PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

Verify installation:

```bash
gometrum --help
```

**Optional: system-wide installation**<br>
For production systems, you may want the binary available system-wide:

```bash
sudo cp "$(go env GOPATH)/bin/gometrum" /usr/local/bin/gometrum
sudo chmod +x /usr/local/bin/gometrum
```

Verify:

```bash
which gometrum
sudo gometrum -help
```

### 2. Configuration

GoMetrum uses a **single explicit YAML configuration file**.

Generate an example configuration

```bash
gometrum --generate-config
```

Or write it directly to a file

```bash
gometrum --generate-config --config /etc/gometrum.yaml
```

Edit the configuration and define **only the sensors you want to exist**.<br>
No sensors are created unless explicitly configured.

You can validate the configuration at any time

```bash
gometrum --config /etc/gometrum.yaml --validate
```

### Running GoMetrum

After installation and configuration, choose one of the following modes.

#### Option A: Manual (foreground) run

This mode is useful for:

- testing configuration,
- debugging MQTT or discovery,
- running GoMetrum temporarily.

Run the agent manually

```bash
gometrum --config /etc/gometrum.yaml
```

Logs are written to stdout/stderr.

Stop the agent with `Ctrl+C`.
On shutdown, availability is updated correctly (`offline`).

#### Option B: systemd service (recommended for production)

**Generate service file**<br>
Generate a systemd unit file

```bash
gometrum --generate-service
```

Or write it directly to disk

```bash
gometrum --generate-service --service-path ./gometrum.service
```

Install and enable the service

```bash
sudo cp gometrum.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now gometrum
```

Check status and logs

```bash
systemctl status gometrum
journalctl -u gometrum -f
```

The agent will now:

- start automatically on boot,
- reconnect to MQTT if needed,
- publish availability and sensor states continuously.

### Cleanup (MQTT discovery purge)

To remove all Home Assistant entities registered by GoMetrum:

```bash
gometrum --purge --config /etc/gometrum.yaml
```

This publishes empty retained discovery messages, causing Home Assistant to remove the entities.

## Features

- System metrics collection (CPU, memory, disk, network, GPU)
- Explicit sensor configuration (no implicit entities)
- One Home Assistant entity per configured sensor (or per include entry)
- Multiple sensors per mount point (disk usage)
- Per-sensor refresh intervals
- MQTT retained state publishing
- Home Assistant MQTT Discovery integration
- Availability reporting (`online` / `offline`)
- Deterministic startup and shutdown behavior
- Discovery cleanup (`--purge` mode)

## Status

The project is in active development.
