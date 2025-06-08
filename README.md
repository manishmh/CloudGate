# 🚀 CloudGate - Enterprise SSO Portal

**Secure Single Sign-On for SaaS Applications**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Next.js](https://img.shields.io/badge/Next.js-14-000000?style=flat&logo=next.js)](https://nextjs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=flat&logo=typescript)](https://www.typescriptlang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

CloudGate is a comprehensive Single Sign-On (SSO) solution that enables secure access to multiple SaaS applications through a unified authentication portal. Built with modern technologies and enterprise-grade security.

---

## 🎯 **Features**

### **🔐 Authentication & Security**
- **Unified SSO**: Single login for multiple SaaS applications
- **Keycloak Integration**: Enterprise-grade identity management
- **OAuth 2.0 & 1.0a**: Support for modern and legacy OAuth flows
- **JWT Validation**: Secure token-based authentication
- **Audit Logging**: Comprehensive security event tracking

### **🔗 SaaS Integrations**
- **Google Workspace**: Gmail, Drive, Calendar, Contacts
- **Microsoft 365**: Outlook, OneDrive, Teams, SharePoint
- **Slack**: Workspaces, Channels, Direct Messages
- **GitHub**: Repositories, Organizations, Issues
- **Trello**: Boards, Cards, Lists (OAuth 1.0a)
- **Extensible**: Easy to add new OAuth providers

### **🎨 Modern Interface**
- **Responsive Design**: Works on desktop, tablet, and mobile
- **Real-time Updates**: Live connection status and notifications
- **OAuth Testing**: Built-in testing interface for developers
- **User Dashboard**: Centralized app management

---

## 🏗️ **Architecture**

### **Tech Stack**

#### **Backend**
- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL (Neon DB Cloud)
- **Authentication**: Keycloak
- **ORM**: GORM
- **Security**: CORS, JWT validation, OAuth flows

#### **Frontend**
- **Framework**: Next.js 14 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: React Hooks
- **HTTP Client**: Fetch API

#### **Infrastructure**
- **Containerization**: Docker & Docker Compose
- **Database**: Neon PostgreSQL (Production)
- **Identity Provider**: Keycloak
- **Development**: Local development stack

### **System Architecture**
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │    │   External      │
│   (Next.js)     │◄──►│   (Go/Gin)      │◄──►│   Services      │
│                 │    │                 │    │                 │
│ • Dashboard     │    │ • OAuth Flows   │    │ • Google        │
│ • Auth UI       │    │ • Token Mgmt    │    │ • Microsoft     │
│ • App Mgmt      │    │ • User Mgmt     │    │ • Slack         │
└─────────────────┘    └─────────────────┘    │ • GitHub        │
                                              │ • Trello        │
┌─────────────────┐    ┌─────────────────┐    └─────────────────┘
│   Keycloak      │    │   Database      │
│   (Identity)    │    │   (Neon DB)     │
│                 │    │                 │
│ • User Auth     │    │ • User Data     │
│ • JWT Tokens    │    │ • OAuth Tokens  │
│ • SSO Sessions  │    │ • Audit Logs    │
└─────────────────┘    └─────────────────┘
```

---

## 🚀 **Quick Start**

### **Prerequisites**
- **Docker & Docker Compose** (for Keycloak)
- **Go 1.21+** (for backend)
- **Node.js 18+** (for frontend)
- **Git** (for version control)

### **1. Clone & Setup**
```bash
# Clone repository
git clone https://github.com/your-username/CloudGate.git
cd CloudGate

# Start Keycloak
docker compose up -d keycloak

# Setup Keycloak realm and users
chmod +x scripts/complete-keycloak-setup.sh
./scripts/complete-keycloak-setup.sh
```

### **2. Backend Setup**
```bash
cd backend

# Configure environment
cp env.example .env
# Edit .env with your Neon DB connection and OAuth credentials

# Install dependencies and run
go mod tidy
go run .
```

### **3. Frontend Setup**
```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev
```

### **4. Access Application**
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8081
- **Keycloak Admin**: http://localhost:8080/admin (admin/admin)
- **OAuth Testing**: http://localhost:3000/dashboard/oauth-test

---

## ⚙️ **Configuration**

### **Environment Variables**

#### **Backend Configuration (`backend/.env`)**
```env
# Database (Neon DB)
NEON_DATABASE_URL=postgresql://user:pass@host/db?sslmode=require

# Keycloak
KEYCLOAK_URL=http://localhost:8080
KEYCLOAK_REALM=cloudgate
KEYCLOAK_CLIENT_ID=cloudgate-frontend

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URI=http://localhost:8081/oauth/google/callback

# Microsoft OAuth
MICROSOFT_CLIENT_ID=your_microsoft_client_id
MICROSOFT_CLIENT_SECRET=your_microsoft_client_secret
MICROSOFT_REDIRECT_URI=http://localhost:8081/oauth/microsoft/callback

# Slack OAuth
SLACK_CLIENT_ID=your_slack_client_id
SLACK_CLIENT_SECRET=your_slack_client_secret
SLACK_REDIRECT_URI=http://localhost:8081/oauth/slack/callback

# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_REDIRECT_URI=http://localhost:8081/oauth/github/callback

# Trello OAuth (OAuth 1.0a)
TRELLO_CLIENT_ID=your_trello_api_key
TRELLO_CLIENT_SECRET=your_trello_secret
TRELLO_REDIRECT_URI=http://localhost:8081/oauth/trello/callback

# Security
ALLOWED_ORIGINS=http://localhost:3000
FRONTEND_URL=http://localhost:3000
```

### **OAuth Provider Setup**

#### **Google OAuth**
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create project and enable APIs (Gmail, Drive, Calendar)
3. Configure OAuth consent screen
4. Create OAuth 2.0 credentials
5. Add redirect URI: `http://localhost:8081/oauth/google/callback`

#### **Microsoft OAuth**
1. Go to [Azure Portal](https://portal.azure.com/)
2. Register new application
3. Configure API permissions (Microsoft Graph)
4. Create client secret
5. Add redirect URI: `http://localhost:8081/oauth/microsoft/callback`

#### **Trello OAuth (OAuth 1.0a)**
1. Go to [Trello Power-Ups](https://trello.com/power-ups/admin)
2. Create new Power-Up
3. Generate API Key and Secret
4. Add allowed origins: `http://localhost:3000`, `http://localhost:8081`

---

## 🔗 **API Reference**

### **Health & Information**
```bash
GET /health                 # Health check
GET /health/db             # Database health
GET /api/info              # API information
```

### **Authentication**
```bash
POST /token/introspect     # Validate JWT token
GET /user/info             # Get user information
```

### **OAuth Flows**
```bash
# OAuth 2.0 Providers
GET /oauth/google/connect      # Initiate Google OAuth
GET /oauth/google/callback     # Google OAuth callback
GET /oauth/microsoft/connect   # Initiate Microsoft OAuth
GET /oauth/slack/connect       # Initiate Slack OAuth
GET /oauth/github/connect      # Initiate GitHub OAuth

# OAuth 1.0a Providers
GET /oauth/trello/connect      # Initiate Trello OAuth
GET /oauth/trello/callback     # Trello OAuth callback
```

### **Application Management**
```bash
GET /apps                  # List connected applications
POST /apps/connect         # Connect new application
POST /apps/launch          # Launch application
```

### **User Management**
```bash
GET /user/profile          # Get user profile
PUT /user/profile          # Update user profile
GET /user/sessions         # Get user sessions
DELETE /user/sessions      # Invalidate sessions
```

---

## 🧪 **Development & Testing**

### **Local Development**
```bash
# Backend with live reload (install air: go install github.com/cosmtrek/air@latest)
cd backend && air

# Frontend with hot reload
cd frontend && npm run dev

# Database migrations (automatic on startup)
cd backend && go run .
```

### **OAuth Testing Interface**
Visit `http://localhost:3000/dashboard/oauth-test` to:
- Test backend connectivity
- Verify OAuth endpoint responses
- Initiate real OAuth flows
- Monitor connection status
- Debug OAuth issues

### **API Testing**
```bash
# Health checks
curl http://localhost:8081/health
curl http://localhost:8081/health/db

# OAuth endpoints (without credentials)
curl http://localhost:8081/oauth/google/connect
# Expected: {"error":"Google OAuth not configured"}

# With credentials configured
curl http://localhost:8081/oauth/google/connect
# Expected: {"auth_url":"https://accounts.google.com/o/oauth2/...", "provider":"google"}
```

---

## 🔒 **Security Features**

### **Authentication Security**
- **JWT Validation**: All protected endpoints validate JWT tokens
- **Keycloak Integration**: Enterprise-grade identity management
- **Session Management**: Secure session handling and timeout
- **CORS Protection**: Configurable cross-origin resource sharing

### **OAuth Security**
- **State Parameters**: CSRF protection for OAuth 2.0 flows
- **PKCE Support**: Enhanced security for public clients
- **Signature Validation**: HMAC-SHA1 signatures for OAuth 1.0a
- **Token Encryption**: Secure storage of OAuth tokens
- **Scope Limitation**: Minimal necessary permissions

### **Data Protection**
- **Environment Variables**: Sensitive data in environment configuration
- **Database Encryption**: Encrypted token storage
- **Audit Logging**: Comprehensive security event tracking
- **HTTPS Enforcement**: Production HTTPS requirements

---

## 🚀 **Deployment**

### **Production Checklist**
- [ ] Configure Neon DB production database
- [ ] Set up OAuth applications with production URLs
- [ ] Configure environment variables securely
- [ ] Enable HTTPS/TLS
- [ ] Set up monitoring and logging
- [ ] Configure backup strategies

### **Environment URLs**
```bash
# Development
Frontend: http://localhost:3000
Backend:  http://localhost:8081
Keycloak: http://localhost:8080

# Production (example)
Frontend: https://cloudgate.your-domain.com
Backend:  https://api.cloudgate.your-domain.com
Keycloak: https://auth.cloudgate.your-domain.com
```

### **Docker Deployment**
```bash
# Build and run with Docker Compose
docker compose up -d

# Production with custom configuration
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

---

## 📊 **Project Status**

### **✅ Completed Features**
- [x] **Phase 1**: Core infrastructure, database, Keycloak setup
- [x] **Phase 2**: OAuth integration (Google, Microsoft, Slack, GitHub, Trello)
- [x] **Security**: JWT validation, CORS, audit logging
- [x] **Frontend**: Dashboard, OAuth testing interface
- [x] **Documentation**: Comprehensive setup and usage guides

### **🔄 In Progress**
- [ ] **Token Refresh**: Automatic OAuth token renewal
- [ ] **Advanced UI**: Enhanced dashboard and user experience
- [ ] **API Access**: Direct SaaS API integration features
- [ ] **Admin Panel**: Administrative interface and controls

### **🎯 Roadmap**
- [ ] **Phase 3**: Additional SaaS providers (Salesforce, Jira, Notion)
- [ ] **Enterprise Features**: SAML, LDAP integration
- [ ] **Analytics**: Usage analytics and reporting
- [ ] **Mobile App**: Native mobile application

---

## 🤝 **Contributing**

### **Development Workflow**
1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### **Code Standards**
- **Go**: Follow Go conventions, use `gofmt` and `golint`
- **TypeScript**: Use ESLint and Prettier for formatting
- **Commits**: Use conventional commit messages
- **Testing**: Add tests for new features

### **Development Setup**
```bash
# Install development tools
go install github.com/cosmtrek/air@latest  # Live reload for Go
npm install -g prettier eslint            # Frontend tools

# Run linters
cd backend && golint ./...
cd frontend && npm run lint
```

---

## 📄 **License**

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🆘 **Support & Troubleshooting**

### **Common Issues**

#### **OAuth "not configured" errors**
- Check environment variables are set correctly
- Restart backend after updating `.env`
- Verify OAuth app credentials

#### **Database connection issues**
- Verify Neon DB connection string
- Check network connectivity
- Ensure database exists and is accessible

#### **Keycloak authentication issues**
- Verify Keycloak is running: `docker compose ps`
- Check realm configuration
- Ensure client settings are correct

### **Getting Help**
- **Issues**: Create an issue in the repository
- **Discussions**: Use GitHub Discussions for questions
- **Documentation**: Check the comprehensive guides in the repository

### **Useful Commands**
```bash
# Check service status
docker compose ps

# View logs
docker compose logs keycloak
cd backend && go run . 2>&1 | grep -i error
cd frontend && npm run dev

# Reset development environment
docker compose down -v
docker compose up -d
```

---

**🚀 Ready to secure your SaaS ecosystem with CloudGate!**

For detailed setup instructions, visit the OAuth testing interface at `http://localhost:3000/dashboard/oauth-test` after starting the application. 