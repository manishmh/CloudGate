# CloudGate SSO Portal - Project Vision

## Overview
CloudGate is an enterprise-grade SSO portal designed to mirror HENNGE One's capabilities, providing unified identity management across 370+ SaaS applications with advanced security features.

## 1. SSO Core Components

| Component | Technology/Framework | Purpose |
|-----------|---------------------|---------|
| **Protocols** | OAuth 2.0, OpenID Connect (modern apps), SAML 2.0 (legacy apps) | Standardized authentication/authorization for SaaS integrations |
| **Identity Provider** | Keycloak (open-source) or Auth0 (cloud-based) | Centralized identity management, session handling, and token issuance |
| **User Directory** | LDAP (OpenLDAP), Azure Active Directory, or AWS Cognito | Store and manage user identities, roles, and permissions |
| **Backend** | Node.js (Express), Python (Django), or Go (Gin) | Handle business logic, API integrations, and policy enforcement |
| **Frontend** | React or Vue.js with OIDC Client libraries | User-friendly login portal and dashboard for app access |

## 2. MFA Integration

| Layer | Tools/APIs | Implementation Example |
|-------|------------|----------------------|
| **MFA Methods** | - Time-based OTP: Google Authenticator, Authy<br>- Push Notifications: Duo Security, Okta Verify<br>- Biometrics: WebAuthn, FIDO2 | Integrate via SDKs (e.g., speakeasy for TOTP, @simplewebauthn/browser for WebAuthn) |
| **Adaptive Policies** | Okta Adaptive MFA or Cisco Duo | Configure risk-based rules (e.g., require biometrics for new devices) |

## 3. Adaptive Authentication Engine

| Component | Technology | Use Case |
|-----------|------------|----------|
| **Risk Scoring** | Elasticsearch + Python ML (scikit-learn/TensorFlow) | Analyze login context (IP, device, time) to calculate risk scores |
| **Policy Engine** | Open Policy Agent (OPA) or AWS Cedar | Enforce dynamic access rules (e.g., block logins from high-risk regions) |
| **Context Signals** | - Device fingerprinting: Fingerprint.js<br>- IP Geolocation: MaxMind API | Detect anomalies (e.g., unrecognized devices or locations) |

## 4. SaaS Integrations

| Application | Integration Method | Code Snippet Example (OIDC) |
|-------------|-------------------|----------------------------|
| **Google Workspace** | OpenID Connect with google-auth-library | `const {OAuth2Client} = require('google-auth-library');` |
| **Microsoft 365** | MSAL.js for OAuth 2.0 | `const msalConfig = { auth: { clientId: "YOUR_CLIENT_ID" } };` |
| **Slack** | SAML 2.0 via Keycloak/Auth0 | Configure SAML assertions in IdP metadata |

## 5. Deployment & Security

| Layer | Tools | Best Practices |
|-------|-------|----------------|
| **Infrastructure** | Kubernetes (EKS/GKE) + Terraform | Isolate authentication services in private subnets |
| **Secrets Management** | HashiCorp Vault or AWS Secrets Manager | Rotate API keys and certificates automatically |
| **Audit Logging** | ELK Stack (Elasticsearch, Logstash, Kibana) | Track SSO events for compliance (e.g., GDPR, ISO 27001) |

## Key Features Mirroring HENNGE One

### Unified Identity Fabric
- Use Keycloak as the central IdP to manage SSO across 370+ SaaS apps via pre-built connectors

### Step-Up Authentication
- Implement risk-based MFA using Duo Security's API:
```python
if risk_score > 0.7:
    duo_auth_api.trigger_push(user_device_id)
```

### Zero-Trust Policies
- Enforce device trust with certificates (e.g., AWS Certificate Manager) and session timeouts after inactivity

## Phase 1 Implementation Plan

### 1. Next.js Frontend (App Router)
- Initialize with TypeScript using `create-next-app`
- Create `/login` page with Keycloak OIDC integration using `@react-keycloak/ssr`
- Configure SSO client settings via environment variables:
  - `NEXT_PUBLIC_KEYCLOAK_URL`
  - `NEXT_PUBLIC_KEYCLOAK_REALM`
  - `NEXT_PUBLIC_KEYCLOAK_CLIENT_ID`

### 2. Go Backend (Gin Framework)
- Create `/token/introspect` endpoint for JWT validation
- Integrate Keycloak Admin REST API client
- Add middleware for CORS and secure headers

### 3. Docker Setup
- Keycloak container with pre-configured realm for SSO
- PostgreSQL container for Keycloak database
- Go API service container

### 4. Basic CI/CD Pipeline
- Multi-stage Dockerfile for Go backend
- GitHub Actions workflow for linting/testing

### 5. TypeScript Interfaces
- `UserSession` (id, email, roles)
- `TokenResponse` (access_token, refresh_token)

### 6. Security Configuration
- Sample `.env.local` template with security best practices 