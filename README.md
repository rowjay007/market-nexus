# MarketNexus

Production-grade Phase 1 through Phase 4 scaffold for a DDD, event-driven, multi-vendor marketplace migration from a PHP monolith using the Strangler Fig pattern.

## Phase 1 to Phase 4 Scope

Implemented bounded contexts:
- Catalog BC
- Inventory BC
- Ordering BC
- Pricing BC
- Payment BC
- Fulfillment BC
- Search BC
- Review/Trust BC
- Analytics BC

Implemented architectural guarantees:
- DDD-style aggregates and domain events per BC
- Strict bounded context data/model ownership
- ACL-based integration from Ordering to Catalog read model
- Choreography-style saga step for OrderPlaced -> InventoryReserved
- Compensation test path for reservation failure
- Checkout choreography for Pricing -> Payment -> Fulfillment with compensation
- Search indexing and query ranking scoped by vendor isolation
- Review/Trust review submission, dispute signals, and vendor rating metrics
- Analytics behavior ingestion and recommendation precompute cache model
- Federation supergraph composition and gateway config artifacts
- Strangler facade routing config with progressive feature flag rollout
- Kafka topic and Protobuf event contract definitions

## Repository Layout

- `services/catalog` - Catalog BC service
- `services/inventory` - Inventory BC service
- `services/ordering` - Ordering BC service
- `services/pricing` - Pricing BC service
- `services/payment` - Payment BC service
- `services/fulfillment` - Fulfillment BC service
- `services/search` - Search BC service
- `services/reviewtrust` - Review/Trust BC service
- `services/analytics` - Analytics BC service
- `contracts/proto` - Event schema contracts (Protobuf)
- `contracts/kafka-topics.yaml` - Kafka topic ownership and schema mapping
- `infra/migrations` - Per-BC migration files
- `deploy/facade/nginx/nginx.conf` - Strangler facade route proxying
- `deploy/facade/flags/phase1.yaml` - Feature rollout policy
- `docs/adr` - Architecture decision records
- `tests/saga` - Cross-BC saga integration tests

## Key Design Notes

- No database sharing across bounded contexts.
- No cross-context model imports across public boundaries.
- Ordering consumes Catalog events via an ACL projection read model.
- Inventory reservations are isolated by `vendor_id` + `sku`.
- Compensation behavior is explicitly tested.

## Run Tests

```bash
go test ./...
```

## Current Status

This repository currently contains production-grade architecture scaffolding and domain-level implementation for Phase 1 through Phase 4.

What is intentionally left as stubs for next iterations:
- Full GraphQL resolver/server runtime wiring (schema and gqlgen config included)
- Real Kafka producer/consumer adapter and outbox pattern
- Real persistent repositories (MongoDB/PostgreSQL adapters)
- Kong/API Gateway runtime and complete federation supergraph composition
- Legacy PHP compatibility contract test harness

## ADRs

- `docs/adr/0001-ddd-bounded-context-isolation.md`
- `docs/adr/0002-strangler-facade-feature-flags.md`
- `docs/adr/0003-ordering-saga-choreography.md`
- `docs/adr/0004-phase2-checkout-saga.md`
- `docs/adr/0005-phase3-search-reviewtrust.md`
- `docs/adr/0006-phase4-analytics-federation.md`
