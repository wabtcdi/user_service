#!/bin/bash
# PostgreSQL Deployment Script for Minikube
# Date: January 17, 2026

echo "=== PostgreSQL Deployment Script ==="
echo ""

# Check if minikube is running
echo "Step 1: Checking Minikube status..."
if ! minikube status | grep -q "Running"; then
    echo "✗ Minikube is not running"
    echo "Please start Minikube first: minikube start"
    exit 1
fi
echo "✓ Minikube is running"
echo ""

# Get Minikube IP
MINIKUBE_IP=$(minikube ip)
echo "Minikube IP: $MINIKUBE_IP"
echo ""

# Apply Kubernetes manifests
echo "Step 2: Deploying PostgreSQL to Kubernetes..."
echo ""

echo "Creating Secret..."
kubectl apply -f k8s/postgres-secret.yaml
echo ""

echo "Creating PersistentVolumeClaim..."
kubectl apply -f k8s/postgres-pvc.yaml
echo ""

echo "Creating Deployment..."
kubectl apply -f k8s/postgres-deployment.yaml
echo ""

echo "Creating Service..."
kubectl apply -f k8s/postgres-service.yaml
echo ""

# Wait for deployment
echo "Step 3: Waiting for PostgreSQL pod to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres --timeout=300s

if [ $? -eq 0 ]; then
    echo "✓ PostgreSQL pod is ready"
else
    echo "✗ PostgreSQL pod failed to become ready"
    echo ""
    echo "Checking pod status:"
    kubectl get pods -l app=postgres
    echo ""
    echo "Pod logs:"
    kubectl logs -l app=postgres --tail=50
    exit 1
fi
echo ""

# Display deployment information
echo "=== Deployment Complete! ==="
echo ""
echo "PostgreSQL Resources:"
kubectl get all -l app=postgres
echo ""
echo "Persistent Volume Claim:"
kubectl get pvc postgres-pvc
echo ""

echo "Connection Information:"
echo "  Host: $MINIKUBE_IP"
echo "  Port: 30432"
echo "  Database: userdb"
echo "  User: postgres"
echo "  Password: postgres123"
echo ""
echo "Connection string:"
echo "  postgresql://postgres:postgres123@$MINIKUBE_IP:30432/userdb"
echo ""
echo "Test connection with:"
echo "  psql -h $MINIKUBE_IP -p 30432 -U postgres -d userdb"
echo ""
echo "Useful commands:"
echo "  kubectl get pods -l app=postgres           - Check pod status"
echo "  kubectl logs -l app=postgres               - View logs"
echo "  kubectl describe pod -l app=postgres       - Detailed pod info"
echo "  kubectl exec -it <pod-name> -- psql -U postgres -d userdb  - Access database"
echo ""
