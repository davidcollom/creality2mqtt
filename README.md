# Creality K1x (SE, MAX, etc.) WebSocket → MQTT Bridge

A small Go service that connects to a Creality printer's LAN WebSocket interface, decodes live telemetry, normalises the data, and publishes structured MQTT topics suitable for Home Assistant automations.

Highlights:

- Lightweight runtime: reconnecting WS client, minimal shared state
- Domain mappers: temps, job, device state, and CFS box sensors
- HA Discovery ready: stable MQTT schema, retained configs, LWT
- Strong tests: testify assert/require, 70–80% coverage target
- Tooling: Cobra CLI, pre-commit, Staticcheck, Goreleaser, Docker

---

## Repository Structure

```text
.
├── cmd/
│   └── creality2mqtt/          # CLI commands (main entry + subcommands)
├── internal/
│   ├── mapper/
│   │   ├── mapper.go           # generic key → topic mapping + filters
│   │   ├── temps.go            # domain: temperature topics
│   │   ├── job.go              # domain: print-job topics
│   │   ├── state.go            # domain: connectivity/device-state topics
│   │   └── box.go              # domain: CFS box humidity/temperature/state
│   ├── discovery/              # Home Assistant MQTT Discovery payloads
│   │   ├── discovery.go        # aggregate discovery builders
│   │   ├── sensors.go          # sensors (temp/status/fan/progress)
│   │   ├── binary_sensors.go   # binary sensors (printing/part fan)
│   │   ├── switches.go         # switch (light)
│   │   ├── camera.go           # camera stream hints
│   │   └── cfs.go              # dynamic CFS sensor discovery
│   ├── mqttclient/             # MQTT wrapper (rate limiting, helpers)
│   │   └── client.go
│   ├── wsclient/               # reconnecting WebSocket client
│   │   └── client.go
├── internal/types/
│   └── types.go               # shared type: MqttMessage
├── Dockerfile                  # container build
├── docker-compose.yml          # local broker/dev setup
├── .goreleaser.yaml            # release pipelines (GHCR images)
├── .pre-commit-config.yaml     # lint/format hooks
├── Makefile                    # common tasks (lint, test, build)
├── CONTRIBUTING.md             # contributing and local dev guide
├── CODE_OF_CONDUCT.md
├── README.md
├── go.mod / go.sum
└── coverage.out / cover.html   # generated test coverage artifacts
```

---

## Requirements

- Go **1.23+**
- MQTT broker (Mosquitto, EMQX, HiveMQ, etc.)
- Creality printer on LAN (K1, K1 Max, K1 SE) with WebSocket enabled
  - WebSocket URL format: `ws://<printer-ip>:9999/`

---

## How It Works

### 1. WebSocket Consumer

`wsclient` connects to the printer's LAN WebSocket and streams JSON messages.
These messages contain partial printer state (temperatures, progress, position, etc.).
No message types exist — instead the printer emits **complete + delta snapshots**.

### 2. Generic Mapper

`mapper.go` maps top-level keys directly to MQTT topics using:

- camelCase → snake_case normalisation
- implicit type coercion (`"219.900000"` → `"219.900"`)
- filtering out irrelevant video/streaming keys

Example:

| Key             | Topic                           |
| --------------- | ------------------------------- |
| `nozzleTemp`    | `3dprinter/k1se/nozzle_temp`    |
| `printProgress` | `3dprinter/k1se/print_progress` |
| `deviceState`   | `3dprinter/k1se/device_state`   |

### 3. Domain Mappers

To keep the project maintainable, domain-specific logic lives in separate files:

| Domain             | File       | Outputs                                                       |
| ------------------ | ---------- | ------------------------------------------------------------- |
| **Temperature**    | `temps.go` | `temperature/nozzle/current`, `temperature/bed0/target`, etc. |
| **Job / Print**    | `job.go`   | `printing`, `job/progress`, `job/file_name`, etc.             |
| **Device / State** | `state.go` | `online`, `tf_card_present`                                   |

These derived topics make Home Assistant automations much simpler.

### 4. MQTT Publishing

Messages are published through a small wrapper around the Paho MQTT client using:

All messages are QoS 0 by default. The MQTT client supports per-topic rate limiting to reduce noise during prints.

CLI/env to set minimum publish interval:

```shell
--mqtt-min-interval 1s
CREALITY_MQTT_MIN_INTERVAL=1s
```

---

## Installing & Running

### Clone the Repository

```bash
git clone https://github.com/davidcollom/creality2mqtt.git
cd creality2mqtt
```

### Build

```bash
go build ./cmd/creality2mqtt
```

### Run

```bash
./creality2mqtt \
  --ws-url ws://192.168.1.50:9999/ \
  --mqtt-broker tcp://192.168.1.10:1883 \
  --mqtt-base-topic 3dprinter/k1se \
  --mqtt-min-interval 1s
```

### Example Output

MQTT topics published:

```text
3dprinter/k1se/nozzle_temp               219.900
3dprinter/k1se/temperature/nozzle/current 219.900
3dprinter/k1se/job/progress               58
3dprinter/k1se/printing                   true
3dprinter/k1se/online                     true
```

---

## Running Tests

### Run all tests

```bash
go test ./...
```

### Test coverage

```bash
go test ./... -cover
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o cover.html
```

### Mapper & Discovery Tests

Tests live alongside the mapper package and validate:

- Key normalisation
- Numeric coercion
- Noise filtering
- Derived domain topics
- Online/printing state interpretation
- Discovery payload schemas (sensors, binary sensors, switches, camera, CFS)

All tests use `github.com/stretchr/testify/assert` and `github.com/stretchr/testify/require`. Mocks use `github.com/stretchr/testify/mock`. Coverage target is 70–80%.

Use these tests as a reference when adding support for new message fields or discovery entities.

---

## Tooling & VS Code Setup

Recommended extensions:

- Go (golang.go)
- EditorConfig
- YAML

Enable:

- `gopls`
- `staticcheck`
- `go test` on save (optional)

Recommended workspace settings:

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "staticcheck",
  "go.testFlags": ["-count=1"],
  "editor.formatOnSave": true
}
```

---

## Contributing

Please read `CONTRIBUTING.md` for detailed guidance. In brief:

- Keep domain logic isolated (`temps.go`, `job.go`, `state.go`, `box.go`)
- Add tests and maintain 70–80% coverage for all changes
- Keep MQTT schemas stable once published; treat discovery topics as contracts
- Run `go vet`, `staticcheck`, and `go test ./...` before opening PRs

This structure keeps the project clean, maintainable, and portable — whether you're consuming the data in Home Assistant, exporting to Prometheus, or building automations like room heating control during print jobs.
