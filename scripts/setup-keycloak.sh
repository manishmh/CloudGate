#!/bin/bash

# CloudGate Keycloak Setup Script
# This script configures Keycloak with the necessary realm and client settings

set -e

echo "🔐 CloudGate Keycloak Configuration Setup"
echo "========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Keycloak configuration
KEYCLOAK_URL="http://localhost:8080"
ADMIN_USER="admin"
ADMIN_PASSWORD="admin_password"
REALM_NAME="cloudgate"
CLIENT_ID="cloudgate-frontend"

# Wait for Keycloak to be ready
print_status "Waiting for Keycloak to be ready..."
timeout=60
while ! curl -f ${KEYCLOAK_URL}/health/ready >/dev/null 2>&1; do
    sleep 2
    timeout=$((timeout - 2))
    if [ $timeout -le 0 ]; then
        print_error "Keycloak is not ready. Please ensure it's running."
        exit 1
    fi
done
print_success "Keycloak is ready!"

# Get admin access token
print_status "Getting admin access token..."
ADMIN_TOKEN=$(curl -s -X POST "${KEYCLOAK_URL}/realms/master/protocol/openid-connect/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=${ADMIN_USER}" \
    -d "password=${ADMIN_PASSWORD}" \
    -d "grant_type=password" \
    -d "client_id=admin-cli" | jq -r '.access_token')

if [ "$ADMIN_TOKEN" = "null" ] || [ -z "$ADMIN_TOKEN" ]; then
    print_error "Failed to get admin access token. Check Keycloak admin credentials."
    exit 1
fi
print_success "Admin access token obtained!"

# Check if realm exists
print_status "Checking if realm '${REALM_NAME}' exists..."
REALM_EXISTS=$(curl -s -H "Authorization: Bearer ${ADMIN_TOKEN}" \
    "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}" | jq -r '.realm // empty')

if [ -z "$REALM_EXISTS" ]; then
    print_status "Creating realm '${REALM_NAME}'..."
    curl -s -X POST "${KEYCLOAK_URL}/admin/realms" \
        -H "Authorization: Bearer ${ADMIN_TOKEN}" \
        -H "Content-Type: application/json" \
        -d '{
            "realm": "'${REALM_NAME}'",
            "enabled": true,
            "displayName": "CloudGate",
            "displayNameHtml": "<div class=\"kc-logo-text\"><span>CloudGate</span></div>",
            "registrationAllowed": true,
            "registrationEmailAsUsername": false,
            "rememberMe": true,
            "verifyEmail": false,
            "loginWithEmailAllowed": true,
            "duplicateEmailsAllowed": false,
            "resetPasswordAllowed": true,
            "editUsernameAllowed": false,
            "bruteForceProtected": true,
            "permanentLockout": false,
            "maxFailureWaitSeconds": 900,
            "minimumQuickLoginWaitSeconds": 60,
            "waitIncrementSeconds": 60,
            "quickLoginCheckMilliSeconds": 1000,
            "maxDeltaTimeSeconds": 43200,
            "failureFactor": 30,
            "defaultRoles": ["default-roles-'${REALM_NAME}'", "offline_access", "uma_authorization"],
            "requiredCredentials": ["password"],
            "passwordPolicy": "length(8)",
            "otpPolicyType": "totp",
            "otpPolicyAlgorithm": "HmacSHA1",
            "otpPolicyInitialCounter": 0,
            "otpPolicyDigits": 6,
            "otpPolicyLookAheadWindow": 1,
            "otpPolicyPeriod": 30,
            "sslRequired": "external",
            "accessTokenLifespan": 300,
            "accessTokenLifespanForImplicitFlow": 900,
            "ssoSessionIdleTimeout": 1800,
            "ssoSessionMaxLifespan": 36000,
            "offlineSessionIdleTimeout": 2592000,
            "accessCodeLifespan": 60,
            "accessCodeLifespanUserAction": 300,
            "accessCodeLifespanLogin": 1800,
            "actionTokenGeneratedByAdminLifespan": 43200,
            "actionTokenGeneratedByUserLifespan": 300,
            "enabled": true,
            "sslRequired": "external",
            "registrationAllowed": true,
            "registrationEmailAsUsername": false,
            "rememberMe": true,
            "verifyEmail": false,
            "loginWithEmailAllowed": true,
            "duplicateEmailsAllowed": false,
            "resetPasswordAllowed": true,
            "editUsernameAllowed": false,
            "bruteForceProtected": true
        }'
    print_success "Realm '${REALM_NAME}' created!"
