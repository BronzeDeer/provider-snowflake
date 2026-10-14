# Provider Snowflake

`provider-snowflake` is a [Crossplane](https://crossplane.io/) provider
snowflake that is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the Snowflake
API.

## Getting Started

This template serves as a starting point for generating a new [Crossplane Provider](https://docs.crossplane.io/latest/packages/providers/) using the [`upjet`](https://github.com/crossplane/upjet) tooling. Please follow the guide linked below to generate a new Provider:

https://github.com/crossplane/upjet/blob/main/docs/generating-a-provider.md

## Developing

Run code-generation pipeline:
```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Integration tests

`integration-tests/` contains [kyverno-chainsaw](https://kyverno.github.io/chainsaw/)
test suites for every managed resource, each in a cluster-scoped and a
namespaced variant:

```console
integration-tests/<resource>/{cluster,namespaced}/
```

The suites only need the provider installed and a default `ProviderConfig`
(cluster-scoped) or `ClusterProviderConfig` (namespaced) named `default`
(`cluster/test/setup.sh` creates both). They do not depend on `uptest`, so they
can be run directly:

```console
chainsaw test integration-tests/warehouse/cluster
chainsaw test integration-tests/warehouse/namespaced
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/BronzeDeer/provider-snowflake/issues).
