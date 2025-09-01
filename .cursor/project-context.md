# 🚀 Go Food Delivery Microservices - Project Context

## 📋 **Project Overview**

`Go Food Delivery Microservices` is an imaginary and practical food delivery microservices application, built with Golang and implementing various software architecture patterns and technologies. This project serves as a **template for building backend microservice projects in Go**.

## 🏗️ **System Architecture**

### **High-Level Architecture Patterns**

- **Microservices Architecture**: Independent, loosely coupled services
- **Vertical Slice Architecture**: Feature-based organization instead of technical layers
- **Event-Driven Architecture**: Asynchronous communication via RabbitMQ
- **CQRS Pattern**: Command Query Responsibility Segregation
- **Domain-Driven Design (DDD)**: Business logic centered design
- **Event Sourcing**: Audit trail and state reconstruction
- **Dependency Injection**: Using Uber FX framework

### **Service Architecture**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Orders Service │    │ Catalogs Write  │    │  Catalogs Read  │
│                 │    │    Service      │    │    Service      │
│ • Event Sourcing│    │ • CRUD Operations│   │ • Read Models   │
│ • CQRS          │    │ • Domain Logic  │    │ • Search & Query│
│ • Audit Based   │    │ • Event Store   │    │ • Elasticsearch │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### **Communication Patterns**

- **Asynchronous**: RabbitMQ for event-driven communication
- **Synchronous**: REST and gRPC for real-time operations
- **Event Bus**: Custom implementation for message routing
- **Service Discovery**: Internal service communication

## 🎯 **Business Domain**

### **Core Entities**

- **Order**: Food delivery orders with lifecycle management
- **Product**: Catalog items with categories and pricing
- **Customer**: User profiles and preferences
- **Restaurant**: Vendor information and menus
- **Driver**: Delivery personnel management

### **Business Use Cases**

- Order creation and lifecycle management
- Product catalog management
- Customer profile management
- Restaurant operations
- Delivery tracking and management

## 🛠️ **Technology Stack**

### **Core Framework & Language**

- **Go 1.22+**: Modern Go with modules support
- **Uber FX**: Dependency injection framework
- **Vertical Slice Architecture**: Feature-based organization
- **CQRS**: Command/Query separation pattern

### **Data Storage**

- **PostgreSQL**: Primary write database with ACID transactions
- **EventStoreDB**: Event sourcing and audit trail
- **MongoDB**: Read models and document storage
- **Elasticsearch**: Search and analytics
- **Redis**: Caching and session management

### **Messaging & Events**

- **RabbitMQ**: Message broker for event-driven architecture
- **Custom Event Bus**: Message routing and handling
- **Event Sourcing**: Complete change history
- **Domain Events**: Business event modeling

### **API & Communication**

- **Echo Framework**: High-performance HTTP framework
- **gRPC**: Internal service communication
- **Swagger**: API documentation with swaggo
- **Validation**: go-playground/validator and ozzo-validation

### **Observability & Monitoring**

- **OpenTelemetry**: Distributed tracing and metrics
- **Jaeger/Zipkin**: Trace visualization
- **Prometheus**: Metrics collection
- **Grafana**: Monitoring dashboards
- **Zap**: Structured logging

### **Testing & Quality**

- **Testify**: Testing framework with mocks
- **Mockery**: Mock generation
- **Testcontainers**: Integration testing with Docker
- **golangci-lint**: Code quality and linting

## 📁 **Code Organization**

### **Vertical Slice Architecture**

- **Feature-based organization**: Each business feature is a vertical slice
- **Minimal coupling**: Slices are independent of each other
- **Technical folders within slices**: Controllers, handlers, models per feature
- **CQRS integration**: Commands and queries organized by feature

### **CQRS Implementation**

- **Commands**: Write operations with business logic
- **Queries**: Read operations optimized for performance
- **Handlers**: Separate command and query handlers
- **Mediator Pattern**: Using Go-MediatR for operation routing

## 🔄 **Data Flow Patterns**

