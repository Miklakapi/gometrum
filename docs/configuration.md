# Configuration

GoMetrum uses a **single explicit YAML configuration file**.

The configuration file fully defines:

- logging behavior,
- which sensors exist,
- how often they are collected,
- how data is published to Home Assistant.

No sensors are created unless explicitly defined.

## Configuration model

GoMetrum follows an **explicit configuration model**:

- sensors are enabled by presence,
- logging is fully defined in configuration,
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

- `log` - logging configuration
- `mqtt` - MQTT connection and discovery settings
- `agent` - device metadata used by Home Assistant
- `sensors` - sensor definitions and collection intervals

## Log section

The `log` section defines how GoMetrum logs information and where logs are sent.

### Log level

`log.level` defines the global minimum log level.

Supported values: `debug`, `info`, `warn`, `error`.

If not specified, defaults to `info`.

### Log sinks

Sinks define additional destinations for logs.

Console logging is always enabled.

Each sink has the following common fields:

- `name` - human-readable sink identifier (used for diagnostics and internal logging)
- `type` - sink type (`udp` or `http`)
- `level` - minimum log level for this sink (inherits from `log.level` if not specified)

#### UDP sink

UDP sends logs as plain text (syslog-like format).

Properties:

- `name` - descriptive identifier of the sink
- `addr` - destination in `host:port` format
- `level` - minimum level for this sink (optional, inherits from global level)

UDP logs are sent as plain text messages.<br>
No authentication or TLS is supported for UDP.

#### HTTP sink

HTTP sends logs to a REST endpoint.

Properties:

- `name` - descriptive identifier of the sink
- `url` - target endpoint (`http` or `https`)
- `method` - HTTP method (default: `POST`)
- `timeout` - request timeout
- `headers` - optional HTTP headers
- `level` - minimum level for this sink (optional, inherits from global level)
- `codec` - payload format

Supported HTTP codecs:

- `event_json` – single structured JSON event per request
- `text` – plain text log line
- `ndjson` – newline-delimited JSON (multiple events per request)
- `loki` – Grafana Loki push API format

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

The `sensors` section defines which system metrics are collected.

Each sensor:

- exists only if present in the configuration,
- is identified by its key (e.g. cpu_usage, memory_usage),
- may define its own `interval`,
- may provide Home Assistant overrides under the `ha` key.

### Sensor activation model

Sensors are enabled strictly by presence.<br>
There are no global enable/disable switches.

- Commenting out a sensor disables it.
- Removing a sensor removes the corresponding Home Assistant entity.
- Startup behavior is deterministic.

`interval`

Each sensor may define an `interval` using Go duration format (e.g. `5s`, `30s`, `1m`, `5m`).

If a sensor does not define its own interval, it inherits `mqtt.default_interval`.

Intervals must be greater than zero.

Internally, sensors sharing the same interval are grouped
for efficient scheduling.

`name`

A sensor may optionally define a human-readable `name`.

If not provided, a default name is used.

`ha` **(Home Assistant overrides)**

The optional `ha` block allows overriding Home Assistant entity metadata.

Supported override fields include:

- `icon`
- `unit`
- `device_class`
- `state_class`

If not specified, sensible defaults are applied where applicable.

### Sensor-specific options

Some sensors expose additional configuration fields.

Example:

- `include_mounts` (for disk usage sensors)

Additional options are validated per sensor type.

### Validation rules

Configuration validation ensures:

- interval values are valid durations,
- sensor keys are recognized,
- required options (if any) are provided,
- sensor-specific constraints are respected.

Hardware availability is not validated at configuration time.

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
