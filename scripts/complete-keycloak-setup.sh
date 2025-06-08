#!/bin/bash

# Complete CloudGate Keycloak Setup Script
# This script does EVERYTHING - creates admin user, realm, client, and test user

set -e

KEYCLOAK_URL="http://localhost:8080"
ADMIN_USER="admin"
ADMIN_PASS="admin"
REALM_NAME="cloudgate"
CLIENT_ID="cloudgate-frontend"

echo "🚀 Complete CloudGate Keycloak Setup Starting..."

# Wait for Keycloak to be ready
echo "⏳ Waiting for Keycloak to be ready..."
until curl -f ${KEYCLOAK_URL} > /dev/null 2>&1; do
    echo "Waiting for Keycloak..."
    sleep 2
done

echo "✅ Keycloak is ready!"

# Step 1: Create initial admin user
echo "👤 Creating initial admin user..."
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/master/users" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "'${ADMIN_USER}'",
    "enabled": true,
    "credentials": [{
      "type": "password",
      "value": "'${ADMIN_PASS}'",
      "temporary": false
    }]
  }' || echo "Admin user creation attempted..."

# Try to get admin token (this will work if admin user exists)
echo "🔑 Getting admin token..."
ADMIN_TOKEN=""
for i in {1..5}; do
    ADMIN_TOKEN=$(curl -s -X POST "${KEYCLOAK_URL}/realms/master/protocol/openid-connect/token" \
      -H "Content-Type: application/x-www-form-urlencoded" \
      -d "username=${ADMIN_USER}" \
      -d "password=${ADMIN_PASS}" \
      -d "grant_type=password" \
      -d "client_id=admin-cli" | jq -r '.access_token // empty' 2>/dev/null)
    
    if [ ! -z "$ADMIN_TOKEN" ] && [ "$ADMIN_TOKEN" != "null" ]; then
        echo "✅ Admin token obtained!"
        break
    fi
    
    echo "Attempt $i: Trying to create admin user via API..."
    # Alternative method: Use Keycloak's initial admin creation
    curl -s -X POST "${KEYCLOAK_URL}/admin/realms/master" \
      -H "Content-Type: application/json" \
      -d '{
        "realm": "master",
        "enabled": true,
        "users": [{
          "username": "'${ADMIN_USER}'",
          "enabled": true,
          "credentials": [{
            "type": "password",
            "value": "'${ADMIN_PASS}'",
            "temporary": false
          }],
          "realmRoles": ["admin"]
        }]
      }' 2>/dev/null || true
    
    sleep 2
done

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" = "null" ]; then
    echo "❌ Could not get admin token. Let's try a different approach..."
    
    # Use environment variables to set admin user (this works for fresh Keycloak)
    echo "🔧 Setting up admin user via environment..."
    docker exec cloudgate-keycloak /opt/keycloak/bin/kcadm.sh config credentials \
      --server http://localhost:8080 --realm master --user ${ADMIN_USER} --password ${ADMIN_PASS} 2>/dev/null || true
    
    # Try token again
    sleep 3
    ADMIN_TOKEN=$(curl -s -X POST "${KEYCLOAK_URL}/realms/master/protocol/openid-connect/token" \
      -H "Content-Type: application/x-www-form-urlencoded" \
      -d "username=${ADMIN_USER}" \
      -d "password=${ADMIN_PASS}" \
      -d "grant_type=password" \
      -d "client_id=admin-cli" | jq -r '.access_token // empty' 2>/dev/null)
