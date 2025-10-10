# Contributing to Terraform Provider Wallix Bastion

Thank you for your interest in contributing to the Terraform Provider for Wallix Bastion! This document provides guidelines and information for contributors.

## Code of Conduct

Please note we have a [code of conduct](CODE_OF_CONDUCT.md), please follow it in all your interactions with the project.

## Getting Started

### Prerequisites

Before you begin contributing, ensure you have the following installed:

- [Go](https://golang.org/doc/install) version 1.22 or 1.23
- [Terraform](https://www.terraform.io/downloads.html) version 1.0 or later
- [Git](https://git-scm.com/downloads)
- Access to a Wallix Bastion instance for testing (recommended)

### Development Setup

1. **Fork and Clone the Repository**

   ```bash
   git clone https://github.com/your-username/terraform-provider-wallix-bastion.git
   cd terraform-provider-wallix-bastion
   ```

2. **Set Up Development Environment**

   ```bash
   # Install dependencies and development tools
   make setup-dev
   
   # Verify setup
   make dev-check
   ```

3. **Build the Provider Locally**

   ```bash
   make build
   make install
   ```

## Types of Contributions

We welcome several types of contributions:

### 🐛 Bug Reports

- Use the GitHub issue template
- Include provider version, Terraform version, and Bastion API version
- Provide minimal reproduction steps
- Include relevant logs and error messages

### ✨ Feature Requests

- Check existing issues to avoid duplicates
- Describe the use case and expected behavior
- Consider if the feature aligns with the provider's scope

### 📖 Documentation Improvements

- Fix typos, clarify instructions, or add examples
- Follow our [documentation guidelines](#documentation-guidelines)

### 🔧 Code Contributions

- Bug fixes, new resources, or enhancements
- Follow our [development guidelines](#development-guidelines)

## Development Guidelines

### Project Structure

```ini
├── bastion/                    # Provider implementation
│   ├── client.go              # API client
│   ├── provider.go            # Provider configuration
│   ├── resource_*.go          # Resource implementations
│   ├── data_source_*.go       # Data source implementations
│   └── *_test.go              # Tests
├── docs/                      # Generated documentation
├── examples/                  # Usage examples
├── templates/                 # Documentation templates
├── scripts/                   # Build and maintenance scripts
└── Makefile                   # Build automation
```

### Adding New Resources

1. **Create Resource Implementation**

   ```bash
   # Create the resource file
   touch bastion/resource_new_resource.go
   
   # Create the test file
   touch bastion/resource_new_resource_test.go
   ```

2. **Implement Required Functions**
   - `CreateContext` - Create resource
   - `ReadContext` - Read resource state
   - `UpdateContext` - Update resource (if applicable)
   - `DeleteContext` - Delete resource
   - `Schema` - Define resource schema

3. **Add to Provider**
   Update `bastion/provider.go` to register the new resource.

4. **Write Tests**
   - Unit tests for each CRUD operation
   - Acceptance tests for real API interactions
   - Test edge cases and error conditions

### Code Style and Standards

- **Go Code**: Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- **Formatting**: Use `make fmt` to format code consistently
- **Linting**: Run `make lint` to check for issues
- **Testing**: Maintain test coverage with `make test-coverage`

### API Client Guidelines

- Use the existing client patterns in `bastion/client.go`
- Handle API errors gracefully with informative messages
- Implement proper retry logic for transient failures
- Log API requests/responses for debugging (when appropriate)

### Resource Implementation Best Practices

1. **Schema Definition**

   ```go
   func resourceNewResource() *schema.Resource {
       return &schema.Resource{
           CreateContext: resourceNewResourceCreate,
           ReadContext:   resourceNewResourceRead,
           UpdateContext: resourceNewResourceUpdate,
           DeleteContext: resourceNewResourceDelete,
           
           Schema: map[string]*schema.Schema{
               "name": {
                   Type:        schema.TypeString,
                   Required:    true,
                   Description: "Name of the resource",
               },
           },
       }
   }
   ```

2. **Error Handling**

   ```go
   if err != nil {
       return diag.FromErr(fmt.Errorf("failed to create resource: %w", err))
   }
   ```

3. **State Management**
   - Always set the resource ID
   - Handle partial state updates gracefully
   - Use `d.Partial(true)` for complex updates

### Testing Guidelines

#### Unit Tests

```bash
# Run unit tests
make test

# Run with coverage
make test-coverage

# Run specific test
go test -v ./bastion -run TestResourceNewResource
```

#### Acceptance Tests

```bash
# Set up test environment
export WALLIX_BASTION_HOST="test-bastion"
export WALLIX_BASTION_TOKEN="test-token"
export WALLIX_BASTION_USER="admin"
export TF_ACC=1

# Run acceptance tests
make testacc

# Run specific acceptance test
TF_ACC=1 go test -v ./bastion -run TestAccResourceNewResource_basic
```

#### Test Best Practices

- Test both success and failure scenarios
- Use test fixtures for consistent data
- Clean up resources in test teardown
- Mock external dependencies when possible

## Pull Request Process

### Before Submitting

1. **Discuss Changes**: For significant changes, open an issue first to discuss the approach
2. **Check Existing Work**: Search for existing issues and PRs to avoid duplication
3. **Follow Guidelines**: Ensure your contribution follows this guide

### PR Requirements

1. **Clean Build**: Ensure `make dev-check` passes without errors
2. **Tests**: Add or update tests for your changes
3. **Documentation**: Update documentation for any user-facing changes
4. **Examples**: Add usage examples when introducing new resources
5. **Changelog**: Update CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/) format

### Submission Steps

1. **Create Feature Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make Changes**
   - Implement your feature/fix
   - Add comprehensive tests
   - Update documentation

3. **Verify Changes**

   ```bash
   # Run full verification
   make test-all
   make docs-verify
   make lint
   ```

4. **Commit Changes**

   ```bash
   git add .
   git commit -m "feat: add new resource for XYZ"
   git push origin feature/your-feature-name
   ```

5. **Open Pull Request**
   - Use the PR template
   - Link related issues
   - Provide detailed description
   - Request review from maintainers

### PR Review Process

- **Automated Checks**: CI/CD pipeline runs tests and linting
- **Code Review**: At least one maintainer review required
- **Documentation Review**: Verify docs are complete and accurate
- **Testing**: Acceptance tests should pass
- **Approval**: Two approvals required for merge

### Commit Message Guidelines

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```ini
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

**Types:**

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding tests
- `refactor`: Code refactoring
- `chore`: Maintenance tasks

**Examples:**

```ini
feat(resources): add support for connection policies
fix(client): handle API timeout errors gracefully
docs: update installation instructions
test: add acceptance tests for user resource
```

## Documentation Guidelines

When adding new resources or modifying existing ones, follow these documentation guidelines:

### For New Resources/Data Sources

1. **Generate Templates**: Run `make docs` to auto-generate missing templates
2. **Edit Templates**: Customize the generated template in `templates/resources/` or `templates/data-sources/`
3. **Add Examples**: Include comprehensive usage examples covering different scenarios
4. **Verify Quality**: Run `make docs-verify` to check documentation quality

### Documentation Requirements

- Include at least 3-5 different usage examples
- Document all required and optional parameters
- Add security considerations and best practices
- Include troubleshooting information
- Provide import examples
- Cross-reference related resources

### Documentation Workflow

```bash
# 1. Generate documentation templates
make docs

# 2. Edit templates in templates/ directory
# Add examples, usage notes, best practices

# 3. Regenerate and verify documentation
make docs-verify

# 4. Commit both templates and generated docs
git add templates/ docs/
git commit -m "docs: add documentation for new resource"
```

For detailed documentation guidelines, see [DOCUMENTATION.md](./DOCUMENTATION.md).

## Release Process

### Version Management

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Commands

```bash
# Prepare patch release (0.14.7 -> 0.14.8)
make release-patch

# Prepare minor release (0.14.7 -> 0.15.0)
make release-minor

# Prepare major release (0.14.7 -> 1.0.0)
make release-major
```

## Community and Support

### Getting Help

- **Documentation**: Check [Provider Registry](https://registry.terraform.io/providers/wallix/wallix-bastion/latest/docs)
- **Issues**: Search existing [GitHub Issues](https://github.com/wallix/terraform-provider-wallix-bastion/issues)
- **Discussions**: Use [GitHub Discussions](https://github.com/wallix/terraform-provider-wallix-bastion/discussions)

### Reporting Security Issues

For security-related issues, please follow our [Security Policy](SECURITY.md) rather than opening a public issue.

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and community discussions
- **Pull Requests**: Code contributions and reviews

## Maintenance

### Dependency Updates

Dependencies are updated regularly:

```bash
# Update Go dependencies
make update-deps

# Check for outdated dependencies
go list -u -m all
```

### Regular Maintenance Tasks

```bash
# Run full maintenance (deps, lint, test, build)
make maintenance

# Quick development checks
make dev-check
```

## Additional Resources

### External Documentation

- [Terraform Plugin Development](https://developer.hashicorp.com/terraform/plugin)
- [Terraform Provider Framework](https://developer.hashicorp.com/terraform/plugin/framework)
- [Wallix Bastion API Documentation](https://docs.wallix.com/)

### Learning Resources

- [Terraform Provider Development Program](https://www.terraform.io/docs/extend/index.html)
- [HashiCorp Learn - Build Providers](https://learn.hashicorp.com/tutorials/terraform/provider-setup)

### Repository Specific

- [Project README](./README.md) - Getting started and usage
- [Documentation Guide](./DOCUMENTATION.md) - How to work with docs
- [Release Notes](./RELEASE.md) - Release process details
- [Security Policy](./SECURITY.md) - Security reporting

## Recognition

Contributors are recognized in:

- Release notes for significant contributions
- GitHub contributor listings
- Special thanks in README for major features

## Questions?

Don't hesitate to ask questions! Open a [GitHub Discussion](https://github.com/wallix/terraform-provider-wallix-bastion/discussions) or reach out to the maintainers.

Thank you for contributing to make this provider better! 🚀
