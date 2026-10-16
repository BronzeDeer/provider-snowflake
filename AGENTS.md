# AGENTS.md

Guidance for agents working in this repository.

## What this repository is

`provider-snowflake` is a [Crossplane](https://crossplane.io/) provider for
Snowflake. It is generated with [Upjet](https://github.com/crossplane/upjet)
(`upjet/v2 v2.2.0`), which wraps the Terraform provider
`snowflakedb/snowflake` (pinned via `TERRAFORM_PROVIDER_VERSION` in the
`Makefile`, currently `2.21.0`). It depends on
`crossplane-runtime/v2 v2.1.0` and targets the Crossplane v2 API style.

Key facts:

- Kubernetes API group suffix: `snowflake.bronze-deer.de` (cluster-scoped) and
  `snowflake.m.bronze-deer.de` (namespaced).
- Resource groups: `database`, `schema`, `rbac`, `warehouse`, `procedure`.
- Managed resources: `Database`, `Schema`, `AccountRole`, `AccountRoleGrant`,
  `AccountRolePrivilegeGrant`, `User`, `ServiceUser`, `Warehouse`,
  `ProcedurePython`.

## Repository layout (human-written vs generated)

Human-written (safe to edit):

- `config/` — Upjet configuration. Per-resource tweaks live in
  `config/cluster/<resource>/config.go` and
  `config/namespaced/<resource>/config.go`; global reference/external-name
  config is in `config/external_name.go` and `config/provider.go`.
- `cmd/` — entrypoints, including the generator (`cmd/generator/main.go`).
- `internal/` — controllers/clients/features (the `zz_*` files inside are
  generated; everything else is hand-written).
- `examples/` — hand-maintained example manifests.
- `cluster/test/setup.sh` — e2e setup (credentials secret + default provider
  configs).
- `integration-tests/` — chainsaw test suites (see below).
- `Makefile`, `flake.nix`, `build/`, `hack/`, `scripts/`.

Generated (do NOT edit by hand; regenerate instead):

- `apis/**/zz_*.go` — types, terraformed metadata, deepcopy, managed
  resource methods, and reference resolvers.
- `internal/controller/**/zz_*.go` — controllers and `zz_setup.go`.
- `package/crds/*.yaml` — CRD manifests.
- `examples-generated/` — examples scraped from Terraform docs.
- `config/provider-metadata.yaml`, `config/schema.json`.

## CRDs and APIs

- CRDs: `package/crds/`. Each resource exists twice: cluster-scoped
  (`<group>.snowflake.bronze-deer.de_<plural>.yaml`) and namespaced
  (`<group>.snowflake.m.bronze-deer.de_<plural>.yaml`).
- Go types: `apis/cluster/<group>/v1alpha1/` and
  `apis/namespaced/<group>/v1alpha1/`. Provider config types use `v1beta1`.

## Crossplane-runtime references (`value` / `valueRef` / `valueLabelSelector`)

For most scalar or object fields that reference another managed resource, Upjet
generates a triplet:

- `<field>` — the literal value to use (e.g. `database: my-db`).
- `<field>Ref` — a `*v1.Reference` (`{name}`) pointing at a managed resource of
  the same provider to source the value from.
- `<field>Selector` — a `*v1.Selector` (label selector) matching managed
  resources to source the value from.

Resolution happens in generated `ResolveReferences` methods in
`apis/**/zz_generated.resolvers.go`, which call
`reference.NewAPIResolver(...).Resolve(...)`. The value is extracted from the
referenced resource's **external-name** (`reference.ExternalName()`), not its
`metadata.name` or `spec.forProvider.name`. This is why external-name
configuration in `config/external_name.go` matters for references (e.g. a
`Schema`'s external name is its schema name, so a `ProcedurePython.schemaRef`
resolves to the schema name).

## Crossplane semantics (eventual consistency)

Crossplane and the generated provider are eventually consistent. Each managed
resource reconciles independently: once its required fields are populated, its
controller syncs it against the external API and retries until it converges,
including resolving `valueRef` / `valueLabelSelector` references once the
referenced resource exists. As a result there is typically **no explicit
ordering dependency** between managed resources, and none should be assumed;
each resource converges on its own schedule.

This applies beyond this repo (e.g. when writing chainsaw tests): do not encode
apply or delete order to satisfy a perceived dependency. Sequential applies are
fine as long as all applies happen before the assertions; the assertions are the
part that waits for convergence and verify the desired converged state.

## Code generation

The pipeline is declared in `apis/generate.go`:

```console
go generate ./apis/...
# or
go run cmd/generator/main.go "$PWD"
```

It scrapes Terraform docs, runs the Upjet generator, then runs `controller-gen`
(CRDs) and `angryjet` (managed-resource method sets).

## Tool versions

The `Makefile` (via `build/makelib/k8s_tools.mk`) defines fallback / minimum
versions for tools such as `chainsaw`. When the `Makefile` version and the
version pinned by `flake.nix` / the devShell disagree, the devShell version
takes precedence.

## Integration tests (chainsaw)

Location and layout:

- `integration-tests/<resource>/{cluster,namespaced}/`, one suite per resource
  per scope. All non-providerconfig CRDs are covered.
- Each suite is a `chainsaw.kyverno.io/v1alpha1` `Test` plus resource manifests
  and matching `*-assert.yaml` files.

Conventions used by these suites:

- A single binding `testName` built from `(join('-', ['<prefix>', random('[a-z0-9]{8}')]))`.
- Every resource name shares the `$testName` prefix (suffixes like
  `(join('-', [$testName, 'parent']))` for multiple resources of the same kind).
- Every resource carries the label
  `chainsaw.snowflake.bronze-deer.de/test-name: ($testName)`.
- `concurrent: true`; resource manifests all applied before assertions run
  (Crossplane is eventually consistent and resolves refs on retry, so apply and
  delete order need not be encoded; see the eventual-consistency section).
- Ready/Synced assertions use
  `(conditions[?type == 'Ready']): - status: 'True'` (and `Synced`).
- Steps: `create` (apply all -> assert all -> `catch` describe) and `delete`
  (delete all -> `error` against the assert files).
- Prerequisites are created in-suite (e.g. `Schema` needs `Database`;
  `ProcedurePython` needs `Database` + `Schema`; grants need `AccountRole`s).

Setup / expectations:

- Cluster-scoped suites assume the default `ProviderConfig`
  (`snowflake.bronze-deer.de/v1beta1`, name `default`).
- Namespaced suites assume the default `ClusterProviderConfig`
  (`snowflake.m.bronze-deer.de/v1beta1`, kind `ClusterProviderConfig`, name
  `default`).
- `cluster/test/setup.sh` is one way to create the credentials secret and both
  default configs (it consumes `UPTEST_CLOUD_CREDENTIALS` and requires
  `kubectl`). The suites themselves only assume those configs already exist.

Running:

- The suites are plain `chainsaw.kyverno.io/v1alpha1` tests and run with
  `kyverno-chainsaw` alone, e.g.
  `chainsaw test integration-tests/<resource>/<scope>`. They do not depend on
  uptest or `cluster/test/setup.sh`; they only require the provider to be
  installed and the default `ProviderConfig` / `ClusterProviderConfig` to exist.
- `make e2e` is a separate, example-based path (`local-deploy` + `uptest`). The
  top-level `integration-tests/` suites are not wired into the `Makefile`.
- Validate suite YAML without a cluster:
  `yq eval '.' integration-tests/**/*/*.yaml > /dev/null`.

## Working agreements

- Do NOT add comments to YAML or code unless explicitly asked.
- Do NOT run integration tests without asking first.
- Do NOT push/pull/fetch git remotes.
- Only commit on non-main branches.
- Be direct; if a request is unclear or seems off, pause and ask.
- If a required tool or domain is missing, pause and ask the user to add it
  rather than attempting workarounds.
