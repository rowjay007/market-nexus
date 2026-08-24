# Rollback Runbook

## Triggers
- Contract test failure against legacy API baseline
- 5xx rate breaches rollout gates for > 5 minutes
- Checkout P99 > 200ms for > 10 minutes

## Steps
1. Set feature flags for impacted endpoint to previous safe percentage.
2. Reload Nginx facade configuration.
3. Verify legacy route traffic recovery and error stabilization.
4. Run contract tests and smoke tests.
5. Create incident report with root-cause timeline.

## Verification
- Error rate back within SLO budget
- Contract tests pass
- No cross-vendor data leakage signals
