# ADR-0001: Single binary with internal boundaries

## Status
Accepted

## Context
The project is intentionally compact, but it should read like an internal service rather than a script.

## Decision
Use a single Go binary with separated internal packages for auth, policy, audit, transport, and orchestration.

## Consequences
- Easier local development
- Clear package boundaries
- Straightforward future extraction if the service grows
