# Cutover Runbook

## Preconditions

- All unit/integration/contract tests green
- SLOs satisfied in current rollout stage
- Rollback drill executed in the last 7 days

## Steps

1. Promote traffic: 0 -> 10 -> 50 -> 100
2. Monitor gates and alerts after each stage
3. Hold each stage for the required observation window
4. On breach, execute rollback runbook

## Post-Cutover

- Mark endpoint as fully strangled from legacy
- Remove temporary proxy fallback for that endpoint after 14 days stability
