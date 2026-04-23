# gatekeeper-go

A compact internal access gateway in Go with policy enforcement, audit logging, and service-token authentication.

## Why this exists

Gatekeeper is a production-style sample service that sits in front of internal endpoints and decides whether a request should be allowed, denied, or challenged. The codebase is intentionally small, but the repository is structured like a real service: configuration, middleware, policy evaluation, audit trails, tests, architecture notes, and an API contract.

## Core capabilities

- Service token authentication
- Policy-based access control
- Request audit logging
- Health and readiness endpoints
- Versioned HTTP API
- Config-driven bootstrapping
- OpenAPI contract
- Sample policy data

## Repository layout

```text
cmd/gatekeeper/           application entrypoint
internal/auth/            token parsing and request identity
internal/policy/          policy engine and rule evaluation
internal/audit/           audit event model and in-memory sink
internal/httpserver/      router, middleware, handlers
internal/service/access/  orchestration layer
internal/config/          configuration loading
api/openapi.yaml          API contract
configs/config.yaml       default service config
docs/architecture.md      system design notes
docs/adr/                 architecture decision records
sample_data/              example policy input
```

## Request flow

1. Client sends request with bearer token.
2. Authentication middleware resolves subject identity.
3. Access service builds an authorization context.
4. Policy engine evaluates action/resource match.
5. Decision is written to the audit sink.
6. HTTP handler returns decision payload.

## Quick start

```bash
make run
```

Then call:

```bash
curl -s http://localhost:8080/v1/authorize \
  -H 'Authorization: Bearer svc:ops-bot|roles=ops,readonly' \
  -H 'Content-Type: application/json' \
  -d '{
    "action": "read",
    "resource": "incidents",
    "attributes": {
      "environment": "prod"
    }
  }'
```

## Example response

```json
{
  "decision": "allow",
  "reason": "matched policy rule",
  "trace_id": "b1f0f2b1-3e6a-41d5-a39e-4ecb5223b6be"
}
```

## Tokens

This repository uses a deliberately simple service-token format for demonstration:

```text
svc:<subject>|roles=ops,readonly
```

Example:

```text
svc:deploy-agent|roles=ops,deploy
```

## Endpoints

- `GET /healthz`
- `GET /readyz`
- `POST /v1/authorize`
- `GET /v1/audit/events`

## Configuration

Default configuration lives in `configs/config.yaml`.

Supported settings:

- service name
- bind address
- audit sink size
- default decision mode
- policy file path

## Architecture notes

See `docs/architecture.md` for the service boundaries and data flow.

## Design goals

- Keep the codebase small, but structured like a maintainable internal service
- Make policy decisions explicit and testable
- Keep handlers thin and business logic isolated
- Preserve a clear audit trail for every authorization decision

## Running tests

```bash
make test
```

## Future work

- persistent audit storage
- signed JWT validation
- policy hot reload
- metrics export
- structured policy DSL

## License

MIT
