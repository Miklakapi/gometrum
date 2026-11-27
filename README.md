<p align="center">
  <img src="docs/images/gometrum-logo.svg" alt="GoMetrum — Know your system." width="360">
</p>

![license](https://img.shields.io/badge/license-MIT-blue)
![linux](https://img.shields.io/badge/os-Linux-green)
![language](https://img.shields.io/badge/language-Go_1.25.1-blue)
![version](https://img.shields.io/github/v/tag/Miklakapi/gometrum)

A lightweight system metrics agent written in Go.

**GoMetrum collects host-level metrics and exposes optional host actions to Home Assistant using MQTT Discovery.**

The project focuses on **explicit configuration**, **deterministic behavior**, and **low runtime overhead**, while avoiding implicit defaults and hidden logic.

<p align="center">
    <img src="docs/images/sensors.png">
</p>

## Installation

Debian, Ubuntu and Raspberry Pi OS 64-bit:

```bash
curl -fsSL https://miklakapi.github.io/gometrum/install.sh | sudo bash
```

This adds the signed APT repository and installs the package. Go is not required.
See [installation and updates](docs/installation.md) for first-time configuration,
manual setup and migration from a manual installation.

[Download releases](https://github.com/Miklakapi/gometrum/releases/latest)

## Documentation

**Project documentation:**

- [Installation](docs/installation.md)
- [Configuration](docs/configuration.md)
- [Running GoMetrum](docs/running.md)
- [Command-line interface](docs/cli.md)

## Overview

GoMetrum is a background agent designed to run on Linux hosts and report system-level metrics to Home Assistant.

Key characteristics:

- single explicit YAML configuration file,
- sensors and buttons enabled **only by presence**,
- predictable Home Assistant entity set,
- clean startup and shutdown behavior.
- publishes states via MQTT,

## Explicit configuration model

Each entity (sensor or button) exists **only if it is explicitly defined** in the configuration file.

There are:

- no global enable/disable switches,
- no implicit or auto-generated sensors,
- no “magic defaults”.

As a result:

- the number of Home Assistant entities is always predictable,
- the configuration file fully describes the resulting system state,
- behavior is deterministic and easy to reason about.

## Features

- System metrics collection (CPU, memory, disk, network, GPU)
- Explicit per-sensor configuration
- Per-sensor refresh intervals
- MQTT retained state publishing
- Home Assistant MQTT Discovery integration
- Home Assistant button entities for executing host commands
- Availability reporting (`online` / `offline`)
- Discovery cleanup (`--purge` mode)
- Structured logging sinks (UDP and HTTP with multiple codecs and batching)
- Low runtime overhead
- Bash and Zsh tab completion, installed automatically by the Debian package

## Technologies

Project is created with:

- Go 1.25.1
- MQTT (Eclipse Paho client)
- Home Assistant MQTT Discovery
- gopsutil
- systemd (optional)
