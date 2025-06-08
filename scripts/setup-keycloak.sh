#!/bin/bash

# CloudGate Keycloak Setup Script
# This script sets up the Keycloak realm and client for CloudGate SSO

set -e

KEYCLOAK_URL="http://localhost:8080"
ADMIN_USER="admin"
ADMIN_PASS="admin"
REALM_NAME="cloudgate"
CLIENT_ID="cloudgate-frontend"

echo "🚀 Setting up CloudGate Keycloak configuration..."

# Wait for Keycloak to be ready
echo "⏳ Waiting for Keycloak to be ready..."
until curl -f ${KEYCLOAK_URL} > /dev/null 2>&1; do
    echo "Waiting for Keycloak..."
    sleep 2
done

echo "✅ Keycloak is ready!"

# Get admin token
echo "🔑 Getting admin token..."
ADMIN_TOKEN=$(curl -s -X POST "${KEYCLOAK_URL}/realms/master/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=${ADMIN_USER}" \
  -d "password=${ADMIN_PASS}" \
  -d "grant_type=password" \
  -d "client_id=admin-cli" | jq -r '.access_token')

if [ "$ADMIN_TOKEN" = "null" ] || [ -z "$ADMIN_TOKEN" ]; then
    echo "❌ Failed to get admin token. Please check your admin credentials."
    echo "💡 Go to http://localhost:8080/admin/ and create an admin user first."
    exit 1
fi

echo "✅ Admin token obtained!"

# Create realm
echo "🏗️  Creating CloudGate realm..."
curl -s -X POST "${KEYCLOAK_URL}/admin/realms" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "realm": "'${REALM_NAME}'",
    "enabled": true,
    "displayName": "CloudGate SSO",
    "registrationAllowed": true,
    "loginWithEmailAllowed": true,
    "duplicateEmailsAllowed": false,
    "resetPasswordAllowed": true,
    "editUsernameAllowed": false,
    "bruteForceProtected": true
  }' || echo "Realm might already exist"

echo "✅ CloudGate realm created!"

# Create client
echo "🔧 Creating frontend client..."
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/clients" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "clientId": "'${CLIENT_ID}'",
    "name": "CloudGate Frontend",
    "description": "CloudGate SSO Frontend Application",
    "enabled": true,
    "clientAuthenticatorType": "client-secret",
    "publicClient": true,
    "standardFlowEnabled": true,
    "directAccessGrantsEnabled": true,
    "serviceAccountsEnabled": false,
    "authorizationServicesEnabled": false,
    "redirectUris": ["http://localhost:3000/*"],
    "webOrigins": ["http://localhost:3000"],
    "attributes": {
      "pkce.code.challenge.method": "S256"
    }
  }' || echo "Client might already exist"

echo "✅ Frontend client created!"

# Create test user
echo "👤 Creating test user..."
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/users" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@cloudgate.com",
    "firstName": "Test",
    "lastName": "User",
    "enabled": true,
    "emailVerified": true
  }' || echo "User might already exist"

# Get user ID
USER_ID=$(curl -s -X GET "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/users?username=testuser" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" | jq -r '.[0].id')

# Set user password
echo "🔒 Setting user password..."
curl -s -X PUT "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/users/${USER_ID}/reset-password" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "password",
    "value": "password123",
    "temporary": false
  }'

echo "✅ Test user created!"

echo ""
echo "🎉 CloudGate Keycloak setup complete!"
echo ""
echo "📋 Configuration Summary:"
echo "   Keycloak URL: ${KEYCLOAK_URL}"
echo "   Realm: ${REALM_NAME}"
echo "   Client ID: ${CLIENT_ID}"
echo "   Test User: testuser / password123"
echo ""
echo "🌐 Access points:"
echo "   Admin Console: ${KEYCLOAK_URL}/admin/"
echo "   CloudGate Frontend: http://localhost:3000"
echo ""
echo "🔄 Now restart your frontend to connect to the configured realm!" 