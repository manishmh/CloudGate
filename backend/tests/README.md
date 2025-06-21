# CloudGate Backend Test Suite

This directory contains comprehensive tests for the CloudGate backend services.

## 📁 Directory Structure

```
tests/
├── services/           # Service layer unit tests
│   ├── user_service_test.go
│   ├── mfa_service_test.go
│   ├── oauth_monitoring_service_test.go
│   ├── risk_service_test.go
│   ├── session_service_test.go
│   └── user_settings_service_test.go
├── handlers/          # HTTP handler tests (future)
├── integration/       # Integration tests (future)
├── fixtures/          # Test data fixtures (future)
├── run_tests.sh       # Comprehensive test runner script
└── README.md          # This file
```

## 🚀 Running Tests

### Quick Start

```bash
# Run all tests
./run_tests.sh

# Run specific test categories
./run_tests.sh unit        # Unit tests only
./run_tests.sh coverage    # Generate coverage report
./run_tests.sh bench       # Run benchmarks
./run_tests.sh race        # Race condition tests
```

### Manual Test Execution

```bash
# Run all service tests
go test -v ./tests/services/...

# Run specific test file
go test -v ./tests/services/user_service_test.go

# Run with coverage
go test -v -coverprofile=coverage.out ./tests/services/...
go tool cover -html=coverage.out -o coverage.html

# Run benchmarks
go test -bench=. -benchmem ./tests/services/...

# Run with race detection
go test -race ./tests/services/...
```

## 📊 Test Categories

### Service Tests (`tests/services/`)

**User Service Tests** (`user_service_test.go`)
- ✅ User creation and updates
- ✅ User retrieval by ID
- ✅ Demo user management
- ✅ User validation and error handling
- ✅ Benchmark tests for performance

**MFA Service Tests** (`mfa_service_test.go`)
- ✅ MFA setup storage and retrieval
- ✅ MFA enable/disable functionality
- ✅ Backup codes management
- ✅ Backup code usage and validation
- ✅ Error handling for invalid inputs

**OAuth Monitoring Service Tests** (`oauth_monitoring_service_test.go`)
- ✅ Application connection monitoring
- ✅ Connection health checks
- ✅ Usage recording and statistics
- ✅ Security event creation
- ✅ Trusted device management

**Risk Service Tests** (`risk_service_test.go`)
- ✅ Risk assessment storage and retrieval
- ✅ Risk threshold management
- ✅ Device fingerprinting
- ✅ WebAuthn credential management
- ✅ Risk assessment history tracking

**Session Service Tests** (`session_service_test.go`)
- ✅ Session creation and validation
- ✅ Session refresh functionality
- ✅ Session invalidation
- ✅ Session cleanup for expired sessions
- ✅ Session statistics generation

**User Settings Service Tests** (`user_settings_service_test.go`)
- ✅ Default settings creation
- ✅ Settings updates (full and partial)
- ✅ Single setting updates
- ✅ Settings reset to defaults
- ✅ Comprehensive validation

## 🔧 Test Infrastructure

### Database Setup
- Uses in-memory SQLite databases for isolation
- Each test gets its own database instance
- Automatic schema migration for required models
- Proper cleanup after test completion

### Test Patterns
- **Arrange-Act-Assert** pattern for test structure
- **Table-driven tests** for multiple scenarios
- **Subtests** for organized test grouping
- **Benchmark tests** for performance validation
- **Error case testing** for comprehensive coverage

### Mocking Strategy
- Global database variable mocking for service functions
- In-memory database for realistic data persistence
- Isolated test environments to prevent cross-test contamination

## 📈 Coverage Goals

Current test coverage targets:
- **Service Layer**: 90%+ coverage
- **Critical Paths**: 100% coverage (auth, MFA, risk assessment)
- **Error Handling**: 100% coverage
- **Edge Cases**: Comprehensive coverage

## 🏃 Performance Testing

### Benchmarks Available
- User service operations (create, retrieve)
- MFA setup operations
- Session management operations
- Risk assessment calculations

### Running Benchmarks
```bash
# Run all benchmarks
./run_tests.sh bench

# Run specific benchmark
go test -bench=BenchmarkUserService_CreateOrUpdateUser ./tests/services/

# Run benchmarks with memory profiling
go test -bench=. -benchmem ./tests/services/
```

## 🔍 Test Utilities

### Available Helper Functions
- `setupTestDB()` - Creates isolated test database
- `createTestUser()` - Creates test user with default values
- `setupTestService()` - Initializes service with test database
- Database cleanup and migration helpers

### Test Data Management
- Consistent test data across test files
- UUID generation for unique identifiers
- Realistic test scenarios and edge cases

## 🐛 Debugging Tests

### Common Issues
1. **Database Connection Errors**: Ensure SQLite driver is available
2. **Test Isolation**: Each test uses separate database instance
3. **Timing Issues**: Use appropriate timeouts for async operations
4. **Memory Leaks**: Tests clean up database connections properly

### Debug Commands
```bash
# Verbose output
go test -v ./tests/services/...

# Race condition detection
go test -race ./tests/services/...

# CPU profiling
go test -cpuprofile=cpu.prof ./tests/services/...

# Memory profiling
go test -memprofile=mem.prof ./tests/services/...
```

## 📋 Test Checklist

Before submitting changes, ensure:
- [ ] All existing tests pass
- [ ] New features have corresponding tests
- [ ] Error cases are tested
- [ ] Benchmarks are included for performance-critical code
- [ ] Test coverage meets minimum requirements
- [ ] Tests are properly isolated and don't depend on external services
- [ ] Test data is cleaned up properly

## 🔮 Future Enhancements

Planned test additions:
- **Handler Tests**: HTTP endpoint testing
- **Integration Tests**: Full request-response cycle testing
- **Load Tests**: Performance under high load
- **Contract Tests**: API contract validation
- **End-to-End Tests**: Complete user journey testing

## 📞 Support

For questions about testing:
1. Check existing test patterns in the codebase
2. Review this README for common patterns
3. Run `./run_tests.sh help` for script usage
4. Ensure all tests pass before submitting changes

---

**Happy Testing! 🧪** 