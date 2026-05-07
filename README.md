# gcp-sts-proxy

A thin HTTP proxy that injects GCP access tokens into proxied requests using [Application Default Credentials](https://cloud.google.com/docs/authentication/application-default-credentials).
Intended to run in Kubernetes pods with Workload Identity enabled.

## How it works

```
Client                    gcp-sts-proxy                       GCP
  │                            │                               │
  │  GET /?url=<target-url>    │                               │
  ├──────────────────────────► │                               │
  │                            │  get token                    │
  │                            ├──────────────────────────────►│
  │                            │◄──────────────────────────────┤
  │                            │                               │
  │                            │  GET <target-url>             │
  │                            │  Authorization: Bearer <tok>  │
  │                            ├──────────────────────────────►│
  │                            │◄──────────────────────────────┤
  │◄───────────────────────────┤                               │
```

1. Client sends any HTTP request to the proxy with a `?url=<target>` query parameter.
2. Proxy fetches a GCP access token via `google.DefaultTokenSource` (ADC).
3. Proxy forwards the original request to `<target>` with `Authorization: Bearer <token>`.
4. Response is streamed back to the client unchanged.

## Environment variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | no | `8080` | Port the proxy listens on |

Token acquisition is handled entirely by ADC. Configure credentials via [standard ADC mechanisms](https://cloud.google.com/docs/authentication/application-default-credentials) (Workload Identity, `GOOGLE_APPLICATION_CREDENTIALS`, `gcloud auth application-default login`, etc.).

## Usage

```
GET http://localhost:8080/?url=https://storage.googleapis.com/...
```

All request methods, headers, and body are forwarded. The `Authorization` header is always overwritten with the GCP token.
