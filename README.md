# Zignurat

A suite of modules designed to handle identity, registers and documents of individuals and organizations.

- **Identirat** — Identity & Access Management. Owns authentication and account roles, acting as the single point of access to the suite.
- **Notariat** — Parish/diocese registry. Manages register books and baptism records, and generates official PDF certificates (with signature, seal and QR-code verification).

## Local development

Run the following commands to get the suite running locally.

```bash
cp .env.example .env
docker compose up --build
```

### Local URLs

- **Identirat**: http://identirat.localhost
- **Notariat**: http://notariat.localhost

## Tech stack

- Go, Gin, GORM + PostgreSQL (one database per service, provisioned via `docker/*/initdb.d`)
- Docker Compose
- JWT-based auth shared through the `utils` module
- Caddy as reverse proxy
- Kubernetes (minikube), Helm, Gateway API (Envoy Gateway)

## Project structure

```
identirat/   IAM service — accounts, roles, JWT issuance
notariat/    Parish registry, register books, baptism records and PDF certificates
utils/       Shared Go module — JWT middleware, config and parsing
webserver/   Simple Caddy reverse proxy config
docker/      Postgres initialization SQL scripts per module
charts/      Helm chart for the whole suite (apps, databases, routes)
k8s/         Cluster-wide platform setup (GatewayClass)
```

## Kubernetes (minikube)

```bash
# 1. Platform setup (once per cluster)
helm install eg oci://docker.io/envoyproxy/gateway-helm --version v1.8.1 \
  -n envoy-gateway-system --create-namespace
kubectl apply -f k8s/platform/gatewayclass.yaml

# 2. Secrets (once per cluster; never committed)
kubectl create namespace zignurat
kubectl -n zignurat create secret generic zignurat-secrets \
  --from-literal=JWT_SECRET="$(openssl rand -hex 32)" \
  --from-literal=SESSION_SECRET="$(openssl rand -hex 32)" \
  --from-literal=IDENTIRAT_DB_PW=changeme \
  --from-literal=NOTARIAT_DB_PW=changeme

# 3. Build images and deploy
minikube image build -t zignurat/identirat:0.3.0 -f identirat/Dockerfile .
minikube image build -t zignurat/notariat:0.3.0 -f notariat/Dockerfile .
helm install zignurat charts/zignurat -n zignurat

# 4. Access (keep running): http://identirat.localhost, http://notariat.localhost
minikube tunnel
```