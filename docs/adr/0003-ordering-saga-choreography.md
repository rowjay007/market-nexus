# ADR 0003: Choreography Saga for Phase 1

## Status
Accepted

## Context
Ordering flow must avoid central orchestrator and support compensating actions.

## Decision
Ordering emits OrderPlaced and reacts to Inventory outcomes. On reservation failure, OrderCancelled is emitted and previously reserved inventory is released.

## Consequences
Independent service autonomy and resilient eventual consistency behavior.