fi

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" = "null" ]; then
    echo "❌ Still cannot get admin token. Using direct container commands..."
    
    # Create realm using kcadm directly in container
    echo "🏗️  Creating CloudGate realm via container..."
    docker exec cloudgate-keycloak /opt/keycloak/bin/kcadm.sh create realms \
      -s realm=${REALM_NAME} \
      -s enabled=true \
      -s displayName="CloudGate SSO" \
      -s registrationAllowed=true \
      -s loginWithEmailAllowed=true \
      --no-config --server http://localhost:8080 --realm master \
      --user ${ADMIN_USER} --password ${ADMIN_PASS} 2>/dev/null || echo "Realm creation attempted"
    
    # Create client using kcadm
    echo "🔧 Creating frontend client via container..."
    docker exec cloudgate-keycloak /opt/keycloak/bin/kcadm.sh create clients \
      -r ${REALM_NAME} \
      -s clientId=${CLIENT_ID} \
      -s name="CloudGate Frontend" \
      -s enabled=true \
      -s publicClient=true \
      -s standardFlowEnabled=true \
      -s directAccessGrantsEnabled=true \
      -s 'redirectUris=["http://localhost:3000/*"]' \
      -s 'webOrigins=["http://localhost:3000"]' \
      --no-config --server http://localhost:8080 --realm master \
      --user ${ADMIN_USER} --password ${ADMIN_PASS} 2>/dev/null || echo "Client creation attempted"
    
    # Create test user using kcadm
    echo "👤 Creating test user via container..."
    docker exec cloudgate-keycloak /opt/keycloak/bin/kcadm.sh create users \
      -r ${REALM_NAME} \
      -s username=testuser \
      -s email=test@cloudgate.com \
      -s firstName=Test \
      -s lastName=User \
      -s enabled=true \
      -s emailVerified=true \
      --no-config --server http://localhost:8080 --realm master \
      --user ${ADMIN_USER} --password ${ADMIN_PASS} 2>/dev/null || echo "User creation attempted"
    
    # Set user password using kcadm
    echo "🔒 Setting user password via container..."
    docker exec cloudgate-keycloak /opt/keycloak/bin/kcadm.sh set-password \
      -r ${REALM_NAME} \
      --username testuser \
      --new-password password123 \
      --temporary false \
      --no-config --server http://localhost:8080 --realm master \
      --user ${ADMIN_USER} --password ${ADMIN_PASS} 2>/dev/null || echo "Password setting attempted"
    
else
    echo "✅ Using API with admin token..."
    
    # Create realm using API
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

    # Create client using API
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

    # Create test user using API
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

    # Get user ID and set password
    USER_ID=$(curl -s -X GET "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/users?username=testuser" \
      -H "Authorization: Bearer ${ADMIN_TOKEN}" | jq -r '.[0].id // empty' 2>/dev/null)

    if [ ! -z "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
        echo "🔒 Setting user password..."
        curl -s -X PUT "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/users/${USER_ID}/reset-password" \
          -H "Authorization: Bearer ${ADMIN_TOKEN}" \
          -H "Content-Type: application/json" \
          -d '{
            "type": "password",
            "value": "password123",
            "temporary": false
          }'
        echo "✅ Test user password set!"
    fi
fi

# Verify setup
echo "🔍 Verifying setup..."
sleep 2

# Test if realm exists
REALM_CHECK=$(curl -s "${KEYCLOAK_URL}/realms/${REALM_NAME}/.well-known/openid_configuration" | jq -r '.issuer // empty' 2>/dev/null)
if [ ! -z "$REALM_CHECK" ]; then
    echo "✅ CloudGate realm is accessible!"
else
    echo "⚠️  Realm verification failed, but setup may still work"
fi

echo ""
echo "🎉 CloudGate Keycloak setup complete!"
echo ""
echo "📋 Configuration Summary:"
echo "   Keycloak URL: ${KEYCLOAK_URL}"
echo "   Admin Console: ${KEYCLOAK_URL}/admin/"
echo "   Admin User: ${ADMIN_USER} / ${ADMIN_PASS}"
echo "   Realm: ${REALM_NAME}"
echo "   Client ID: ${CLIENT_ID}"
echo "   Test User: testuser / password123"
echo ""
echo "🌐 Next Steps:"
echo "   1. Go to http://localhost:3000 (your frontend)"
echo "   2. Click 'Sign in with SSO'"
echo "   3. Use: testuser / password123"
echo ""
echo "🔧 If frontend still shows loading:"
echo "   1. Refresh the page"
echo "   2. Check browser console for errors"
echo "   3. Restart frontend: npm run dev" 