### **Write Side (Commands)**

1. **Command** → **Command Handler** → **Domain Logic**
2. **Domain Events** → **Event Store** → **Event Bus**
3. **Projections** → **Read Models** (MongoDB/Elasticsearch)

### **Read Side (Queries)**

1. **Query** → **Query Handler** → **Read Repository**
2. **Read Models** → **Response DTOs** → **Client**

### **Event Flow**

1. **Domain Event** → **Event Store** (EventStoreDB)
2. **Event Bus** → **Event Handlers** (RabbitMQ)
3. **Projections** → **Read Models** (MongoDB/Elasticsearch)
4. **Synchronization** → **Other Services**

## 🧪 **Testing Strategy**

### **Testing Levels**

- **Unit Tests**: Individual components with mocked dependencies
- **Integration Tests**: Service integration with real dependencies
- **End-to-End Tests**: Complete workflow testing
- **Performance Tests**: Load and stress testing

### **Testing Tools**

- **Testify**: Testing framework and assertions
- **Mockery**: Mock generation for interfaces
- **Testcontainers**: Docker-based test dependencies
- **Coverage**: Code coverage analysis

## 🚀 **Deployment & DevOps**

### **Containerization**

- **Docker**: Service containerization
- **Docker Compose**: Local development environment
- **Multi-stage builds**: Optimized production images

### **CI/CD Pipeline**

- **GitHub Actions**: Automated build and test
- **Conventional Commits**: Standardized commit messages
- **Husky**: Git hooks for quality gates
- **Commitlint**: Commit message validation

### **Development Tools**

- **Air**: Live reloading for development
- **golangci-lint**: Code quality enforcement
- **gofumpt**: Advanced Go formatting
- **goimports-reviser**: Import organization

## 📚 **Development Standards**

### **Code Quality**

- **Go Modules**: Dependency management
- **Conventional Commits**: Standardized commit messages
- **golangci-lint**: Comprehensive linting rules
- **Code formatting**: Automated with pre-commit hooks

### **API Design**

- **RESTful principles**: Standard HTTP methods and status codes
- **OpenAPI**: Swagger documentation
- **Validation**: Input validation and error handling
- **Versioning**: API versioning strategy

### **Error Handling**

- **Structured errors**: Using emperror/errors package
- **Error wrapping**: Context preservation
- **Logging**: Structured logging with Zap
- **Monitoring**: Error tracking and alerting

## 🔮 **Future Roadmap**

### **Planned Features**

- **Kubernetes deployment**: Helm charts and K8s manifests
- **Outbox Pattern**: Guaranteed message delivery
- **Inbox Pattern**: Idempotency and exactly-once delivery
- **Advanced DDD**: Enhanced domain modeling

### **Technical Improvements**

- **GraphQL**: Flexible query language
- **WebSockets**: Real-time communication
- **Circuit Breaker**: Resilience patterns
- **Rate Limiting**: API protection

## 💡 **AI Development Guidelines**

### **When Generating Code**

- Follow **Vertical Slice Architecture** principles
- Implement **CQRS pattern** for operations
- Use **Event Sourcing** for audit-based services
- Apply **DDD principles** for domain logic
- Include **proper error handling** and logging
- Write **comprehensive tests** with mocking

### **When Reviewing Code**

- Verify **architecture compliance** with patterns
- Check **error handling** and logging
- Validate **testing coverage** and quality
- Review **performance considerations**
- Ensure **observability** and monitoring

### **When Suggesting Architecture**

- Maintain **service independence**
- Consider **scalability** and performance
- Plan for **resilience** and fault tolerance
- Include **monitoring** and alerting
- Follow **microservices best practices**

## 🎯 **Project Goals**

This project demonstrates:

- **Modern Go development** practices
- **Enterprise architecture** patterns
- **Microservices design** principles
- **Event-driven systems** implementation
- **Testing strategies** for complex systems
- **DevOps practices** and tooling

**Use this project as a template for building your own Go microservices applications!**
