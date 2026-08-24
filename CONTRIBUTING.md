# Contributing to MarketNexus

Thanks for your interest in contributing.

## Development Flow

1. Sync with `main` and create a feature branch.
2. Keep each PR scoped to one primary concern or bounded context.
3. Add/adjust tests for changed behavior.
4. Run checks:
   - `go test ./...`
   - `go test ./tests/saga/...` if saga behavior changes
   - `go test ./tests/contracts/...` if API contract shapes change
5. Use Conventional Commits.
6. Open a pull request using the project PR template and include evidence.

## Local Setup

```bash
git clone <repo-url>
cd market-nexus
go test ./...
```

If tests are green, you are ready to contribute.

## Branch Naming

Use one of:

- `feat/<short-description>`
- `fix/<short-description>`
- `docs/<short-description>`
- `chore/<short-description>`

## Engineering Guidelines

- Preserve bounded context boundaries.
- Do not share databases across contexts.
- Prefer events and ACLs over direct model coupling.
- Add tests for all behavior changes and compensation paths.
- Do not introduce cross-vendor access patterns.
- Update ADR/docs when architectural behavior changes.

## Commit Convention

Examples:
- `feat(ordering): add checkout compensation branch`
- `fix(inventory): correct optimistic lock check`
- `docs: update architecture handbook`

## Pull Request Expectations

Your PR should clearly state:

1. Which bounded context(s) changed
2. Whether contracts/migrations changed
3. How backward compatibility is preserved
4. What tests were run and their results

## Pull Request Checklist

- [ ] Tests pass locally
- [ ] New behavior is covered by tests
- [ ] No cross-context model leakage introduced
- [ ] README/docs updated when required
- [ ] Contract and rollout assets updated when migration behavior changes
