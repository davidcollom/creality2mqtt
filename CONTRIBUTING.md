# Contributing

Thanks for your interest in contributing to creality2mqtt! This guide explains how to get started, the project tooling, and the expectations for tests and code quality.

## Getting Started

- Clone: `git clone https://github.com/davidcollom/creality2mqtt.git && cd creality2mqtt`
- Install Go 1.23+ and the VS Code Go extension (golang.go)
- Install tooling: `go install honnef.co/go/tools/cmd/staticcheck@latest`
- Tidy modules after adding deps: `go mod tidy`

## Running Locally

- Start a local MQTT broker (optional) using Docker Compose:

```bash
docker compose up -d
```

- Build the CLI:

```bash
go build ./cmd/creality2mqtt
```

- Run against your printer:

```bash
./creality2mqtt \
  --ws-url ws://<printer-ip>:9999/ \
  --mqtt-broker tcp://127.0.0.1:1883 \
  --mqtt-base-topic 3dprinter/k1se \
  --mqtt-min-interval 1s
```

## Tests & Coverage (Required)

- All code MUST maintain 70–80% test coverage
- Run tests before any commit:

```bash
go test ./...
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o cover.html
```

- Use `github.com/stretchr/testify/assert` and `github.com/stretchr/testify/require`
- Use `github.com/stretchr/testify/mock` for mocks
- Include edge cases and error handling; keep tests in table format where appropriate

## Linting & Quality

- Format: `gofmt -s -w .`
- Vet: `go vet ./...`
- Static analysis: `staticcheck ./...`

## Project Structure Expectations

- Keep domain logic isolated in mapper files (`temps.go`, `job.go`, `state.go`, `box.go`)
- Discovery payloads are contracts; changes must be additive and documented
- MQTT schemas must remain stable once published
- Avoid adding unnecessary complexity; prefer small focused changes

## Git Hooks / Pre-Commit

This repo includes a `.pre-commit-config.yaml`. To enable pre-commit checks locally:

```bash
pip3 install pre-commit
pre-commit install
```

Alternatively, you can use Git hooks:

```bash
git config core.hooksPath .githooks
```

## Releases

- Releases are built via Goreleaser using multi-arch Docker images to GHCR
- Keep Dockerfile and `.goreleaser.yaml` up to date when adding binaries or assets

## Opening a Pull Request

- Rebase your branch on `main`
- Ensure tests pass and coverage meets the target
- Include a concise PR description: context, changes, tests, and any discovery schema updates

Thank you for contributing — this project aims to be clean, maintainable, and friendly to Home Assistant automations.
