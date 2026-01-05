# Contributing to Toutā Inertia

Thank you for your interest in contributing to Toutā Inertia!

## Development Setup

1. Clone the repository:
```bash
git clone git@github.com:toutaio/toutago-inertia.git
cd toutago-inertia
```

2. Install dependencies:
```bash
go mod download
```

3. Run tests:
```bash
go test ./...
```

4. Run linter:
```bash
golangci-lint run
```

## Testing

We follow Test-Driven Development (TDD):

1. Write tests first
2. Run tests (they should fail)
3. Implement the feature
4. Run tests (they should pass)
5. Refactor if needed

Aim for >85% test coverage.

## Code Style

- Follow standard Go conventions
- Run `gofmt` before committing
- Keep functions small and focused
- Add comments for exported functions
- Use descriptive variable names

## Commit Messages

Keep commit messages brief and descriptive:

```
Add Inertia protocol implementation
Fix SSR rendering edge case
Update TypeScript codegen for nested structs
```

## Pull Requests

1. Fork the repository
2. Create a feature branch
3. Write tests for your changes
4. Implement your changes
5. Ensure all tests pass
6. Update CHANGELOG.md
7. Submit a pull request

## Questions?

Open an issue or start a discussion on GitHub.
