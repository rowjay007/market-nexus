# ADR 0005: Phase 3 Search and Review/Trust Domains

## Status
Accepted

## Context
Phase 3 adds search discovery and trust signals while preserving bounded-context ownership and tenant isolation.

## Decision
Search BC owns indexable documents and ranking behavior scoped by vendor. Review/Trust BC owns reviews, disputes, and trust metrics. Integration uses domain events and no cross-context DB joins.

## Consequences
Improves product discovery and trust visibility without violating vendor isolation or DDD boundaries.
