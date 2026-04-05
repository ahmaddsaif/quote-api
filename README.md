# Quote API (Golang)

This repository is now tracked as part of our **Learning Golang with Practical AWS** project.

## Purpose in the learning project

This service gives us a compact Go API we can use to practice AWS deployment and operations patterns end-to-end.

## Current app scope

- `GET /quote` — returns a random quote (with Redis caching).
- `GET /health` — basic health endpoint.

## AWS practice roadmap for this repo

1. Containerize and publish image (ECR).
2. Deploy on ECS Fargate (or EKS) behind ALB.
3. Add managed Redis (ElastiCache) and config via environment variables.
4. Add edge protections (WAF/rate limiting) and observability (CloudWatch).
5. Add CI/CD (GitHub Actions + deploy pipeline).

## Local run

```bash
docker compose up --build
```

Then call:

- `http://localhost:8080/quote`
- `http://localhost:8080/health`
