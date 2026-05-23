# Contributing to GoW

Thank you for your interest in contributing to GoW — the Laravel-inspired web framework for Go!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/gow.git`
3. Install dependencies: `go mod download`
4. Run tests: `make test` (or `go test ./... -race`)
5. Run linter: `make lint` (requires golangci-lint)

## Development Guidelines

- Follow existing code style and package structure
- Add tests for new features (target >70% coverage on core packages)
- Run `make fmt && make vet` before submitting
- Update documentation in `docs/` when changing public APIs
- Use conventional commit messages

## Pull Request Process

1. Create a feature branch from `main`
2. Make focused, atomic commits
3. Ensure all tests and lint pass (`make all`)
4. Update `CHANGELOG.md` for user-facing changes
5. Open a PR with clear description and linked issues

## Reporting Issues

Use GitHub Issues with the appropriate template.

## Security

Report security vulnerabilities privately via email or GitHub Security tab. Do **not** open public issues for security problems.

## Code of Conduct

Be respectful and inclusive. See `CODE_OF_CONDUCT.md` (to be added).

We appreciate every contribution — from bug reports to major features!
