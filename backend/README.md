# CloudGate Backend

## Architecture Overview

The CloudGate backend is organized into separate modules for better maintainability and code organization.

## File Structure

### `main.go` (39 lines)
- **Purpose**: Entry point of the application
- **Contains**: Only the `main()` function and server initialization
- **Responsibilities**: 
  - Load configuration
  - Setup Gin router
  - Configure middleware
  - Setup routes
  - Start the server

### `config.go` (35 lines)
- **Purpose**: Application configuration management
- **Contains**: 
  - `Config` struct definition
  - `LoadConfig()` function
  - `getEnv()` helper function
- **Responsibilities**: Environment variable handling and configuration loading

### `types.go` (49 lines)
- **Purpose**: Data structure definitions
- **Contains**: All struct definitions used throughout the application
- **Types**:
  - `TokenIntrospectionResponse`
  - `UserSession`
  - `TokenIntrospectionRequest`
  - `HealthResponse`
  - `APIInfoResponse`

### `middleware.go` (28 lines)
- **Purpose**: HTTP middleware functions
- **Contains**:
  - `SetupCORS()` - CORS configuration
  - `SecurityHeadersMiddleware()` - Security headers
- **Responsibilities**: Request/response processing and security

### `routes.go` (20 lines)
- **Purpose**: Route configuration and organization
- **Contains**: `SetupRoutes()` function
- **Responsibilities**: API endpoint registration and routing

### `handlers.go` (157 lines)
- **Purpose**: HTTP request handlers
- **Contains**: All endpoint handler functions
- **Handlers**:
  - `HealthCheckHandler` - Health check endpoint
  - `APIInfoHandler` - API information endpoint
  - `TokenIntrospectionHandler` - JWT token validation
  - `UserInfoHandler` - User information retrieval

## Benefits of This Structure

1. **Separation of Concerns**: Each file has a single responsibility
2. **Maintainability**: Easy to locate and modify specific functionality
3. **Testability**: Individual components can be tested in isolation
4. **Readability**: Clean, organized code structure
5. **Scalability**: Easy to add new features without cluttering main.go

## Building and Running

```bash
# Build the application
docker build -t cloudgate-backend .

# Run with Docker Compose
docker-compose up backend

# Or run directly (requires Go 1.23+)
go run .
```

## Environment Variables

See `config.go` for all supported environment variables:
- `KEYCLOAK_URL` - Keycloak server URL
- `KEYCLOAK_REALM` - Keycloak realm name
- `KEYCLOAK_CLIENT_ID` - Client ID for authentication
- `PORT` - Server port (default: 8081)
- `ALLOWED_ORIGINS` - CORS allowed origins 