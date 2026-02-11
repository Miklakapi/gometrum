# Command-line interface

GoMetrum is configured and controlled primarily via command-line flags.

The CLI supports:

- normal agent operation,
- one-shot execution,
- configuration and service generation,
- validation and cleanup tasks.

## Common flags

### `--config`, `-c`

Path to the YAML configuration file.

Example:

```bash
gometrum --config /etc/gometrum.yaml
```

## Agent execution modes

### Normal mode (default)

Runs the agent continuously, collecting and publishing metrics according to the configured sensor intervals.

```bash
gometrum --config gometrum.yaml
```

---

### `--once`

Collect and publish all configured sensors **once**, then exit.

This mode:

- publishes data to MQTT normally,
- performs Home Assistant discovery if needed,
- exits immediately after one collection cycle.

Useful for:

- testing configuration,
- scripting,
- debugging sensor output.

```bash
gometrum --once --config gometrum.yaml
```

---

### `--dry-run`

Collect metrics but do not publish them to MQTT.

Instead of publishing, all MQTT messages(discovery, state, availability) are printed to stdout.

This flag can be combined with `--once`.

```bash
gometrum --dry-run --config gometrum.yaml
```

## One-shot / exit modes

The following flags perform a single action and then exit.
Only one of these flags can be used at a time.

---

### `--generate-config`

Generate an example YAML configuration and exit.

```bash
gometrum --generate-config
```

---

### `--validate`

Validate the configuration file and exit.

- YAML structure,
- required fields,
- sensor definitions,
- interval formats.

```bash
gometrum --validate --config gometrum.yaml
```

---

### `--generate-service`

Generate a systemd service file and exit.

Prints to stdout by default.

```bash
gometrum --generate-service
```

Use --service-path to write directly to a file:

```bash
gometrum --generate-service --service-path ./gometrum.service
```

---

### `--purge`

Remove all Home Assistant entities defined in the configuration.

This publishes empty retained MQTT discovery messages,
causing Home Assistant to delete the entities.

```bash
gometrum --purge --config gometrum.yaml
```

---

### `--version`, `-v`

Print version information and exit.

```bash
gometrum --version
```

---

### Flag compatibility rules

- Only one exit mode flag can be used at a time:
    - --generate-config
    - --validate
    - --purge
    - --generate-service
    - --version

- Flags `--once` and `--dry-run` can be used together.
- When combined, GoMetrum collects metrics once, prints all MQTT messages
  to stdout, and exits.

Invalid flag combinations result in an error.
