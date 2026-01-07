#!/bin/bash

# Deploy tc-fiap-50 Application to EKS
# This script deploys the application assuming the database is already deployed

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}[DEPLOY]${NC} $1"
}

# Function to check if kubectl is available
check_kubectl() {
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl is not installed or not in PATH"
        exit 1
    fi
    print_status "✅ kubectl is available"
}

# Function to check if we can connect to the cluster
check_cluster_connection() {
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to Kubernetes cluster"
        print_warning "Make sure you have the correct kubeconfig set up"
        exit 1
    fi
    print_status "✅ Connected to Kubernetes cluster"
}

# Function to check if database is available
check_database() {
    print_header "Checking database availability"
    
    # Check if using RDS (external database)
    if [ "$DATABASE_TYPE" = "rds" ]; then
        print_status "Using RDS database - skipping local database checks"
        return 0
    fi
    
    # Check if using local Kubernetes database
    if ! kubectl get service postgres-service &> /dev/null; then
        print_error "PostgreSQL service not found"
        print_warning "Please deploy the database first using tc-fiap-database project"
        exit 1
    fi
    
    if ! kubectl get pods -l app=postgres --field-selector=status.phase=Running &> /dev/null; then
        print_error "PostgreSQL pod is not running"
        print_warning "Please ensure the database is deployed and running"
        exit 1
    fi
    
    print_status "✅ Database is available"
}

# Function to deploy application secrets
deploy_secrets() {
    print_header "Deploying Application Secrets"
    
    # Check if app secrets already exist
    if kubectl get secret app-secret &> /dev/null; then
        print_warning "App secret already exists, skipping..."
        return 0
    fi
    
    # Create app secrets
    kubectl create secret generic app-secret \
        --from-literal=MERCADO_PAGO_BASEURL=https://api.mercadopago.com \
        --from-literal=MERCADO_PAGO_ACCESS_TOKEN=test_token \
        --from-literal=MERCADO_PAGO_CLIENT_ID=test_client \
        --from-literal=MERCADO_PAGO_POS_ID=test_pos \
        --from-literal=MERCADO_PAGO_WEBHOOK_SECRET=test_secret \
        --from-literal=MERCADO_PAGO_WEBHOOK_CALLBACK_URL=https://test.com/webhook \
        --dry-run=client -o yaml | kubectl apply -f -
    
    print_status "✅ Application secrets created"
}

# Function to deploy application
deploy_app() {
    print_header "Deploying tc-fiap-50 Application"
    
    kubectl apply -f app-deployment.yaml
    print_status "✅ Application deployment created"
    
    kubectl apply -f app-service.yaml
    print_status "✅ Application service created"
    
    kubectl apply -f app-hpa.yaml
    print_status "✅ Horizontal Pod Autoscaler created"
}

# Function to wait for application to be ready
wait_for_app() {
    print_header "Waiting for application to be ready"
    
    kubectl wait --for=condition=ready pod -l app=my-app --timeout=300s
    print_status "✅ Application is ready"
}

# Function to show deployment status
show_status() {
    print_header "Application Deployment Status"
    
    echo "Pods:"
    kubectl get pods -l app=my-app
    
    echo ""
    echo "Services:"
    kubectl get services -l app=my-app
    
    echo ""
    echo "HPA:"
    kubectl get hpa app-hpa
    
    # Get external URL
    EXTERNAL_IP=$(kubectl get service app-service -o jsonpath='{.status.loadBalancer.ingress[0].hostname}')
    if [ -n "$EXTERNAL_IP" ]; then
        echo ""
        print_status "🌐 Application is accessible at: http://${EXTERNAL_IP}:8080"
        print_status "📚 Swagger documentation: http://${EXTERNAL_IP}:8080/swagger/index.html"
    fi
}

# Main execution
main() {
    print_header "Starting tc-fiap-50 Application Deployment"
    
    # Pre-flight checks
    check_kubectl
    check_cluster_connection
    check_database
    
    # Deploy components
    deploy_secrets
    deploy_app
    wait_for_app
    
    # Show status
    show_status
    
    print_status "🎉 tc-fiap-50 application deployment completed successfully!"
}

# Run main function
main "$@"
