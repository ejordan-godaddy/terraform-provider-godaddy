# terraform-provider-godaddy

[Terraform](https://www.terraform.io/) provider for managing GoDaddy domain records.

[![Release](https://github.com/ejordan-godaddy/terraform-provider-godaddy/workflows/release/badge.svg)](https://github.com/ejordan-godaddy/terraform-provider-godaddy/actions)
[![Test](https://github.com/ejordan-godaddy/terraform-provider-godaddy/workflows/test/badge.svg)](https://github.com/ejordan-godaddy/terraform-provider-godaddy/actions)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.22 (for building from source)

## API Key

An [API key](https://developer.godaddy.com/keys/) is required to use the GoDaddy APIs. The key pair can be stored in environment variables:

```bash
export GODADDY_API_KEY=abc
export GODADDY_API_SECRET=123
```

## Provider Configuration

If `key` and `secret` aren't provided in the provider block, they are expected as environment variables `GODADDY_API_KEY` and `GODADDY_API_SECRET`.

```terraform
provider "godaddy" {
  key    = "abc"
  secret = "123"
}
```

## Domain Record Resource

A `godaddy_domain_record` resource requires a `domain`. If the domain is not registered under the account that owns the key, an optional `customer` number can be specified.

Records can be defined using the `record` block, or using shorthand `addresses` and `nameservers` lists. Supported record types: A, AAAA, CAA, CNAME, MX, NS, SOA, SRV, TXT.

```terraform
resource "godaddy_domain_record" "gd-fancy-domain" {
  domain   = "fancy-domain.com"
  customer = "1234" # required if provider key does not belong to customer

  record {
    name = "www"
    type = "CNAME"
    data = "fancy.github.io"
    ttl  = 3600
  }

  record {
    name     = "@"
    type     = "MX"
    data     = "aspmx.l.google.com."
    ttl      = 600
    priority = 1
  }

  record {
    name     = "@"
    type     = "SRV"
    data     = "host.example.com"
    ttl      = 3600
    service  = "_ldap"
    protocol = "_tcp"
    port     = 389
  }

  addresses   = ["192.168.1.2", "192.168.1.3"]
  nameservers = ["ns7.domains.com", "ns6.domains.com"]
}
```

## Import

Existing domain records can be imported:

```bash
terraform import godaddy_domain_record.gd-fancy-domain fancy-domain.com
```

> If your zone contains existing data, ensure your Terraform configuration includes all existing records, otherwise they will be removed.

## Building

```bash
# Build and install locally
make local

# Build for Linux via Docker
make linux

# Run tests
make test
```

## License

[Apache License 2.0](LICENSE)
