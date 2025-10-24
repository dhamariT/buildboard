# Docker Deployment Guide

This guide covers deploying the BuildBoard backend using Docker and Docker Compose.

## Quick Start with Docker Compose

### Development Setup

```bash
# Start all services (backend + PostgreSQL)
make docker-up

# Or with database UI tool (Adminer)
make docker-up-tools
```

The backend will be available at:
- API: http://localhost:8080
- Adminer (DB UI): http://localhost:8081 (if using tools profile)

### View Logs

```bash
make docker-logs
```

### Stop Services

```bash
make docker-down
```

## What's Included

The docker-compose setup includes:

1. **PostgreSQL Database** (postgres:16-alpine)
   - Persistent data volume
   - Health checks
   - Port: 5432

2. **BuildBoard API** (Go backend)
   - Multi-stage build for small image size (~20MB)
   - Non-root user for security
   - Health checks
   - Port: 8080

3. **Adminer** (optional - use `--profile tools`)
   - Web-based database management
   - Port: 8081

## Docker Files Overview

### Dockerfile
Multi-stage build for production:
- **Builder stage**: Compiles Go binary with static linking
- **Runtime stage**: Minimal Alpine Linux image
- **Size**: ~20MB (vs ~1GB with full Go image)
- **Security**: Runs as non-root user

### docker-compose.yml
Development environment with:
- Development settings (GIN_MODE=debug)
- Local database
- Optional email configuration
- Volume persistence

### docker-compose.prod.yml
Production environment with:
- Production settings (GIN_MODE=release)
- Required environment validation
- Resource limits
- SSL enforcement for database

## Environment Configuration

### Development
Edit `docker-compose.yml` directly or use environment variables:

```yaml
environment:
  DB_HOST: postgres
  DB_PASSWORD: postgres
  FRONTEND_URL: http://localhost:3000
  # Email optional for development
```

### Production
Create `.env` file with all required values:

```bash
# Copy template
make setup

# Edit .env
nano .env
```

Required production variables:
```env
ENVIRONMENT=production
DB_PASSWORD=secure_password_here
FRONTEND_URL=https://yourdomain.com
BACKEND_URL=https://api.yourdomain.com
SMTP_HOST=smtp.gmail.com
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password
FROM_EMAIL=your-email@gmail.com
```

## Docker Commands

### Local Development

```bash
# Build Docker image
make docker-build

# Start services
make docker-up

# View logs
make docker-logs

# Check status
make docker-ps

# Restart services
make docker-restart

# Stop services
make docker-down

# Clean up (removes volumes too)
make docker-clean
```

### Production Deployment

```bash
# Build production image
make docker-prod-build

# Start production services
make docker-prod-up

# Stop production services
make docker-prod-down
```

## Manual Docker Commands

If you prefer not to use Make:

### Build Image
```bash
docker build -t buildboard-backend:latest .
```

### Run Container (with existing PostgreSQL)
```bash
docker run -d \
  --name buildboard-api \
  -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_PASSWORD=postgres \
  -e FRONTEND_URL=http://localhost:3000 \
  -e BACKEND_URL=http://localhost:8080 \
  buildboard-backend:latest
```

### Run with Docker Compose
```bash
# Start
docker-compose up -d

# Stop
docker-compose down

# View logs
docker-compose logs -f api

# Rebuild and restart
docker-compose up -d --build
```

## Database Management

### Access Database via Adminer
1. Start services with tools profile:
   ```bash
   make docker-up-tools
   ```
2. Open http://localhost:8081
3. Login with:
   - System: PostgreSQL
   - Server: postgres
   - Username: postgres
   - Password: postgres
   - Database: buildboard_db

### Access Database via psql
```bash
docker-compose exec postgres psql -U postgres -d buildboard_db
```

### Backup Database
```bash
docker-compose exec postgres pg_dump -U postgres buildboard_db > backup.sql
```

### Restore Database
```bash
docker-compose exec -T postgres psql -U postgres buildboard_db < backup.sql
```

## Health Checks

The API container includes health checks:

```yaml
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3
```

Check container health:
```bash
docker ps
# Look for "healthy" in STATUS column
```

## Troubleshooting

### Container won't start
```bash
# Check logs
docker-compose logs api

# Check if database is ready
docker-compose logs postgres
```

### Database connection fails
```bash
# Verify database is healthy
docker-compose ps

# Check database logs
docker-compose logs postgres

# Test database connection
docker-compose exec postgres pg_isready -U postgres
```

### Port already in use
```bash
# Change port in docker-compose.yml
ports:
  - "8081:8080"  # Host:Container
```

### Permission denied errors
The container runs as non-root user (uid 1000). Ensure volumes have correct permissions.

### Email not sending
1. Check SMTP credentials in environment
2. For Gmail: use app-specific password
3. Check logs: `docker-compose logs api | grep -i email`

## Production Deployment Platforms

### Google Cloud Run

```bash
# Build and tag
docker build -t gcr.io/YOUR_PROJECT/buildboard-backend .

# Push to registry
docker push gcr.io/YOUR_PROJECT/buildboard-backend

# Deploy
gcloud run deploy buildboard-backend \
  --image gcr.io/YOUR_PROJECT/buildboard-backend \
  --platform managed \
  --region us-central1 \
  --add-cloudsql-instances YOUR_PROJECT:us-central1:buildboard-db \
  --set-env-vars "DB_HOST=/cloudsql/YOUR_PROJECT:us-central1:buildboard-db"
```

### AWS ECS/Fargate

```bash
# Tag for ECR
docker tag buildboard-backend:latest 123456789.dkr.ecr.us-east-1.amazonaws.com/buildboard-backend

# Push to ECR
docker push 123456789.dkr.ecr.us-east-1.amazonaws.com/buildboard-backend

# Create task definition and service via AWS Console or CLI
```

### Heroku

```bash
# Login to Heroku container registry
heroku container:login

# Build and push
heroku container:push web -a your-app-name

# Release
heroku container:release web -a your-app-name
```

### DigitalOcean App Platform

Upload `Dockerfile` via App Platform dashboard or use `doctl` CLI.

## Resource Requirements

### Minimum (Development)
- CPU: 0.5 cores
- Memory: 256MB
- Disk: 1GB

### Recommended (Production)
- CPU: 1 core
- Memory: 512MB
- Disk: 10GB (with database)

### Database Storage
- Initial: ~10MB
- Growth: ~1KB per user signup

## Security Best Practices

1. **Non-root user**: Container runs as uid 1000
2. **Static binary**: No runtime dependencies
3. **Minimal base image**: Alpine Linux (~5MB)
4. **Health checks**: Automatic container restart on failure
5. **Environment secrets**: Use Docker secrets or env files (not committed)
6. **SSL/TLS**: Enable DB_SSL_MODE=require in production
7. **Resource limits**: Set in docker-compose.prod.yml

## Monitoring

### Container Stats
```bash
docker stats buildboard-api
```

### Health Status
```bash
curl http://localhost:8080/health
```

### Application Logs
```bash
docker-compose logs -f --tail=100 api
```

## Updating the Application

```bash
# Pull latest code
git pull

# Rebuild and restart
docker-compose up -d --build

# Or with make
make docker-build
make docker-restart
```

## Multi-Service Deployment (with Frontend)

To deploy frontend + backend together, see the root docker-compose.yml (if created), or deploy separately:

```bash
# Backend
cd backend && make docker-up

# Frontend (in another terminal)
cd frontend && docker-compose up -d
```

## Support

For issues or questions:
- Check logs: `make docker-logs`
- Health check: `curl http://localhost:8080/health`
- Database status: `make docker-ps`