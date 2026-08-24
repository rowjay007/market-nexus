# Phase 5 SLO and Error Budget Policy

## SLO Targets
- Search latency P99 < 50ms
- Checkout path latency P99 < 200ms
- Recommendation serving P99 < 10ms
- API success rate >= 99.9%

## Error Budget Policy
- Monthly budget: 0.1% failed requests
- Freeze progressive rollout when 30% of budget is burned in 7 days
- Immediate rollback when 5xx exceeds rollout gate thresholds

## Promotion Gates
- 0% -> 10% only when contract tests pass and no sev-1 incidents in 24h
- 10% -> 50% only when SLOs are met for 48h
- 50% -> 100% only when SLOs are met for 7 days and rollback drills pass