else
    print_success "Realm '${REALM_NAME}' already exists!"
fi

# Check if client exists
print_status "Checking if client '${CLIENT_ID}' exists..."
CLIENT_EXISTS=$(curl -s -H "Authorization: Bearer ${ADMIN_TOKEN}" \
    "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/clients?clientId=${CLIENT_ID}" | jq -r '.[0].id // empty')

if [ -z "$CLIENT_EXISTS" ]; then
    print_status "Creating client '${CLIENT_ID}'..."
    curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/clients" \
        -H "Authorization: Bearer ${ADMIN_TOKEN}" \
        -H "Content-Type: application/json" \
        -d '{
            "clientId": "'${CLIENT_ID}'",
            "name": "CloudGate Frontend",
            "description": "CloudGate Frontend Application",
            "enabled": true,
            "clientAuthenticatorType": "client-secret",
            "redirectUris": [
                "http://localhost:3000/*",
                "http://localhost:3001/*"
            ],
            "webOrigins": [
                "http://localhost:3000",
                "http://localhost:3001"
            ],
            "protocol": "openid-connect",
            "publicClient": true,
            "standardFlowEnabled": true,
            "implicitFlowEnabled": false,
            "directAccessGrantsEnabled": true,
            "serviceAccountsEnabled": false,
            "authorizationServicesEnabled": false,
            "fullScopeAllowed": true,
            "nodeReRegistrationTimeout": -1,
            "defaultClientScopes": [
                "web-origins",
                "role_list",
                "profile",
                "roles",
                "email"
            ],
            "optionalClientScopes": [
                "address",
                "phone",
                "offline_access",
                "microprofile-jwt"
            ],
            "access": {
                "view": true,
                "configure": true,
                "manage": true
            }
        }'
    print_success "Client '${CLIENT_ID}' created!"
else
    print_success "Client '${CLIENT_ID}' already exists!"
fi

# Create a test user
print_status "Creating test user..."
TEST_USER_EXISTS=$(curl -s -H "Authorization: Bearer ${ADMIN_TOKEN}" \
    "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/users?username=testuser" | jq -r '.[0].id // empty')

if [ -z "$TEST_USER_EXISTS" ]; then
    curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${REALM_NAME}/users" \
        -H "Authorization: Bearer ${ADMIN_TOKEN}" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "testuser",
            "email": "test@cloudgate.dev",
            "firstName": "Test",
            "lastName": "User",
            "enabled": true,
            "emailVerified": true,
            "credentials": [{
                "type": "password",
                "value": "testpass123",
                "temporary": false
            }]
        }'
    print_success "Test user created! (username: testuser, password: testpass123)"
else
    print_success "Test user already exists! (username: testuser, password: testpass123)"
fi

echo ""
print_success "🎉 Keycloak configuration completed successfully!"
echo ""
echo "📋 Configuration Summary:"
echo "========================"
echo "🌐 Keycloak URL:     ${KEYCLOAK_URL}"
echo "🏰 Realm:            ${REALM_NAME}"
echo "📱 Client ID:        ${CLIENT_ID}"
echo "👤 Admin Console:    ${KEYCLOAK_URL}/admin"
echo "   Username:         ${ADMIN_USER}"
echo "   Password:         ${ADMIN_PASSWORD}"
echo ""
echo "🧪 Test User:"
echo "   Username:         testuser"
echo "   Password:         testpass123"
echo "   Email:            test@cloudgate.dev"
echo ""
echo "🚀 You can now test authentication at: http://localhost:3000"
echo ""
print_success "Happy coding! 🎉" 