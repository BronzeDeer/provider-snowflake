# DataProduct

Crossplane v2 composition bundling the provider's Snowflake resources into a single namespaced `DataProduct` (`platform.bronze-deer.de/v1alpha1`).

## Inputs

| Field | Description |
|---|---|
| `name` | Name of the data product, 3 to 20 characters |
| `id` | Numeric ID, greater than 0 |
| `env` | One of `INTE`, `ZINT`, `GSEC`, `KOMO`, `MIGR`, `TEST`, `DEV`, `PROD`, `ABNA` |
| `division` | One of `PNC`, `HEALTH` |
| `warehouse.size` | Warehouse size, defaults to `XSMALL` |
| `warehouse.autoscale` | Multi-cluster warehouse (1 to 3 clusters), defaults to `false` |
| `viewers` | Snowflake users granted the viewer role |
| `contributors` | Snowflake users granted the developer role |
| `maintainers` | Snowflake users granted the maintainer role |

## Composed Snowflake objects

For `id: 1001`, `env: DEV`, `division: PNC`, `name: sales-orders`:

| Object | Name |
|---|---|
| Database | `D1001_DEV` |
| Schemas | `SDP`, `OPS`, `OUT`, `DPZ` |
| Warehouse | `D1001_DEV_WH` |
| Account roles | `D1001_DEV_PNC_SALES_ORDERS_{D,M,V,C,X}` (developer, maintainer, viewer, consumer, x) |
| Python procedures (stubs) | `SDP.CREATE_DEV_SCHEMA(VARCHAR)`, `SDP.CREATE_DBT_TASK(VARIANT)` |
| Service users | `SVC_D1001_DEV_GITHUB_CI`, `SVC_D1001_DEV`, both granted the developer role |

Role hierarchy: developer > maintainer > viewer, consumer and developer > x. The x role is also granted to `SVC_D1001_DEV_GITHUB_CI`.

The privileges of each role are defined in the `privileges` template at the top of the composition. Schema object privileges are only granted on future objects.

## Install

```console
kubectl apply -f functions.yaml
kubectl apply -f definition.yaml
kubectl apply -f composition.yaml
kubectl apply -f example-users.yaml
kubectl apply -f example.yaml
```

`example.yaml` uses all inputs, its users must already exist in Snowflake, `example-users.yaml` creates them. `example-minimal.yaml` only sets the required inputs.

Composed resources use the `ClusterProviderConfig` named `default`.

When running the provider out of cluster with `make run`, Crossplane is missing the RBAC a provider package would install, apply `local-dev.yaml` (RBAC and a default `ClusterProviderConfig` using the `example-creds` secret) first.

## Render locally

```console
crossplane render example.yaml composition.yaml functions.yaml --xrd=definition.yaml
```
