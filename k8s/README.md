# Kubernetes Deployment Guide

## Prerequisites
- VirtualBox installed
- kubectl installed
- Minikube executable

## Quick Start

### 1. Install and Start Minikube
```bash
chmod +x install-minikube.sh
./install-minikube.sh
```

This script will:
- Verify VirtualBox installation
- Clean up old Minikube installations
- Install Minikube to ~/bin/minikube.exe
- Start Minikube with VirtualBox driver
- Enable useful addons (metrics-server, dashboard)

### 2. Deploy PostgreSQL
```bash
chmod +x deploy-postgres.sh
./deploy-postgres.sh
```

This script will:
- Deploy PostgreSQL 15 with persistent storage
- Create a NodePort service on port 30432
- Wait for the pod to be ready
- Display connection information

## Manual Installation

### Install Minikube
```bash
# Clean up old installation
minikube delete --all --purge
rm -rf ~/.minikube ~/.kube

# Install executable
mkdir -p ~/bin
cp minikube-windows-amd64.exe ~/bin/minikube.exe
chmod +x ~/bin/minikube.exe
export PATH="$HOME/bin:$PATH"

# Start Minikube
minikube start --driver=virtualbox --cpus=2 --memory=4096
```

### Deploy PostgreSQL
```bash
# Apply manifests
kubectl apply -f k8s/postgres-secret.yaml
kubectl apply -f k8s/postgres-pvc.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml

# Check status
kubectl get pods -l app=postgres
kubectl get svc postgres-service
```

## Configuration Files

### PostgreSQL Resources
- `postgres-secret.yaml` - Database credentials
- `postgres-pvc.yaml` - Persistent volume claim (1Gi)
- `postgres-deployment.yaml` - PostgreSQL deployment (1 replica)
- `postgres-service.yaml` - NodePort service (port 30432)

### Database Configuration
- **Image**: postgres:15-alpine
- **Database**: userdb
- **User**: postgres
- **Password**: postgres123
- **Port**: 5432 (internal), 30432 (external NodePort)

## Connection

### Get Minikube IP
```bash
minikube ip
```

### Connect via psql
```bash
psql -h $(minikube ip) -p 30432 -U postgres -d userdb
```

### Connection String
```
postgresql://postgres:postgres123@$(minikube ip):30432/userdb
```

## Useful Commands

### Minikube Management
```bash
minikube status              # Check cluster status
minikube stop                # Stop the cluster
minikube start               # Start the cluster
minikube delete              # Delete the cluster
minikube dashboard           # Open Kubernetes dashboard
minikube ip                  # Get cluster IP address
minikube ssh                 # SSH into the minikube VM
minikube logs                # View minikube logs
```

### Kubernetes Management
```bash
# Pods
kubectl get pods -l app=postgres
kubectl describe pod -l app=postgres
kubectl logs -l app=postgres
kubectl logs -l app=postgres -f          # Follow logs

# Services
kubectl get svc postgres-service
kubectl describe svc postgres-service

# Storage
kubectl get pvc postgres-pvc
kubectl describe pvc postgres-pvc

# Secrets
kubectl get secret postgres-secret
kubectl describe secret postgres-secret

# All resources
kubectl get all -l app=postgres
```

### Database Access
```bash
# Get pod name
POD_NAME=$(kubectl get pod -l app=postgres -o jsonpath="{.items[0].metadata.name}")

# Execute psql in the pod
kubectl exec -it $POD_NAME -- psql -U postgres -d userdb

# Run SQL command
kubectl exec -it $POD_NAME -- psql -U postgres -d userdb -c "SELECT version();"

# Copy SQL file and execute
kubectl cp database/schema.sql $POD_NAME:/tmp/schema.sql
kubectl exec -it $POD_NAME -- psql -U postgres -d userdb -f /tmp/schema.sql
```

## Troubleshooting

### Minikube won't start
```bash
# Check VirtualBox
vboxmanage --version

# Check logs
minikube logs

# Try deleting and starting fresh
minikube delete
minikube start --driver=virtualbox
```

### PostgreSQL pod not starting
```bash
# Check pod status
kubectl get pods -l app=postgres

# View pod events
kubectl describe pod -l app=postgres

# Check logs
kubectl logs -l app=postgres

# Check PVC status
kubectl get pvc postgres-pvc
```

### Cannot connect to database
```bash
# Verify service is running
kubectl get svc postgres-service

# Check if pod is ready
kubectl get pods -l app=postgres

# Test from within the cluster
kubectl run -it --rm debug --image=postgres:15-alpine --restart=Never -- psql -h postgres-service -U postgres -d userdb

# Check NodePort
kubectl get svc postgres-service -o jsonpath='{.spec.ports[0].nodePort}'
```

### Certificate issues
If you encounter TLS certificate errors:
```bash
minikube delete
rm -rf ~/.minikube
minikube start --driver=virtualbox
```

## Addons

### Enable useful addons
```bash
minikube addons enable metrics-server
minikube addons enable dashboard
minikube addons enable ingress
```

### View enabled addons
```bash
minikube addons list
```

## Performance Tuning

### Adjust resources
```bash
minikube start --driver=virtualbox --cpus=4 --memory=8192 --disk-size=20g
```

### PostgreSQL resource limits
Edit `k8s/postgres-deployment.yaml` to adjust:
- `resources.requests.memory`
- `resources.requests.cpu`
- `resources.limits.memory`
- `resources.limits.cpu`

## Cleanup

### Delete PostgreSQL deployment
```bash
kubectl delete -f k8s/postgres-service.yaml
kubectl delete -f k8s/postgres-deployment.yaml
kubectl delete -f k8s/postgres-pvc.yaml
kubectl delete -f k8s/postgres-secret.yaml
```

### Stop Minikube
```bash
minikube stop
```

### Delete Minikube cluster
```bash
minikube delete
```

### Complete cleanup
```bash
minikube delete --all --purge
rm -rf ~/.minikube ~/.kube
```
