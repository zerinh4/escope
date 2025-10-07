# Contributing to Escope

Thank you for your interest in contributing to Escope! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

This project follows a code of conduct that ensures a welcoming environment for all contributors. Please be respectful and constructive in all interactions.

## Getting Started

### Prerequisites

- Go 1.24 or later
- Git
- Make (optional, for using Makefile commands)

### Setting Up Development Environment

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/escope.git
   cd escope
   ```

3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/mertbahardogan/escope.git
   ```

4. **Install dependencies**:
   ```bash
   make tidy
   ```

5. **Build the project**:
   ```bash
   make build
   ```

## Development Workflow

### 1. Create a Feature Branch

Always create a new branch for your changes:

```bash
# Create and switch to a new feature branch
git checkout -b feat/your-feature-name

# Or for bug fixes
git checkout -b fix/your-bug-description

# Or for documentation
git checkout -b docs/your-doc-update
```

**Branch Naming Convention:**
- `feat/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test improvements
- `chore/` - Maintenance tasks

### 2. Make Your Changes

- Write clean, readable code
- Follow the existing code style
- Add comments for complex logic
- Update documentation if needed

### 3. Test Your Changes

**Before submitting any changes, you MUST run the test suite:**

```bash
# Run all command tests
make test-commands

# Run code formatting
make fmt

# Run linting (if golangci-lint is installed)
make lint
```

**The `make test-commands` is mandatory** - it tests all CLI commands and ensures nothing is broken.

### 4. Commit Your Changes

Use clear, descriptive commit messages:

```bash
git add .
git commit -m "feat: add new node distribution analysis command

- Add node dist command to show shard distribution
- Include balance analysis with load statistics
- Update README with new command examples"
```

**Commit Message Format:**
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Test additions/changes
- `chore:` - Maintenance tasks

### 5. Push and Create Pull Request

```bash
# Push your branch
git push origin feat/your-feature-name

# Create a Pull Request on GitHub
```

## Pull Request Process

### Before Submitting

1. **Run all tests**: `make test-commands`
2. **Check code formatting**: `make fmt`
3. **Verify documentation**: Update README.md if needed
4. **Test manually**: Test your changes with real Elasticsearch cluster

### Pull Request Template

When creating a PR, include:

1. **Description**: What changes were made and why
2. **Testing**: How you tested the changes
3. **Screenshots**: If UI changes were made
4. **Breaking Changes**: Any breaking changes (if applicable)
5. **Checklist**: Complete the PR checklist

### PR Checklist

- [ ] Code follows the project's coding standards
- [ ] Self-review of the code has been performed
- [ ] Code has been commented, particularly in hard-to-understand areas
- [ ] Corresponding changes to documentation have been made
- [ ] `make test-commands` has been run and all tests pass
- [ ] `make fmt` has been run and code is properly formatted
- [ ] No new warnings or errors have been introduced
- [ ] Changes have been tested with a real Elasticsearch cluster

### Review Process

1. **Automated Checks**: All CI checks must pass
2. **Code Review**: At least one maintainer must approve
3. **Testing**: Changes must be tested by maintainers
4. **Documentation**: Documentation must be updated if needed

## Coding Standards

### Go Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Keep functions small and focused
- Add comments for exported functions
- Use proper error handling

### Project Structure

```
escope/
├── cmd/                    # CLI commands
│   ├── check/             # Health check command
│   ├── cluster/           # Cluster info command
│   ├── config/            # Configuration command
│   └── ...
├── internal/              # Internal packages
│   ├── config/            # Configuration management
│   ├── elastic/           # Elasticsearch client
│   ├── services/          # Business logic
│   └── ui/                # UI formatters
├── main.go               # Application entry point
├── Makefile              # Build and test commands
└── README.md             # Project documentation
```

### Error Handling

- Always handle errors properly
- Use descriptive error messages
- Log errors when appropriate
- Return meaningful error codes

## Testing

### Command Testing

All commands must be tested using the test suite:

```bash
# Test all commands
make test-commands

# Test specific command (example)
./escope --host http://localhost:9200 cluster
```

### Manual Testing

Before submitting changes:

1. **Test with real Elasticsearch cluster**
2. **Test all command variations**
3. **Test error scenarios**
4. **Verify output formatting**

### Test Data

- Use non-production Elasticsearch clusters
- Create test indices with sample data
- Test with different cluster configurations

## Documentation

### README Updates

When adding new features:

1. **Update command reference table**
2. **Add usage examples**
3. **Update installation instructions if needed**
4. **Add new configuration options**

### Code Documentation

- Document all exported functions
- Add package-level documentation
- Include usage examples in comments
- Document configuration options

## Release Process

### Version Bumping

- **Patch**: Bug fixes, documentation updates
- **Minor**: New features, new commands
- **Major**: Breaking changes

### Release Checklist

- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Version is bumped
- [ ] Changelog is updated
- [ ] Release notes are written

## Getting Help

- **Issues**: Use GitHub Issues for bug reports and feature requests
- **Discussions**: Use GitHub Discussions for questions and ideas
- **Email**: Contact maintainers directly for sensitive issues

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to Escope! 🚀
