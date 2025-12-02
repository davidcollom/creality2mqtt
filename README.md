# Creality K1x (SE,MAX,etc...) WebSocket → MQTT Bridge

A small Go service that connects to a Creality printer's LAN WebSocket interface, decodes live telemetry, normalises the data, and publishes structured MQTT topics suitable for Home Assistant automations (e.g., controlling heating during a print job).

This project is intentionally designed to be:

* **Lightweight** – minimal state, no blocking loops
* **Domain-structured** – separate mappers for temperature, job state, and device state
* **Easily testable** – all mapping logic is pure functions
* **Ready for HA** – MQTT topics are stable and friendly to HA MQTT Discovery
* **VS Code friendly** – go mod tidy, gopls, linting and tests work out of the box

---

## Repository Structure

```
.
├── cmd/
│   └── k1se-bridge/         # main.go entrypoint
├── internal/
│   ├── mapper/
│   │   ├── mapper.go        # generic key → topic mapping
│   │   ├── temps.go         # domain: temperature topics
│   │   ├── job.go           # domain: print-job topics
│   │   └── state.go         # domain: connectivity/device-state topics
│   ├── mqttclient/
│   │   └── mqtt.go          # lightweight wrapper around Eclipse Paho
│   └── wsclient/
│       └── wsclient.go      # reconnecting WebSocket client
├── internal/types/
│   └── types.go             # shared type: MqttMessage
├── tests (optional)
├── go.mod
└── INSTRUCTIONS.md
```

---

## Requirements

* Go **1.23+**
* MQTT broker (Mosquitto, EMQX, HiveMQ, etc.)
* Creality printer on LAN (K1, K1 Max, K1 SE) with WebSocket enabled
  Typically:

  ```shell
  ws://<printer-ip>:9999/
  ```

---

## How It Works

### 1. WebSocket Consumer

`wsclient` connects to the printer's LAN WebSocket and streams JSON messages.
These messages contain partial printer state (temperatures, progress, position, etc.).
No message types exist — instead the printer emits **complete + delta snapshots**.

### 2. Generic Mapper

`mapper.go` maps top-level keys directly to MQTT topics using:

* camelCase → snake_case normalisation
* implicit type coercion (`"219.900000"` → `"219.900"`)
* filtering out irrelevant video/streaming keys

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

```shell
mqtt://host:port
```

All messages are QoS 0 by default (configurable later).

---

## Installing & Running

### Clone the Repository

```bash
git clone https://github.com/<your-user>/k1se-bridge.git
cd k1se-bridge
```

### Build

```bash
go build ./cmd/k1se-bridge
```

### Run

```bash
./k1se-bridge \
  --ws-url ws://192.168.1.50:9999/ \
  --mqtt-broker tcp://192.168.1.10:1883 \
  --mqtt-base-topic 3dprinter/k1se
```

### Example Output

MQTT topics published:

```plain
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
```

### Mapper Tests

Tests live alongside the mapper package and validate:

* Key normalisation
* Numeric coercion
* Noise filtering
* Derived domain topics
* Online/printing state interpretation

Use these tests as a reference when adding support for new message fields.

---

## VS Code Setup

* Install the **Go** extension (golang.go)
* Enable:

  * `gopls`
  * `staticcheck`
  * `go test` on save (optional)
* Run `go mod tidy` after adding dependencies

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

* Keep domain logic isolated (`temps.go`, `job.go`, `state.go`)
* Add tests for every new mapped field
* Keep mappings stable once published (MQTT schema stability matters)
* Use `go vet` and `staticcheck` before PRs

This structure keeps the project clean, maintainable, and portable — whether you're consuming the data in Home Assistant, exporting it to Prometheus, or building higher-level automations like room heating control during print jobs.
