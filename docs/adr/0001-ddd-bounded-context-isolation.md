# ADR 0001: Bounded Context Isolation and DDD

## Status

Accepted

## Context

MarketNexus requires strict domain ownership, no cross-context joins, and anti-corruption at boundaries.

## Decision

Each bounded context (Catalog, Inventory, Ordering) owns its model and datastore. Integration is event-driven through Kafka topics and ACL projections.

## Consequences

Strong decoupling and migration safety. Additional projection and consistency complexity is accepted.
