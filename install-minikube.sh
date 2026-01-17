#!/bin/bash
# Minikube Installation and Setup Script for Windows with VirtualBox
# Date: January 17, 2026

echo "=== Minikube Installation and Setup Script ==="
echo ""

# Step 1: Verify VirtualBox is installed
echo "Step 1: Checking VirtualBox installation..."
if command -v vboxmanage &> /dev/null; then
    echo "✓ VirtualBox is installed:"
    vboxmanage --version
else
    echo "✗ VirtualBox is NOT installed"
    echo "Please install VirtualBox from: https://www.virtualbox.org/wiki/Downloads"
    echo "After installation, run this script again."
    exit 1
fi
echo ""

# Step 2: Remove old Minikube installation and config
echo "Step 2: Cleaning up old Minikube installation..."
minikube delete --all --purge 2>/dev/null || true
rm -rf ~/.minikube
rm -rf ~/.kube
echo "✓ Cleanup complete"
echo ""

# Step 3: Install Minikube executable
echo "Step 3: Installing Minikube..."
if [ -f "minikube-windows-amd64.exe" ]; then
    echo "✓ Minikube executable found in current directory"

    # Create a bin directory in user home if it doesn't exist
    mkdir -p ~/bin

    # Copy and rename the executable
    cp minikube-windows-amd64.exe ~/bin/minikube.exe
    chmod +x ~/bin/minikube.exe

    echo "✓ Minikube installed to ~/bin/minikube.exe"

    # Add to PATH (temporary for this session)
    export PATH="$HOME/bin:$PATH"

    echo ""
    echo "Note: Add the following line to your ~/.bashrc to make it permanent:"
    echo 'export PATH="$HOME/bin:$PATH"'
else
    echo "✗ minikube-windows-amd64.exe not found"
    echo "Downloading Minikube..."
    curl -Lo minikube-windows-amd64.exe https://storage.googleapis.com/minikube/releases/latest/minikube-windows-amd64.exe

    if [ -f "minikube-windows-amd64.exe" ]; then
        mkdir -p ~/bin
        cp minikube-windows-amd64.exe ~/bin/minikube.exe
        chmod +x ~/bin/minikube.exe
        export PATH="$HOME/bin:$PATH"
        echo "✓ Minikube downloaded and installed"
    else
        echo "✗ Failed to download Minikube"
        exit 1
    fi
fi
echo ""

# Step 4: Verify Minikube installation
echo "Step 4: Verifying Minikube installation..."
~/bin/minikube.exe version || minikube version
echo "✓ Minikube is accessible"
echo ""

# Step 5: Start Minikube with VirtualBox driver
echo "Step 5: Starting Minikube with VirtualBox driver..."
echo "This may take several minutes on first run..."
~/bin/minikube.exe start --driver=virtualbox --cpus=2 --memory=4096 || minikube start --driver=virtualbox --cpus=2 --memory=4096

if [ $? -eq 0 ]; then
    echo "✓ Minikube started successfully"
else
    echo "✗ Failed to start Minikube"
    echo "Check the logs with: minikube logs"
    exit 1
fi
echo ""

# Step 6: Verify the cluster
echo "Step 6: Verifying cluster status..."
echo ""
echo "Minikube status:"
minikube status
echo ""

echo "Cluster info:"
kubectl cluster-info
echo ""

echo "Nodes:"
kubectl get nodes
echo ""

# Step 7: Get Minikube IP
echo "Step 7: Getting Minikube IP address..."
MINIKUBE_IP=$(minikube ip)
echo "✓ Minikube IP: $MINIKUBE_IP"
echo ""

# Step 8: Enable required addons
echo "Step 8: Enabling useful addons..."
minikube addons enable metrics-server
minikube addons enable dashboard
echo "✓ Addons enabled"
echo ""

echo "=== Installation Complete! ==="
echo ""
echo "Summary:"
echo "  - Minikube IP: $MINIKUBE_IP"
echo "  - Driver: VirtualBox"
echo ""
echo "Useful commands:"
echo "  minikube status          - Check cluster status"
echo "  minikube stop            - Stop the cluster"
echo "  minikube start           - Start the cluster"
echo "  minikube dashboard       - Open Kubernetes dashboard"
echo "  minikube ip              - Get cluster IP"
echo "  kubectl get pods -A      - List all pods"
echo ""
echo "Next steps:"
echo "  1. Deploy PostgreSQL: ./deploy-postgres.sh"
echo "  2. Check deployment: kubectl get pods"
echo "  3. Connect to PostgreSQL at: $MINIKUBE_IP:30432"
echo ""
