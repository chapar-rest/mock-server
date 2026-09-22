mock-server k8s setup

----
Runs on a k3s cluster with the default Traefik ingress. One pod; all state
is in memory.

#### Create namespace

```bash
kubectl apply -f k8s/namespace.yaml
```

#### Apply manifests

```bash
kubectl apply -f k8s/middleware.yaml -f k8s/deployment.yaml -f k8s/service.yaml -f k8s/ingress.yaml
```

After this, deploys only change the image: `task publish` (or the Deploy
workflow on every push to `main`).

#### DNS / TLS
Point `mocks.chapar.rest` at the cluster. For gRPC through a TLS-terminating
proxy in front of Traefik (e.g. Cloudflare), the proxy must allow gRPC /
HTTP/2 to the origin.

#### Rate limiting
`k8s/middleware.yaml` limits each client (by `CF-Connecting-IP`) to 10 req/s
(burst 30) and 10 concurrent requests. The key relies on traffic coming
through the Cloudflare proxy: restrict the origin to Cloudflare IPs, or a
client can bypass the proxy and the limits.
