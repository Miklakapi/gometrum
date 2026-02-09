# Running GoMetrum

After installation and configuration, GoMetrum can be run in two modes:
manual (foreground) or as a systemd service.

## Foreground mode (manual)

Useful for:

- testing configuration,
- debugging MQTT or discovery,
- temporary runs.

Run the agent:

```bash
gometrum --config gometrum.yaml
```

Logs are written to stdout/stderr.

Stop the agent with `Ctrl+C`.

On shutdown:

- availability is updated to offline,
- the MQTT connection is closed cleanly.

## Production usage (systemd service)

Running GoMetrum as a systemd service is recommended for long-running or production setups.

### Generate service file

Generate a systemd unit file:

```bash
gometrum --generate-service
```

Or write it directly to a file:

```bash
gometrum --generate-service --service-path ./gometrum.service
```

### Install and start the service

```bash
sudo cp gometrum.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now gometrum
```

### Check status and logs

```bash
systemctl status gometrum
journalctl -u gometrum -f
```

### Service behavior

When running as a service, GoMetrum:

- starts automatically on boot,
- reconnects to MQTT if needed,
- publishes availability (online / offline),
- shuts down gracefully on stop or reboot.
