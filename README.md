# Zignurat

Identity and access management, built as a lightweight, self-hosted Go microservice — handles authentication, accounts and roles.

## Local development

Run the following commands to get the service running locally.

```bash
cp .env.example .env
docker compose up --build
```

### Local URLs

- **Identirat**: http://identirat.localhost

## Tech stack

- Go, Gin, GORM + PostgreSQL
- Docker Compose
- JWT-based auth shared through the `utils` module
- Caddy as reverse proxy

## Project structure

```
identirat/   IAM service — accounts, roles, JWT issuance
utils/       Shared Go module — JWT middleware, config and parsing
webserver/   Simple Caddy reverse proxy config
docker/      Postgres initialization SQL scripts
```
