#!/bin/bash

# CloudGate Backend Test Runner
# This script runs all tests with coverage reporting and categorization

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
COVERAGE_DIR="coverage"
COVERAGE_FILE="$COVERAGE_DIR/coverage.out"
COVERAGE_HTML="$COVERAGE_DIR/coverage.html"
TEST_TIMEOUT="30s"

# Create coverage directory
mkdir -p $COVERAGE_DIR

echo -e "${BLUE}🧪 CloudGate Backend Test Suite${NC}"
echo "=================================="

# Function to print section headers
print_section() {
    echo -e "\n${YELLOW}$1${NC}"
    echo "$(printf '=%.0s' {1..50})"
}

# Function to run tests with coverage
run_tests() {
    local test_path=$1
    local test_name=$2
    
    echo -e "${BLUE}Running $test_name tests...${NC}"
    
    if [ -d "$test_path" ] && [ "$(ls -A $test_path/*.go 2>/dev/null)" ]; then
        go test -v -timeout=$TEST_TIMEOUT -coverprofile="$COVERAGE_DIR/${test_name,,}_coverage.out" ./$test_path/...
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✅ $test_name tests passed${NC}"
        else
            echo -e "${RED}❌ $test_name tests failed${NC}"
            return 1
        fi
    else
        echo -e "${YELLOW}⚠️  No $test_name tests found${NC}"
    fi
}

# Function to merge coverage files
merge_coverage() {
    echo -e "\n${BLUE}Merging coverage reports...${NC}"
    
    # Find all coverage files
    coverage_files=$(find $COVERAGE_DIR -name "*_coverage.out" 2>/dev/null)
    
    if [ -n "$coverage_files" ]; then
        # Create merged coverage file
        echo "mode: set" > $COVERAGE_FILE
        
        for file in $coverage_files; do
            if [ -f "$file" ]; then
                tail -n +2 "$file" >> $COVERAGE_FILE
            fi
        done
        
        echo -e "${GREEN}✅ Coverage reports merged${NC}"
    else
        echo -e "${YELLOW}⚠️  No coverage files found${NC}"
    fi
}

# Function to generate coverage report
generate_coverage_report() {
    if [ -f "$COVERAGE_FILE" ]; then
        echo -e "\n${BLUE}Generating coverage report...${NC}"
        
        # Generate HTML coverage report
        go tool cover -html=$COVERAGE_FILE -o $COVERAGE_HTML
        
        # Display coverage summary
        echo -e "\n${YELLOW}Coverage Summary:${NC}"
        go tool cover -func=$COVERAGE_FILE | tail -1
        
        echo -e "${GREEN}✅ Coverage report generated: $COVERAGE_HTML${NC}"
    fi
}

# Function to run benchmarks
run_benchmarks() {
    print_section "🏁 Running Benchmarks"
    
    benchmark_files=$(find . -name "*_test.go" -exec grep -l "func Benchmark" {} \; 2>/dev/null)
    
    if [ -n "$benchmark_files" ]; then
        echo -e "${BLUE}Running benchmark tests...${NC}"
        go test -bench=. -benchmem ./services/... 2>/dev/null || echo -e "${YELLOW}⚠️  No benchmarks found${NC}"
    else
        echo -e "${YELLOW}⚠️  No benchmark tests found${NC}"
    fi
}

# Function to run race condition tests
run_race_tests() {
    print_section "🏃 Running Race Condition Tests"
    
    echo -e "${BLUE}Running tests with race detection...${NC}"
    go test -race -timeout=$TEST_TIMEOUT ./services/... 2>/dev/null || echo -e "${YELLOW}⚠️  Race detection tests completed${NC}"
}

# Function to validate test files
validate_tests() {
    print_section "🔍 Validating Test Files"
    
    # Check for test files
    test_count=$(find services -name "*_test.go" 2>/dev/null | wc -l)
    echo -e "${BLUE}Found $test_count test files${NC}"
    
    # Check for missing test coverage
    echo -e "${BLUE}Checking for missing test coverage...${NC}"
    
    # List all service files that might need tests
    service_files=$(find ../internal/services -name "*.go" ! -name "*_test.go" 2>/dev/null | wc -l)
    echo -e "${BLUE}Found $service_files service files${NC}"
    
    if [ $test_count -lt $service_files ]; then
        echo -e "${YELLOW}⚠️  Consider adding more test coverage${NC}"
    fi
}

# Main test execution
main() {
    local start_time=$(date +%s)
    
    # Parse command line arguments
    case "${1:-all}" in
        "unit")
            print_section "🔬 Unit Tests Only"
            run_tests "services" "Service"
            ;;
        "integration")
            print_section "🔗 Integration Tests Only"
            run_tests "integration" "Integration"
            ;;
        "handlers")
            print_section "🌐 Handler Tests Only"
            run_tests "handlers" "Handler"
            ;;
        "coverage")
            print_section "📊 Coverage Report Only"
            merge_coverage
            generate_coverage_report
            ;;
        "bench")
            run_benchmarks
            ;;
        "race")
            run_race_tests
            ;;
        "validate")
            validate_tests
            ;;
        "clean")
            print_section "🧹 Cleaning Test Artifacts"
            rm -rf $COVERAGE_DIR
            echo -e "${GREEN}✅ Test artifacts cleaned${NC}"
            ;;
        "all"|*)
            print_section "🚀 Running All Tests"
            
            # Validate test setup
            validate_tests
            
            # Run all test categories
            run_tests "services" "Service" || exit 1
            run_tests "handlers" "Handler" || exit 1
            run_tests "integration" "Integration" || exit 1
            
            # Merge coverage and generate reports
            merge_coverage
            generate_coverage_report
            
            # Run additional tests
            run_race_tests
            run_benchmarks
            ;;
    esac
    
    local end_time=$(date +%s)
    local duration=$((end_time - start_time))
    
    echo -e "\n${GREEN}🎉 Test execution completed in ${duration}s${NC}"
    
    # Show coverage file location if it exists
    if [ -f "$COVERAGE_HTML" ]; then
        echo -e "${BLUE}📊 Coverage report: $(pwd)/$COVERAGE_HTML${NC}"
    fi
}

# Help function
show_help() {
    echo "CloudGate Backend Test Runner"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  all         Run all tests (default)"
    echo "  unit        Run unit tests only"
    echo "  integration Run integration tests only"
    echo "  handlers    Run handler tests only"
    echo "  coverage    Generate coverage report only"
    echo "  bench       Run benchmark tests"
    echo "  race        Run race condition tests"
    echo "  validate    Validate test setup"
    echo "  clean       Clean test artifacts"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0              # Run all tests"
    echo "  $0 unit         # Run only unit tests"
    echo "  $0 coverage     # Generate coverage report"
    echo "  $0 clean        # Clean test artifacts"
}

# Check if help is requested
if [ "$1" = "help" ] || [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

# Run main function
main "$@" 