# ADR 0002: Strangler Fig Facade with Progressive Routing

## Status

Accepted

## Context

Migration from PHP monolith must preserve API compatibility with zero breaking changes.

## Decision

Nginx facade proxies all routes with feature-flag controlled rollout per endpoint (0 -> 10 -> 50 -> 100).

## Consequences

Safe incremental migration with observability gates before full cutover.
