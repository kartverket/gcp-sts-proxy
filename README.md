# gcp-sts-proxy

A thin HTTP wrapper around [`golang.org/x/oauth2/google/externalaccount`](https://pkg.go.dev/golang.org/x/oauth2/google/externalaccount) that injects GCP access tokens into proxied requests. 
Intended to run in Kubernetes pods using Workload Identity Federation. 
Requires the cluster to be registered as a trusted OIDC identity provider in GCP, and the pod to have a projected service account token mounted.

## How it works

```
Client                    gcp-sts-proxy                       GCP
  │                            │                               │
  │  GET /?url=<target-url>    │                               │
  ├──────────────────────────► │                               │
  │                            │  exchange k8s SA token        │
  │                            ├──────────────────────────────►│ sts.googleapis.com
  │                            │◄──────────────────────────────┤ (access token)
  │                            │                               │
  │                            │  GET <target-url>             │
  │                            │  Authorization: Bearer <tok>  │
  │                            ├──────────────────────────────►│
  │                            │◄──────────────────────────────┤
  │◄───────────────────────────┤                               │
```

1. Client sends any HTTP request to the proxy with a `?url=<target>` query parameter.
2. Proxy reads the Kubernetes service account token from disk (`TOKEN_FILE`).
3. Proxy exchanges it for a GCP access token via [Workload Identity Federation](https://cloud.google.com/iam/docs/workload-identity-federation) (`sts.googleapis.com`).
4. If `IMPERSONATION_URL` is set, the STS token is further exchanged for a service account token.
5. Proxy forwards the original request to `<target>` with `Authorization: Bearer <token>`.
6. Response is streamed back to the client unchanged.

## Environment variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `AUDIENCE` | yes | — | Workload Identity audience. For GDC: `identitynamespace:<k8s-project>.svc.id.goog:https://gkehub.googleapis.com/projects/<k8s-project>/locations/<cluster-region>/memberships/<cluster-name>` |
| `TOKEN_FILE` | no | `/var/run/secrets/tokens/gcp-ksa/token` | Path to projected Kubernetes service account token |
| `IMPERSONATION_URL` | no | — | Service account impersonation URL. If set, STS token is exchanged for a SA token. Format: `https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/<sa-email>:generateAccessToken` |
| `PORT` | no | `8080` | Port the proxy listens on |

## Usage

```
GET http://localhost:8080/?url=https://storage.googleapis.com/...
```

All request methods, headers, and body are forwarded. The `Authorization` header is always overwritten with the GCP token.
