# ChargeSphere Deployment Guide

This guide covers deploying the OCPI service to DigitalOcean App Platform.

## Deployment Overview

**Stack:**
- **Application**: Go OCPI Service (Docker container)
- **Database**: DigitalOcean Managed MongoDB
- **Cache**: DigitalOcean Managed Redis
- **Platform**: DigitalOcean App Platform

**Estimated Cost:**
- Basic: ~$78/month (2 instances, 1GB MongoDB, 1GB Redis)
- Production: ~$291/month (4 instances, 8GB MongoDB, 4GB Redis)

---

## Prerequisites

1. **DigitalOcean Account**
   - Sign up at https://www.digitalocean.com
   - Add payment method

2. **GitHub Repository**
   - Fork or push code to GitHub
   - Repository: `ikkurthis1998/charge-sphere`

3. **doctl CLI** (Optional)
   ```bash
   # Install doctl
   brew install doctl  # macOS
   # or
   snap install doctl  # Linux

   # Authenticate
   doctl auth init
   ```

---

## Deployment Methods

### Method 1: Using DigitalOcean Dashboard (Easiest)

#### Step 1: Create App

1. Go to https://cloud.digitalocean.com/apps
2. Click **"Create App"**
3. Select **"GitHub"** as source
4. Authorize DigitalOcean to access your repository
5. Select repository: `ikkurthis1998/charge-sphere`
6. Select branch: `main`
7. Click **"Next"**

#### Step 2: Configure Service

1. **Detect Resources**: App Platform should auto-detect the Dockerfile
2. **Source Directory**: Set to `/services/ocpi-service`
3. **Dockerfile Path**: Set to `services/ocpi-service/Dockerfile`
4. **HTTP Port**: 8080
5. **Health Check Path**: `/health`
6. Click **"Next"**

#### Step 3: Add Managed Databases

**MongoDB:**
1. Click **"Add Resource"** → **"Database"**
2. Engine: **MongoDB 7**
3. Plan: **Basic 1GB RAM** ($15/month)
4. Data center: **New York 3**
5. Name: `chargesphere-db`

**Redis:**
1. Click **"Add Resource"** → **"Database"**
2. Engine: **Redis 7**
3. Plan: **Basic 1GB RAM** ($15/month)
4. Data center: **New York 3**
5. Name: `chargesphere-redis`

#### Step 4: Configure Environment Variables

Add the following environment variables:

```
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_MODE=release

MONGODB_URI=${chargesphere-db.DATABASE_URL}
MONGODB_DATABASE=chargesphere
MONGODB_TIMEOUT=10

REDIS_HOST=${chargesphere-redis.HOSTNAME}:${chargesphere-redis.PORT}
REDIS_PASSWORD=${chargesphere-redis.PASSWORD}
REDIS_DB=0

TEMPORAL_HOST=localhost:7233
TEMPORAL_NAMESPACE=default

OCPI_VERSION=2.3
OCPI_COUNTRY_CODE=US
OCPI_PARTY_ID=CSP
```

**Note:** DigitalOcean automatically injects database connection strings.

#### Step 5: Configure Resources

1. **Instance Count**: 2 (for high availability)
2. **Instance Size**: Basic ($6/month each)
3. **Region**: New York 3
4. Click **"Next"**

#### Step 6: Review and Deploy

1. Review all settings
2. **Estimated Cost**: ~$78/month
3. Click **"Create Resources"**
4. Wait 5-10 minutes for deployment

#### Step 7: Verify Deployment

Once deployed, you'll get a URL like: `https://charge-sphere-xxxxx.ondigitalocean.app`

Test it:
```bash
# Health check
curl https://charge-sphere-xxxxx.ondigitalocean.app/health

# Should return:
# {"status":"healthy","service":"ocpi-service"}
```

---

### Method 2: Using App Spec (app.yaml)

#### Step 1: Prepare App Spec

The app spec is already created at `.do/app.yaml`. Review and modify if needed.

#### Step 2: Deploy via doctl

```bash
# Create app
doctl apps create --spec .do/app.yaml

# Get app ID
doctl apps list

# Check deployment status
doctl apps get <app-id>

# View logs
doctl apps logs <app-id> --type run
```

#### Step 3: Update App Spec

If you need to update configuration:

```bash
# Update app
doctl apps update <app-id> --spec .do/app.yaml
```

