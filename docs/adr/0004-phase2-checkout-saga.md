# ADR 0004: Phase 2 Checkout Saga Choreography

## Status
Accepted

## Context
Phase 2 introduces Pricing, Payment, and Fulfillment bounded contexts to complete checkout.

## Decision
Ordering initiates checkout by confirming inventory reservation, requesting pricing quote, capturing payment, and scheduling fulfillment. Compensation sequence on downstream failure includes payment refund and inventory release before leaving order in cancelled state.

## Consequences
Improves service autonomy and keeps context boundaries strict while supporting resilient eventual consistency.
