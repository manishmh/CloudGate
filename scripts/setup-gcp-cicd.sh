#!/bin/bash

# CloudGate CI/CD Setup Script for Google Cloud Platform
# This script creates a service account with proper permissions for GitHub Actions

set -e

# Configuration
PROJECT_ID="routemate-409518"
SERVICE_ACCOUNT_NAME="cloudgate-cicd"
SERVICE_ACCOUNT_EMAIL="${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
KEY_FILE="cloudgate-cicd-key.json"
REGION="us-central1"
REPO_NAME="cloudgate-repo"

echo "🚀 Setting up CloudGate CI/CD for Google Cloud Platform"
echo "Project ID: $PROJECT_ID"
echo "Service Account: $SERVICE_ACCOUNT_EMAIL"
echo ""

# Check if gcloud is installed and authenticated
if ! command -v gcloud &> /dev/null; then
    echo "❌ Error: gcloud CLI is not installed"
    echo "Please install it from: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

# Set the project
echo "📋 Setting project to $PROJECT_ID..."
gcloud config set project $PROJECT_ID

# Enable required APIs
echo "🔧 Enabling required APIs..."
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable artifactregistry.googleapis.com
gcloud services enable iam.googleapis.com

# Create service account if it doesn't exist
echo "👤 Creating service account..."
if gcloud iam service-accounts describe $SERVICE_ACCOUNT_EMAIL &>/dev/null; then
    echo "Service account already exists"
else
    gcloud iam service-accounts create $SERVICE_ACCOUNT_NAME \
        --display-name="CloudGate CI/CD Service Account" \
        --description="Service account for CloudGate GitHub Actions CI/CD pipeline"
fi

# Grant necessary roles to the service account
echo "🔐 Granting permissions to service account..."

# Artifact Registry permissions
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
    --role="roles/artifactregistry.admin"

# Cloud Run permissions
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
    --role="roles/run.admin"

# Cloud Build permissions (for building images)
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
    --role="roles/cloudbuild.builds.editor"

# Storage permissions (for Cloud Build artifacts)
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
    --role="roles/storage.admin"

# IAM permissions (to manage service accounts for Cloud Run)
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$SERVICE_ACCOUNT_EMAIL" \
    --role="roles/iam.serviceAccountUser"

# Create Artifact Registry repository if it doesn't exist
echo "📦 Creating Artifact Registry repository..."
if gcloud artifacts repositories describe $REPO_NAME --location=$REGION &>/dev/null; then
    echo "Artifact Registry repository already exists"
else
    gcloud artifacts repositories create $REPO_NAME \
        --repository-format=docker \
        --location=$REGION \
        --description="CloudGate application images"
fi

# Generate service account key
echo "🔑 Generating service account key..."
if [ -f "$KEY_FILE" ]; then
    echo "Key file already exists. Removing old key..."
    rm "$KEY_FILE"
fi

gcloud iam service-accounts keys create $KEY_FILE \
    --iam-account=$SERVICE_ACCOUNT_EMAIL

echo ""
echo "✅ Setup completed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Add the following secret to your GitHub repository:"
echo "   Secret name: GCP_SA_KEY"
echo "   Secret value: Copy the entire contents of $KEY_FILE"
echo ""
echo "2. To copy the key content, run:"
echo "   cat $KEY_FILE"
echo ""
echo "3. Go to your GitHub repository settings:"
echo "   https://github.com/YOUR_USERNAME/CloudGate/settings/secrets/actions"
echo ""
echo "4. Click 'New repository secret'"
echo "5. Name: GCP_SA_KEY"
echo "6. Value: Paste the entire JSON content from $KEY_FILE"
echo ""
echo "⚠️  IMPORTANT: Keep the $KEY_FILE secure and delete it after adding to GitHub!"
echo "   rm $KEY_FILE"
echo ""
echo "🔧 Artifact Registry repository created:"
echo "   $REGION-docker.pkg.dev/$PROJECT_ID/$REPO_NAME"
echo ""
echo "🚀 Your CI/CD pipeline should now work!" 