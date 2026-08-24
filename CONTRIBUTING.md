# Contributing to MarketNexus

Thanks for your interest in contributing.

## Development Flow

1. Create a branch from `main`.
2. Make focused changes by bounded context.
3. Run checks:
   - `go test ./...`
4. Use Conventional Commits for commit messages.
5. Open a pull request with:
   - Summary of changes
   - Affected bounded contexts
   - Test evidence

## Engineering Guidelines

- Preserve bounded context boundaries.
- Do not share databases across contexts.
- Prefer events and ACLs over direct model coupling.
- Add tests for all behavior changes and compensation paths.

## Commit Convention

Examples:
- `feat(ordering): add checkout compensation branch`
- `fix(inventory): correct optimistic lock check`
- `docs: update architecture handbook`

## Pull Request Checklist

- [ ] Tests pass locally
- [ ] New behavior is covered by tests
- [ ] No cross-context model leakage introduced
- [ ] README/docs updated when required
