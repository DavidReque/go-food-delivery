# 🚀 Development Guide - Go Food Delivery

## 📋 **Prerequisites**

### **Required Tools:**

- **Go 1.22+** - [Download](https://golang.org/dl/)
- **Node.js 18+** - [Download](https://nodejs.org/)
- **npm 9+** - Included with Node.js
- **Docker** - [Download](https://www.docker.com/)
- **Docker Compose** - Included with Docker Desktop

### **Go Tools:**

```bash
# Install Go tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/cosmtrek/air@latest
```

---

## 🛠️ **Environment Setup**

### **1. Install Development Dependencies:**

```bash
# Install all npm dependencies
npm install

# Configure Husky hooks
npm run prepare
```

### **2. Verify Installation:**

```bash
# Verify everything is configured
npm run setup:dev

# Verify Git hooks
ls -la .husky/
```

---

## 🔧 **Workflow with Conventional Commits**

### **Create Standardized Commits:**

```bash
# Use Commitizen for interactive commits
npm run commit

# Or create commits manually following the format:
git commit -m "feat(orders): add order creation endpoint"
git commit -m "fix(auth): resolve JWT validation issue"
git commit -m "docs(api): update endpoint documentation"
git commit -m "test(orders): add integration tests"
```

### **Commit Format:**

```
type(scope): description

[optional body]

[optional footer]
```

### **Commit Types:**

- **`feat`** - New feature
- **`fix`** - Bug fix
- **`docs`** - Documentation changes
- **`style`** - Code style changes (formatting, missing semicolons, etc)
- **`refactor`** - Code refactoring
- **`perf`** - Performance improvements
- **`test`** - Add or modify tests
- **`build`** - Build system changes
- **`ci`** - CI/CD changes
- **`chore`** - Maintenance tasks

### **Project Scopes:**

- **`orders`** - Orders service
- **`catalogs`** - Catalogs service
- **`customers`** - Customers service
- **`payments`** - Payments service
- **`delivery`** - Delivery service
- **`auth`** - Authentication service
- **`api`** - API layer
- **`core`** - Domain logic
- **`infrastructure`** - Infrastructure layer
- **`testing`** - Testing framework
- **`docs`** - Documentation
- **`ci`** - CI/CD
- **`docker`** - Docker configuration
- **`k8s`** - Kubernetes configuration

---

## 🐕 **Git Hooks with Husky**

### **Configured Hooks:**

#### **Pre-commit:**

- ✅ Automatic Go code formatting
- ✅ Import organization
- ✅ Linting with golangci-lint
- ✅ Unit test execution

#### **Commit-msg:**

- ✅ Commit format validation
- ✅ Conventional Commits verification
- ✅ Scope and type validation

#### **Pre-push:**

- ✅ Complete tests with race detection
- ✅ Integration tests (if available)
- ✅ Coverage verification

### **Run Hooks Manually:**

```bash
# Run pre-commit hooks
npx husky run .husky/pre-commit

# Run commit-msg validation
npx husky run .husky/commit-msg

# Run pre-push hooks
npx husky run .husky/pre-push
```

---

## 🧪 **Testing and Code Quality**

### **Automatic Tests:**

```bash
# Unit tests
go test -v ./...

# Tests with race detection
go test -race ./...

# Tests with coverage
go test -cover ./...

# Integration tests
make integration-test

# End-to-end tests
make e2e-test
```

### **Linting and Formatting:**

```bash
# Automatic formatting
go fmt ./...

# Import organization
goimports -w .

# Complete linting
golangci-lint run

# Fast linting (used in pre-commit)
golangci-lint run --fast
```

---

## 📦 **Versioning and Releases**

### **Generate Releases:**

```bash
# Patch release (1.0.0 → 1.0.1)
npm run release:patch

# Minor release (1.0.0 → 1.1.0)
npm run release:minor

# Major release (1.0.0 → 2.0.0)
npm run release:major

# Changelog preview
npm run changelog
```

### **Versioning Configuration:**

- **`feat:`** → Increments minor version
- **`fix:`** → Increments patch version
- **`breaking:`** → Increments major version
- **Other types** → Increment patch version

---

## 🐳 **Development with Docker**

### **Available Services:**

```bash
# Start all services
docker-compose up -d

# Start only database services
docker-compose up -d mongodb elasticsearch rabbitmq

# Run tests with testcontainers
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

### **Local Development:**

```bash
# Use Air for hot reload
cd internal/services/orderservice
air

# Or run directly
go run cmd/app/main.go
```

---

## 📚 **Useful Commands**

### **NPM Scripts:**

```bash
# Complete environment setup
npm run setup:dev

# Create interactive commit
npm run commit

# Generate release
npm run release

# Validate recent commits
npm run validate:commits

# Clean dependencies
npm run clean
```

### **Make Commands:**

```bash
# Build all services
make build

# Unit tests
make unit-test

# Integration tests
make integration-test

# End-to-end tests
make e2e-test

# Install dependencies
make install-dependencies
```

---

## 🔍 **Troubleshooting**

### **Common Issues:**

#### **Husky not working:**

```bash
# Reinstall Husky
npm run clean
npm install
npx husky init
```

#### **Commits rejected:**

```bash
# Verify commit format
npm run validate:commits

# Use Commitizen for correct commits
npm run commit
```

#### **Tests failing:**

```bash
# Verify databases are running
docker-compose ps

# Run tests individually
go test -v ./internal/services/orderservice/...
```

#### **Linting failing:**

```bash
# Format code automatically
go fmt ./...
goimports -w .

# Verify golangci-lint configuration
golangci-lint run --help
```

---

## 📖 **Additional Resources**

### **Documentation:**

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Husky Documentation](https://typicode.github.io/husky/)
- [Commitlint Rules](https://commitlint.js.org/#/reference-rules)
- [Standard Version](https://github.com/conventional-changelog/standard-version)

### **Go Tools:**

- [golangci-lint](https://golangci-lint.run/)
- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
- [Air](https://github.com/cosmtrek/air)

### **Architecture Patterns:**

- [Vertical Slice Architecture](https://jimmybogard.com/vertical-slice-architecture/)
- [CQRS Pattern](https://martinfowler.com/bliki/CQRS.html)
- [Event Sourcing](https://martinfowler.com/eaaDev/EventSourcing.html)

---

## 🎯 **Next Steps**

1. **Configure IDE** with Go tools
2. **Create first feature** following VSA
3. **Implement tests** for the feature
4. **Create commit** using Conventional Commits
5. **Push and verify** CI/CD
6. **Generate release** when ready

Happy coding! 🚀
