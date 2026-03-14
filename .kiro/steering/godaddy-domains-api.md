---
inclusion: auto
---

# GoDaddy Domains API Reference

This provider wraps the [GoDaddy Domains API](https://developer.godaddy.com/doc/endpoint/domains).

Canonical Swagger spec: #[[file:https://developer.godaddy.com/swagger/swagger_domains.json]]

Base URL (production): `https://api.godaddy.com`
Base URL (OTE/testing): `https://api.ote-godaddy.com`

## Authentication

All requests use `sso-key {key}:{secret}` in the `Authorization` header.
Optional `X-Shopper-Id` header for reseller/delegated access (mapped to `customer` attribute in the provider).

Rate limit: 60 requests/minute (enforced client-side at 1 req/sec via `rateLimitedTransport` in `api/client.go`).

## Current Provider Coverage

The provider currently implements a single resource: `godaddy_domain_record`.

### Implemented Endpoints

| API Endpoint | Method | Provider Usage | Code Location |
|---|---|---|---|
| `/v1/domains` | GET | `GetDomains()` — list all domains | `api/domains.go` |
| `/v1/domains/{domain}` | GET | `GetDomain()` — fetch domain details, used to resolve domain ID | `api/domains.go` |
| `/v1/domains/{domain}/records?limit=&offset=` | GET | `GetDomainRecords()` — paginated fetch of all records | `api/domains.go` |
| `/v1/domains/{domain}/records/{type}` | PUT | `UpdateDomainRecords()` — replace records by type | `api/domains.go` |

### Not Yet Implemented

These endpoints exist in the Swagger spec but have no provider resource or data source:

**Domain lifecycle:**
- `POST /v1/domains/purchase` — register a new domain
- `POST /v1/domains/{domain}/renew` — renew a domain
- `DELETE /v1/domains/{domain}` — cancel a domain
- `PATCH /v1/domains/{domain}` — update domain settings (lock, auto-renew, nameservers, WHOIS exposure)
- `POST /v1/domains/{domain}/transfer` — transfer a domain in

**DNS records (granular):**
- `PUT /v1/domains/{domain}/records` — replace ALL records at once
- `PATCH /v1/domains/{domain}/records` — add records (append, not replace)
- `GET /v1/domains/{domain}/records/{type}/{name}` — get records by type+name
- `PUT /v1/domains/{domain}/records/{type}/{name}` — replace records by type+name
- `DELETE /v1/domains/{domain}/records/{type}/{name}` — delete records by type+name

**Domain info & availability:**
- `GET /v1/domains/available` — check domain availability
- `POST /v1/domains/available` — bulk availability check
- `GET /v1/domains/suggest` — domain name suggestions
- `GET /v1/domains/tlds` — list supported TLDs

**Contacts:**
- `PATCH /v1/domains/{domain}/contacts` — update domain contacts
- `POST /v1/domains/contacts/validate` — validate contacts

**Privacy:**
- `DELETE /v1/domains/{domain}/privacy` — cancel privacy
- `POST /v1/domains/{domain}/privacy/purchase` — purchase privacy

**Agreements & validation:**
- `GET /v1/domains/agreements` — legal agreements for TLD
- `GET /v1/domains/purchase/schema/{tld}` — purchase schema
- `POST /v1/domains/purchase/validate` — validate purchase request

**v2 endpoints (newer API):**
- `GET /v2/customers/{customerId}/domains/{domain}` — detailed domain info (v2)
- `PUT /v2/customers/{customerId}/domains/{domain}/nameServers` — update nameservers
- `PATCH /v2/customers/{customerId}/domains/{domain}/dnssecRecords` — add DNSSEC
- `DELETE /v2/customers/{customerId}/domains/{domain}/dnssecRecords` — remove DNSSEC
- `POST /v2/customers/{customerId}/domains/register` — register domain (v2)
- `POST /v2/customers/{customerId}/domains/{domain}/renew` — renew (v2)
- `POST /v2/customers/{customerId}/domains/{domain}/transfer` — transfer (v2)
- `POST /v2/customers/{customerId}/domains/{domain}/redeem` — redeem from redemption
- Domain forwarding CRUD under `/v2/customers/{customerId}/domains/forwards/{fqdn}`
- Actions, notifications, maintenance, and usage endpoints under `/v2/`

## Key Data Types

### DNSRecord (from Swagger)

```json
{
  "type": "A|AAAA|CNAME|MX|NS|SOA|SRV|TXT",
  "name": "string (required)",
  "data": "string (required)",
  "ttl": "integer",
  "priority": "integer (MX/SRV)",
  "weight": "integer (SRV)",
  "port": "integer 1-65535 (SRV)",
  "service": "string (SRV)",
  "protocol": "string (SRV)"
}
```

### Provider DomainRecord (api/types.go)

Maps 1:1 to the Swagger `DNSRecord` definition. Validation in `NewDomainRecord()`:
- Name: 1–255 octets, each label ≤63 chars
- Data: ≤512 chars for TXT, ≤255 for others (SRV exempt)
- TTL: must be ≥0
- Priority: 0–65535
- Weight: 0–100
- Port: 1–65535
- Service/Protocol: must start with `_` if non-empty

### Domain (api/types.go)

```go
type Domain struct {
    ID     int64  `json:"domainId"`
    Name   string `json:"domain"`
    Status string `json:"status"`
}
```

Only `domainId`, `domain`, and `status` are captured. The Swagger `DomainDetail` has many more fields (contacts, privacy, expiration, nameservers, etc.) that could be exposed as computed attributes or a data source.

## Provider Architecture

```
main.go                          → plugin entrypoint
plugin/godaddy/provider.go       → schema.Provider (key, secret, baseurl)
plugin/godaddy/config.go         → Config struct, creates api.Client
plugin/godaddy/resource_dns_record.go → godaddy_domain_record CRUD
api/client.go                    → HTTP client with rate limiting
api/domains.go                   → API methods (GetDomain, GetDomainRecords, UpdateDomainRecords)
api/types.go                     → Domain, DomainRecord, validation functions
```

## Implementation Notes

- `UpdateDomainRecords` iterates over supported types and PUTs each type separately via `/v1/domains/{domain}/records/{type}`. It skips empty NS, SOA, and CAA lists (see `IsDisallowed`).
- On delete (`resourceDomainRecordRestore`), the provider resets to `defaultRecords` (www CNAME → @, _domainconnect CNAME).
- Pagination in `GetDomainRecords` uses `limit=500` and increments `offset` until an empty page is returned.
- The provider uses `terraform-plugin-sdk/v2` (v2.35.0) with `ConfigureContextFunc`.