---

### Method 3: Deploy via GitHub Integration

#### Step 1: Enable Auto-Deploy

In the DigitalOcean dashboard:
1. Go to your app
2. Settings → App Spec
3. Enable **"Autodeploy"**
4. Select branch: `main`

#### Step 2: Push to GitHub

```bash
git add .
git commit -m "Deploy to production"
git push origin main
```

The app will automatically deploy on every push to `main`.

---

## Post-Deployment Configuration

### 1. Set Up Custom Domain (Optional)

1. Go to App → Settings → Domains
2. Add your domain (e.g., `api.chargesphere.com`)
3. Add DNS records:
   ```
   Type: CNAME
   Name: api
   Value: charge-sphere-xxxxx.ondigitalocean.app
   ```
4. DigitalOcean will automatically provision SSL certificate

### 2. Configure Firewall

1. Go to Networking → Firewalls
2. Create firewall rule
3. Allow:
   - HTTP (80)
   - HTTPS (443)
4. Attach to app

### 3. Set Up Monitoring

DigitalOcean provides built-in monitoring:
- CPU usage
- Memory usage
- Network traffic
- Response times

View in: App → Insights

### 4. Configure Alerts

1. Go to App → Settings → Alerts
2. Add alert for:
   - High error rate (> 5%)
   - High response time (> 1s)
   - Low memory (< 100MB available)

---

## Database Setup

### MongoDB Indexes

The application automatically creates indexes on startup. Verify:

```javascript
// Connect to MongoDB
use chargesphere

// Check indexes
db.partners.getIndexes()

// Should show:
// - _id (default)
// - partner_id (unique)
// - credentials.token (unique)
// - status
```

### Redis Configuration

Redis is used for caching and rate limiting. No special setup needed.

---

## Scaling

### Horizontal Scaling

Increase instance count:

```bash
# Via doctl
doctl apps update <app-id> --instance-count 4

# Or via dashboard
App → Settings → Resources → Instance Count: 4
```

### Vertical Scaling

Upgrade instance size:

```bash
# Via dashboard
App → Settings → Resources → Instance Size: Professional (2GB RAM)
```

### Database Scaling

Upgrade database plan:
1. Go to Databases → chargesphere-db
2. Settings → Resize
3. Select larger plan

**Recommended for production:**
- MongoDB: 4GB RAM ($90/month)
- Redis: 2GB RAM ($30/month)

---

## Environment-Specific Deployments

### Development

```yaml
instance_count: 1
instance_size_slug: basic-xxs  # $5/month
databases:
  size: db-s-1vcpu-1gb  # $15/month
```

### Staging

```yaml
instance_count: 2
instance_size_slug: basic-xs  # $6/month
databases:
  size: db-s-1vcpu-2gb  # $30/month
```

### Production

```yaml
instance_count: 4
instance_size_slug: professional-xs  # $12/month
databases:
  size: db-s-2vcpu-4gb  # $90/month
```

---

## CI/CD Pipeline

### GitHub Actions

A GitHub Actions workflow is provided at `.github/workflows/deploy.yml`.

**Triggers:**
- Push to `main` → Deploy to production
- Push to `develop` → Deploy to staging
- Pull request → Run tests only

**Steps:**
1. Checkout code
2. Run unit tests
3. Run integration tests
4. Build Docker image
5. Deploy to DigitalOcean

---

## Monitoring & Logs

### View Logs

```bash
# Via doctl
doctl apps logs <app-id> --type run --follow

# Or via dashboard
App → Runtime Logs
```

### Monitor Performance

DigitalOcean provides:
- Request rate
- Response time (p50, p95, p99)
- Error rate
- CPU usage
- Memory usage

Access: App → Insights

### Set Up External Monitoring (Optional)

**Recommended tools:**
- **Uptime monitoring**: UptimeRobot, Pingdom
- **APM**: New Relic, Datadog
- **Error tracking**: Sentry

---

## Backup & Recovery

### Database Backups

DigitalOcean automatically backs up managed databases:
- **Daily backups** (retained for 7 days)
- **Point-in-time recovery** (last 7 days)

Restore:
1. Go to Databases → chargesphere-db
2. Backups & Restore
3. Select backup
4. Restore to new database or overwrite

### Application State

Application is stateless - no backup needed.

---

## Security Best Practices

### 1. Enable HTTPS Only

