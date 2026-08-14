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

## Project structure

```
identirat/   IAM service — accounts, roles, JWT issuance
notariat/    Parish registry, register books, baptism records and PDF certificates
utils/       Shared Go module — JWT middleware, config and parsing
webserver/   Simple Caddy reverse proxy config
docker/      Postgres initialization SQL scripts per module
```
