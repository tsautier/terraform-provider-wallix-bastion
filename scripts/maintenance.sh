#!/bin/bash

# Quick maintenance script for terraform-provider-wallix-bastion
# Performs common maintenance tasks without creating a release

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

update_deps() {
    log_info "Updating Go dependencies..."
    
    # Update dependencies more conservatively
    log_info "Running go mod tidy..."
    go mod tidy
    
    log_info "Updating direct dependencies..."
    # Update only direct dependencies to avoid tool conflicts
    go get -u github.com/hashicorp/terraform-plugin-sdk/v2
    go get -u github.com/hashicorp/go-cleanhttp
    
    log_info "Final tidy and verify..."
    go mod tidy
    go mod verify
    
    log_success "Dependencies updated"
}

run_lints() {
    log_info "Running linters..."
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run
        log_success "Linting passed"
    else
        log_warning "golangci-lint not found! Please install it from https://golangci-lint.run/usage/install/"
    fi
}

run_tests() {
    log_info "Running tests..."
    go test ./...
    log_success "Tests passed"
}

run_build() {
    log_info "Building provider..."
    make build
    log_success "Build successful"
}

run_install() {
    log_info "Installing provider locally..."
    make install
    log_success "Provider installed locally"
}

clean_build() {
    log_info "Cleaning build artifacts..."
    rm -f terraform-provider-wallix-bastion
    rm -rf dist/
    log_success "Build artifacts cleaned"
}

show_help() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  deps      Update Go dependencies"
    echo "  lint      Run golangci-lint"
    echo "  test      Run unit tests"
    echo "  build     Build the provider"
    echo "  install   Install provider locally"
    echo "  clean     Clean build artifacts"
    echo "  all       Run deps, lint, test, build (default)"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0            # Run all checks"
    echo "  $0 deps       # Only update dependencies"
    echo "  $0 test       # Only run tests"
}

main() {
    local command=${1:-all}
    
    case $command in
        deps)
            update_deps
            ;;
        lint)
            run_lints
            ;;
        test)
            run_tests
            ;;
        build)
            run_build
            ;;
        install)
            run_install
            ;;
        clean)
            clean_build
            ;;
        all)
            update_deps
            run_lints
            run_tests
            run_build
            log_success "All checks completed"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "Unknown command: $command"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

main "$@"