```yaml
# In app.yaml
routes:
  - path: /
    preserve_path_prefix: true
    protocol: HTTPS
```

### 2. Use Secrets for Sensitive Data

Don't commit secrets to git. Use environment variables:

```bash
# Set secret via doctl
doctl apps update <app-id> \
  --env MONGODB_URI=secret:mongodb-uri
```

### 3. Enable Database Trusted Sources

1. Go to Databases → chargesphere-db
2. Settings → Trusted Sources
3. Add only: App Platform IP ranges

### 4. Regular Updates

Keep dependencies updated:

```bash
# Update Go dependencies
go get -u ./...
go mod tidy

# Rebuild Docker image
git commit -am "Update dependencies"
git push
```

---

## Troubleshooting

### Deployment Failed

**Check build logs:**
```bash
doctl apps logs <app-id> --type build
```

**Common issues:**
- Missing Dockerfile
- Incorrect source directory
- Build command failed

### Application Not Starting

**Check runtime logs:**
```bash
doctl apps logs <app-id> --type run
```

**Common issues:**
- Cannot connect to MongoDB
- Config file missing
- Port already in use

### Database Connection Error

1. Verify database is running:
   ```bash
   doctl databases list
   ```

2. Check connection string:
   ```bash
   doctl databases connection <db-id>
   ```

3. Verify trusted sources include App Platform

### High Memory Usage

1. Check metrics in App → Insights
2. Consider increasing instance size
3. Review code for memory leaks

### Slow Response Times

1. Check database query performance
2. Add database indexes if needed
3. Enable Redis caching
4. Consider horizontal scaling

---

## Cost Optimization

### Development/Testing

- Use 1 instance (basic-xxs): $5/month
- Use smaller database (1GB): $30/month
- Pause when not in use
- **Total: ~$35/month**

### Production

- Use 2-4 instances for high availability
- Use appropriate database size (4-8GB)
- Enable auto-scaling
- **Total: $78-$291/month**

### Free Tier

DigitalOcean doesn't have a free tier, but offers:
- $200 credit for new accounts (60 days)
- Use this to test deployment

---

## Rollback

### Via Dashboard

1. Go to App → Deployments
2. Find previous successful deployment
3. Click **"Rollback to this deployment"**

### Via doctl

```bash
# List deployments
doctl apps list-deployments <app-id>

# Rollback to specific deployment
doctl apps rollback <app-id> <deployment-id>
```

---

## Health Checks

The application includes a health check endpoint:

```bash
curl https://your-app.ondigitalocean.app/health
```

Response:
```json
{
  "status": "healthy",
  "service": "ocpi-service"
}
```

DigitalOcean automatically:
- Checks health every 30 seconds
- Restarts unhealthy instances
- Routes traffic only to healthy instances

---

## Next Steps

After successful deployment:

1. ✅ Test all API endpoints
2. ✅ Register a test CPO partner
3. ✅ Verify database persistence
4. ✅ Set up monitoring alerts
5. ✅ Configure custom domain
6. ✅ Set up CI/CD pipeline
7. ✅ Document API for partners
8. ✅ Deploy Locations module (next)

---

## Support

**DigitalOcean:**
- Documentation: https://docs.digitalocean.com/products/app-platform/
- Support: https://cloud.digitalocean.com/support/tickets

**Project Issues:**
- GitHub: https://github.com/ikkurthis1998/charge-sphere/issues

---

## Quick Reference

```bash
# Deploy
doctl apps create --spec .do/app.yaml

# View logs
doctl apps logs <app-id> --follow

# Update app
doctl apps update <app-id> --spec .do/app.yaml

# Scale up
doctl apps update <app-id> --instance-count 4

# Rollback
doctl apps rollback <app-id> <deployment-id>

# Delete app
doctl apps delete <app-id>
```

---

## Deployment Checklist

- [ ] GitHub repository created
- [ ] Code pushed to `main` branch
- [ ] DigitalOcean account created
- [ ] App created in App Platform
- [ ] MongoDB database provisioned
- [ ] Redis database provisioned
- [ ] Environment variables configured
- [ ] Health check passing
- [ ] Test registration successful
- [ ] Custom domain configured (optional)
- [ ] Monitoring alerts set up
- [ ] CI/CD pipeline configured
- [ ] Documentation updated

---

**Your OCPI service is now live! 🚀**
