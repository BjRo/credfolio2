---
# credfolio2-n907
title: Configure environment variables loading
status: todo
type: task
created_at: 2026-01-20T11:26:27Z
updated_at: 2026-01-20T11:26:27Z
parent: credfolio2-jpin
---

Set up environment-based configuration for the Go backend.

## Requirements
- Load configuration from environment variables
- Support .env files for local development
- Validate required configuration at startup
- Type-safe configuration struct

## Configuration Values
- DATABASE_URL (PostgreSQL connection string)
- MINIO_ENDPOINT, MINIO_ACCESS_KEY, MINIO_SECRET_KEY
- ANTHROPIC_API_KEY
- SERVER_PORT (default: 8080)
- LOG_LEVEL (default: info)
- ENV (development, staging, production)

## Acceptance Criteria
- Application fails fast if required config missing
- Config values accessible throughout application
- .env.example documents all variables
- Sensitive values not logged

## Technical Notes
- Consider using envconfig, viper, or similar
- No hardcoded values in source code