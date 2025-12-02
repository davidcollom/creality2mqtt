# Contributing

Thanks for your interest in contributing to creality2mqtt!

Please follow these steps when contributing code:

- Fork the repository and create a feature branch.
- Run `gofmt` on all changed files: `gofmt -s -w .`
- Run `go vet ./...` and `go test ./...` before opening a PR.
- Commit messages should be clear and reference the issue when relevant.

Git hooks

This repo includes a `githooks/pre-commit` script. To enable it locally run:

```bash
git config core.hooksPath .githooks
```

That will cause the pre-commit checks to run automatically on commit.

If you prefer, install `pre-commit` (Python tool) and run `pre-commit install`.
