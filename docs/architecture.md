# Architecture

## Overview

Gatekeeper is a single-process HTTP service with explicit internal boundaries.
The repository is intentionally organized like a service that could live inside a larger platform stack.

## Components

### HTTP server
Owns routing, middleware chaining, and transport concerns.

### Auth module
Resolves a request subject from a bearer token. Authentication is kept isolated from policy logic.

### Access service
Builds the authorization request context and coordinates policy evaluation with audit emission.

### Policy engine
Accepts the subject, action, resource, and attributes. Returns a decision, rule name, and human-readable reason.

### Audit sink
Captures every decision in append-only form. The default implementation is in-memory with bounded capacity.

## Decision path

```text
HTTP request
  -> auth middleware
  -> authorize handler
  -> access service
  -> policy engine
  -> audit sink
  -> response
```

## Boundary rules

- handlers do not contain policy logic
- policy engine does not know about HTTP
- audit sink does not decide access
- config is loaded at startup and passed inward

## Why this shape

The main objective is to show strong repository style:

- transport isolated from business logic
- decisioning isolated from persistence
- docs and API contract versioned with code
- sample data included for reproducibility
- tests focused on policy and token behavior
