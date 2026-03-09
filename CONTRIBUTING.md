# Contributing

Thanks for your interest in contributing to the GoDaddy Terraform Provider.

## Prerequisites

- [Go](https://go.dev/dl/) 1.22+
- [Terraform](https://developer.hashicorp.com/terraform/downloads) 1.0+
- A [GoDaddy API key and secret](https://developer.godaddy.com/keys)
- [golangci-lint](https://golangci-lint.run/welcome/install/) (for linting)

## Getting Started

Fork and clone the repo, then install dependencies:

```sh
git clone https://github.com/<your-fork>/terraform-provider-godaddy.git
cd terraform-provider-godaddy
go mod tidy
```

Build and install the provider locally:

```sh
make local
```

This compiles the binary and places it in `~/.terraform.d/plugins/` so Terraform can discover it.

## Running Tests

Unit tests:

```sh
make test
```

Integration / acceptance tests require real GoDaddy credentials. Set the following environment variables before running them:

```sh
export GODADDY_API_KEY="your-key"
export GODADDY_API_SECRET="your-secret"
export GODADDY_DOMAIN="your-test-domain.com"
```

Then run with the `integration` build tag:

```sh
go test ./... -v -tags=integration
```

> **Warning:** Acceptance tests create and modify real DNS records. Use a domain you control and don't mind being modified.

## Code Style

This project uses [golangci-lint](https://golangci-lint.run/) for static analysis. Run it before submitting a PR:

```sh
golangci-lint run
```

Beyond that, follow standard Go conventions:
- `gofmt` / `goimports` formatted code
- Exported types and functions have doc comments
- Meaningful variable names, short-lived variables can be terse
- Error handling — don't discard errors silently

## Submitting Changes

1. Fork the repository.
2. Create a feature branch from `master` (`git checkout -b my-feature`).
3. Make your changes, add tests where applicable.
4. Run `make test` and `golangci-lint run` to verify.
5. Commit with a clear, descriptive message.
6. Open a pull request against `master`.

Keep PRs focused — one logical change per PR makes review easier.

## A Note on Acceptance Tests

Acceptance tests hit the live GoDaddy API. They are not run in CI by default. If your change affects API interactions (resources, data sources, or the client), please run the relevant acceptance tests locally before submitting. Reviewers may ask for test output as part of the review.
