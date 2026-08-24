# MarketNexus

MarketNexus is a multi-vendor e-commerce marketplace platform built as a domain-driven, event-driven microservice system in Go.

The project demonstrates a full migration path from a legacy PHP monolith using Strangler Fig rollout controls, while preserving API compatibility and enforcing bounded-context ownership.

## Table of Contents

1. Project Vision
2. Architecture Overview
3. Bounded Contexts
4. Key Technical Patterns
5. End-to-End Flows
6. Getting Started
7. Day-to-Day Development
8. Repository Structure
9. Local Validation
10. Migration and Operations
11. Security, Reliability, and SLOs
12. Troubleshooting
13. Current Implementation Status
14. ADR Index

## Project Vision

MarketNexus models a marketplace where:
- Vendors manage listings and stock
- Buyers browse, search, order, and pay
- Fulfillment and trust are managed by dedicated domains
- Analytics computes recommendation candidates from behavior events

The system is intentionally split into independent business domains to reduce coupling and enable safe incremental delivery.

## Architecture Overview

- Language: Go 1.23+
- Style: DDD + CQRS-inspired separation + event choreography
- Integration:
	- Synchronous: GraphQL subgraph boundaries (schema and config scaffolded)
	- Asynchronous: Domain events mapped to Kafka topics
- Migration: Nginx facade with progressive feature-flag routing

### High-Level Service Map

1. Catalog
2. Inventory
3. Ordering
4. Pricing
5. Payment
6. Fulfillment
7. Search
8. Review/Trust
9. Analytics

## Bounded Contexts

### 1. Catalog
- Owns product and variant domain behavior
- Publishes listing events
- Path: `services/catalog`

### 2. Inventory
- Owns stock reservation and release
- Path: `services/inventory`

### 3. Ordering
- Owns order lifecycle and checkout choreography hooks
- Uses ACL projection to consume catalog knowledge
- Path: `services/ordering`

### 4. Pricing
- Owns quoting logic and discount/tax composition
- Path: `services/pricing`

### 5. Payment
- Owns capture and refund lifecycle
- Path: `services/payment`

### 6. Fulfillment
- Owns shipment scheduling and cancellation
- Path: `services/fulfillment`

### 7. Search
- Owns search document indexing and ranking behavior
- Path: `services/search`

### 8. Review/Trust
- Owns review, dispute, and trust-signal events
- Path: `services/reviewtrust`

### 9. Analytics
- Owns behavior ingestion and recommendation precompute model
- Path: `services/analytics`

## Key Technical Patterns

### Domain-Driven Design
- Each context has its own model, application service, and storage contract.
- Cross-context direct model coupling is avoided.

### Anti-Corruption Layer (ACL)
- Ordering consumes translated catalog read-model data via ACL projection.
- Path: `services/ordering/acl/catalogreadmodel`

### Event Choreography and Compensation
- Order and checkout flows rely on event-driven progression and rollback behavior.
- Compensation (release/refund/cancel) is tested explicitly in saga tests.

### Vendor Isolation
- Domain data operations are scoped by `vendor_id`.
- Prevents cross-vendor data leakage by design.

### Strangler Fig Migration
- Legacy fallback and progressive cutover routing:
	- `deploy/facade/nginx/nginx.conf`
	- `deploy/facade/flags/phase1.yaml`

### Contract-First Eventing
- Topic ownership map:
	- `contracts/kafka-topics.yaml`
- Event schemas:
	- `contracts/proto/*/events.proto`

## End-to-End Flows

### Order Reservation Saga (Phase 1)
1. Order created
2. Inventory reservation attempted
3. Success: order confirms
4. Failure: order cancels and reserved stock is released

Test: `tests/saga/order_reservation_saga_test.go`

### Checkout Saga (Phase 2)
1. Order reserved
2. Pricing quote computed
3. Payment captured
4. Fulfillment scheduled
5. Failure path includes refund/release compensation

Test: `tests/saga/checkout_phase2_saga_test.go`

### Search and Trust Flow (Phase 3)
1. Search documents indexed and ranked per vendor
2. Reviews submitted and trust metrics generated

Test: `tests/saga/phase3_search_review_saga_test.go`

### Recommendations Flow (Phase 4)
1. Analytics consumes behavior events
2. Scores recommendations in precompute cache
3. Serves top products by user

Test: `tests/saga/phase4_analytics_recommendations_test.go`

## Getting Started

### Prerequisites

- Go 1.23+
- Git
- A Unix-like shell (macOS/Linux)

### Clone and Validate

```bash
git clone <your-fork-or-origin-url>
cd market-nexus
go test ./...
```

