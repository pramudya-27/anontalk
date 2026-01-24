# Kubernetes Deployment Guide - AnonTalk

## Prerequisites

- Minikube installed and running
- kubectl configured
- Docker images pushed to Docker Hub:
  - `dio713/anontalk-backend:latest`
  - `dio713/anontalk-frontend:latest`

## Quick Start

### 1. Start Minikube

```bash
minikube start
```

### 2. Deploy All Resources

```bash
cd k8s
kubectl apply -f .
```

This will create all resources in the correct order:

- `00-namespace.yaml` - Creates `anontalk` namespace
- `01-configmap.yaml` - Application configuration
- `02-secret.yaml` - Database credentials
- `03-mysql.yaml` - MySQL database with persistent storage
- `04-redis.yaml` - Redis cache
- `05-backend-api.yaml` - Backend API service
- `06-matchmaker.yaml` - Matchmaking worker
- `07-frontend.yaml` - Frontend web app

### 3. Check Deployment Status

```bash
# Check all pods
kubectl get pods -n anontalk

# Check services
kubectl get svc -n anontalk

# Watch pod status
kubectl get pods -n anontalk -w
```

### 4. Access the Application

#### Frontend (NodePort)

```bash
# Get the URL
minikube service frontend -n anontalk --url

# Or access via NodePort
kubectl get svc frontend -n anontalk
# Then access: http://<minikube-ip>:30000
```

#### Backend API (Port Forward)

```bash
kubectl port-forward -n anontalk svc/backend-api 8080:8080
# Access at: http://localhost:8080
```

### 5. Initialize Database

The MySQL database will be created automatically. If you need to run migrations:

```bash
# Get MySQL pod name
kubectl get pods -n anontalk | grep mysql

# Copy migration files
kubectl cp ../backend/migrations <mysql-pod-name>:/tmp/migrations -n anontalk

# Execute migrations (if needed)
kubectl exec -it <mysql-pod-name> -n anontalk -- mysql -uroot -prootpassword anontalk < /tmp/migrations/001_init.sql
```

## Troubleshooting

### Check Logs

```bash
# Backend API logs
kubectl logs -n anontalk -l app=backend-api

# Matchmaker logs
kubectl logs -n anontalk -l app=matchmaker

# Frontend logs
kubectl logs -n anontalk -l app=frontend

# MySQL logs
kubectl logs -n anontalk -l app=mysql

# Redis logs
kubectl logs -n anontalk -l app=redis
```

### Restart Deployments

```bash
kubectl rollout restart deployment/backend-api -n anontalk
kubectl rollout restart deployment/matchmaker -n anontalk
kubectl rollout restart deployment/frontend -n anontalk
```

### Delete and Redeploy

```bash
# Delete all resources
kubectl delete -f . -n anontalk

# Or delete namespace (removes everything)
kubectl delete namespace anontalk

# Redeploy
kubectl apply -f .
```

## Configuration

### Update Environment Variables

Edit `01-configmap.yaml` or `02-secret.yaml` and reapply:

```bash
kubectl apply -f 01-configmap.yaml
kubectl apply -f 02-secret.yaml

# Restart affected pods
kubectl rollout restart deployment/backend-api -n anontalk
kubectl rollout restart deployment/matchmaker -n anontalk
```

### Scale Services

```bash
# Scale backend API
kubectl scale deployment/backend-api --replicas=2 -n anontalk

# Scale matchmaker
kubectl scale deployment/matchmaker --replicas=2 -n anontalk
```

## Architecture

```
┌─────────────────────────────────────────────────┐
│              Namespace: anontalk                │
├─────────────────────────────────────────────────┤
│                                                 │
│  ┌──────────┐      ┌──────────────┐            │
│  │ Frontend │─────▶│ Backend API  │            │
│  │ (NodePort│      │ (ClusterIP)  │            │
│  │  :30000) │      │   :8080      │            │
│  └──────────┘      └───────┬──────┘            │
│                            │                    │
│                    ┌───────┴────────┐           │
│                    │                │           │
│              ┌─────▼─────┐   ┌─────▼─────┐     │
│              │   MySQL   │   │   Redis   │     │
│              │  (PVC 1Gi)│   │           │     │
│              └───────────┘   └─────┬─────┘     │
│                                    │           │
│                              ┌─────▼─────┐     │
│                              │Matchmaker │     │
│                              │  Worker   │     │
│                              └───────────┘     │
└─────────────────────────────────────────────────┘
```

## Notes

- **Database**: MySQL data persists via PVC even if pod restarts
- **Redis**: In-memory only, data lost on restart
- **WebSocket**: Backend handles WebSocket connections for real-time chat
- **Matchmaking**: Worker continuously processes the Redis queue
- **Images**: Using public Docker Hub images (dio713/anontalk-\*)
