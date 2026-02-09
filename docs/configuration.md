# Configuration

GoMetrum uses a **single explicit YAML configuration file**.

The configuration file fully defines:

- which sensors exist,
- how often they are collected,
- how data is published to Home Assistant.

No sensors are created unless explicitly defined.

## Configuration model

GoMetrum follows an **explicit configuration model**:

- sensors are enabled by presence,
- there are no global enable/disable switches,
- startup behavior is deterministic.

Commenting out a sensor block disables it.
Removing a sensor from the file removes the corresponding Home Assistant entity.

## Generate example configuration

Print an example configuration to stdout:

```bash
gometrum --generate-config
```

Write directly to a file:

```bash
gometrum --generate-config --config gometrum.yaml
```

The generated file contains all available configuration options
with inline comments.

## Configuration structure

At a high level, the configuration file consists of four sections:

- `mqtt` - MQTT connection and discovery settings
- `agent` - device metadata used by Home Assistant
- `sensors` - sensor definitions and collection intervals

## MQTT section

The `mqtt` section defines how GoMetrum connects to the MQTT broker
and publishes data.

Key fields:

- `host`, `port` - broker address
- `username`, `password` - optional authentication
- `client_id` - must be unique per running agent
- `discovery_prefix` - Home Assistant MQTT discovery prefix
- `state_prefix` - base topic for sensor state publishing
- `default_interval` - fallback interval for sensors without an explicit interval

## Agent section

The agent section defines device metadata as seen by Home Assistant.

Important fields:

- `device_id` - required, must be unique per host
- `device_name` - human-readable name
- `manufacturer`, `model` - Home Assistant metadata

## Sensors section

The sensors section defines which system metrics are collected.

Each sensor:

- exists only if present in the configuration,
- may define its own `interval`,
- may provide Home Assistant overrides under the `ha` key.

Some sensors support additional options (e.g. `include_mounts` for disk usage).

Multiple sensors can share the same refresh interval internally.

## Validate configuration

You can validate the configuration at any time:

```bash
gometrum --config gometrum.yaml --validate
```

Validation checks:

- YAML structure,
- required fields,
- sensor definitions,
- interval formats.

## File location

Common locations:

- local testing: ./gometrum.yaml

- production: /etc/gometrum.yaml