If tests pass, your local environment is ready.

### What You Can Run Today

This repository is currently a production-grade architecture scaffold with executable domain tests and saga tests. It does not yet run full production infra adapters (Kafka brokers, external DBs, etc.) out of the box.

Use these commands as your default loop:

```bash
go test ./...
go test ./tests/saga/...
go test ./tests/contracts/...
```

## Day-to-Day Development

### Typical Workflow

1. Pick the bounded context to change.
2. Update domain + application behavior in that context only.
3. Update event contract/migration only for the context you changed.
4. Add or update tests (unit + saga/contract where relevant).
5. Run test suite.
6. Commit with a Conventional Commit message.

### Bounded Context Rule of Thumb

- Do not import another context's internal domain types.
- Use ACL or events at boundaries.
- Never share storage across bounded contexts.

### When to Update Which Files

- Business behavior change:
	- `services/<context>/internal/domain/*`
	- `services/<context>/internal/application/*`
- New event payload:
	- `contracts/proto/<context>/events.proto`
	- `contracts/kafka-topics.yaml`
- Storage shape change:
	- `infra/migrations/<context>/*`
- Migration rollout policy:
	- `deploy/facade/flags/phase1.yaml`
	- `deploy/reliability/rollout-gates.yaml`
- Operational behavior:
	- `docs/runbooks/*`
	- `docs/slo/*`

## Repository Structure

- `services/*` - Bounded contexts
- `contracts/proto` - Event contracts
- `contracts/kafka-topics.yaml` - Topic ownership and schema mapping
- `infra/migrations` - Per-context migration files
- `deploy/facade` - Strangler routing and rollout policy
- `deploy/federation` - Supergraph composition artifact
- `deploy/gateway` - Gateway policy baseline
- `deploy/reliability` - Rollout gates
- `deploy/observability` - Alert policies
- `deploy/security` - Security baseline
- `tests/saga` - Cross-context flow validation
- `tests/contracts` - Legacy API compatibility shape tests
- `docs/adr` - Architectural decisions
- `docs/runbooks` - Cutover/rollback operations
- `docs/slo` - Reliability objectives and error budget policy

## Local Validation

Run all tests:

```bash
go test ./...
```

Run coverage:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Focus suites:

```bash
go test ./tests/saga/...
go test ./tests/contracts/...
```

## Migration and Operations

### Progressive Rollout
- Controlled via facade flag policy in `deploy/facade/flags/phase1.yaml`.
- Use gate checks from `deploy/reliability/rollout-gates.yaml`.

### Contract Safety
- Legacy API response shape checks are in `tests/contracts/legacy_contract_test.go`.

### Runbooks
- Cutover procedure: `docs/runbooks/cutover.md`
- Rollback procedure: `docs/runbooks/rollback.md`

## Security, Reliability, and SLOs

- Security baseline: `deploy/security/policies.yaml`
- Alert definitions: `deploy/observability/alerts-rules.yaml`
- SLO and error budget policy: `docs/slo/phase5-slo.md`
- CI quality gates: `.github/workflows/phase5-quality-gates.yml`

## Troubleshooting

### "expected declaration, found 'package'"

This usually means a file accidentally contains duplicate package declarations. Open the file and ensure it starts with a single `package <name>` line.

### Tests pass but coverage gate fails in CI

Run locally:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -n 1
```

Ensure total coverage remains above the quality gate configured in `.github/workflows/phase5-quality-gates.yml`.

### Strangler rollout concerns

Before increasing traffic percentages:

1. Run saga and contract tests
2. Verify SLOs and alert health
3. Use the cutover/rollback runbooks in `docs/runbooks`

## Current Implementation Status

Implemented:
- Full Phase 1-5 architecture scaffold
- Domain logic and in-memory infrastructure for all bounded contexts
- Saga and contract tests
- Migration hardening and operational controls

Still to productionize further:
- Real runtime adapters for Kafka/DB/Redis/Elastic/ClickHouse/Stripe/3PL
- Full GraphQL resolver/runtime composition and production routers
- Real infra provisioning and managed-cloud environment wiring

## ADR Index

- `docs/adr/0001-ddd-bounded-context-isolation.md`
- `docs/adr/0002-strangler-facade-feature-flags.md`
- `docs/adr/0003-ordering-saga-choreography.md`
- `docs/adr/0004-phase2-checkout-saga.md`
- `docs/adr/0005-phase3-search-reviewtrust.md`
- `docs/adr/0006-phase4-analytics-federation.md`
- `docs/adr/0007-phase5-migration-hardening.md`
