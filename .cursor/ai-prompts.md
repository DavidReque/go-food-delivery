# 🤖 AI Prompts for Cursor

## 🎯 **General Project Prompts**

### **Go Code Generation**

```
Generate Go code for [DESCRIPTION] following these patterns:
- Vertical Slice Architecture with feature-based organization
- CQRS pattern with separate command/query handlers
- Event Sourcing for audit-based services
- Domain-Driven Design principles
- Proper error handling with emperror/errors
- Structured logging with Zap
- Comprehensive testing with Testify and Mockery
- Use Uber FX for dependency injection
```

### **Code Review**

```
Review this Go code considering:
- Vertical Slice Architecture compliance
- CQRS pattern implementation
- Error handling and structured logging
- Testing coverage and mocking strategy
- Performance and resource management
- Observability and monitoring
- Microservices best practices
```

## 🏗️ **Architecture-Specific Prompts**

### **Vertical Slice Architecture**

```
Implement a vertical slice for [FEATURE] following:
- Feature-based organization instead of technical layers
- All technical concerns (controllers, handlers, models) in one slice
- Minimal coupling between slices
- CQRS integration within the slice
- Independent deployment capability
```

### **CQRS Implementation**

```
Implement CQRS for [ENTITY] following:
- Separate command and query models
- Command handlers with business logic
- Query handlers optimized for reading
- Event sourcing integration
- Mediator pattern using Go-MediatR
- Performance optimization for read side
```

### **Event Sourcing**

```
Create event sourcing for [ACTION] following:
- Domain events with proper structure
- Event store integration with EventStoreDB
- Event versioning strategy
- Projection updates to read models
- Event replay capabilities
- Audit trail maintenance
```

## 📁 **Layer-Specific Prompts**

### **Domain Layer**

```
Generate domain logic for [ENTITY] following DDD:
- Aggregate root with business methods
- Value objects with immutability
- Domain events with business meaning
- Business rule validation
- Domain service implementation
- Rich domain model design
```

### **Application Layer**

```
Implement application layer for [OPERATION]:
- Command/Query with validation
- Handler with dependency injection
- Business logic orchestration
- Transaction management
- Event publishing
- Integration with domain layer
```

### **Infrastructure Layer**

```
Create infrastructure for [COMPONENT]:
- Repository pattern implementation
- Database integration (PostgreSQL/MongoDB)
- Event store configuration
- Message bus setup (RabbitMQ)
- Caching strategy (Redis)
- Health checks and monitoring
```

### **API Layer**

```
Design API endpoint for [OPERATION]:
- RESTful design with Echo framework
- Request/Response DTOs with validation
- Swagger documentation
- Error handling and status codes
- Authentication and authorization
- Rate limiting and throttling
```

## 🔄 **Pattern Implementation Prompts**

### **Repository Pattern**

```
Implement repository for [ENTITY]:
- Generic repository interface
- CRUD operations with transactions
- Specification pattern support
- Connection pooling and management
- Error handling and logging
- Testing with mocks
```

### **Mediator Pattern**

```
Configure mediator for [SERVICE]:
- Command/Query registration
- Handler discovery and routing
- Middleware pipeline setup
- Error handling and logging
- Performance monitoring
- Cross-cutting concerns
```

### **Event-Driven Architecture**

```
Setup event-driven communication for [SERVICE]:
- Event bus configuration
- RabbitMQ integration
- Event routing and handling
- Dead letter queue setup
- Event persistence and replay
- Monitoring and alerting
```

## 🧪 **Testing Strategy Prompts**

### **Unit Tests**

```
Generate unit tests for [COMPONENT]:
- Test cases covering happy path
- Edge cases and error scenarios
- Mocking of dependencies with Mockery
- Assertion best practices with Testify
- Test naming conventions
- Coverage analysis
```

### **Integration Tests**

```
Create integration tests for [FUNCTIONALITY]:
- Testcontainers setup for dependencies
- Real database integration
- Event bus testing
- End-to-end workflow validation
- Performance assertions
- Cleanup strategies
```

### **Test Data Management**

```
Generate test fixtures for [ENTITY]:
- Realistic test data scenarios
- Edge case data preparation
- Performance test data sets
- Cleanup utilities
- Reusable test helpers
- Data isolation strategies
```

## 🚀 **Performance & Optimization Prompts**

### **Query Optimization**

```
Optimize queries for [OPERATION]:
- MongoDB aggregation pipeline optimization
- Elasticsearch query performance
- Database indexing strategy
- Caching implementation
- Connection pooling
- Performance monitoring
```

### **Concurrency Management**

```
Implement concurrency for [OPERATION]:
- Goroutine management and lifecycle
- Channel usage patterns
- Mutex and RWMutex strategies
- Context cancellation handling
- Resource cleanup
- Race condition prevention
```

### **Caching Strategy**

```
Design caching for [DATA]:
- Cache invalidation strategies
- TTL and expiration policies
- Distributed caching with Redis
- Cache warming techniques
- Performance metrics
- Memory management
```

## 🔍 **Observability Prompts**

### **Logging Implementation**

```
Implement logging for [OPERATION]:
- Structured logging with Zap
- Log levels and correlation IDs
- Performance metrics logging
- Error context and stack traces
- Log aggregation and analysis
- Compliance and audit requirements
```

### **Metrics Collection**

```
Define metrics for [SERVICE]:
- Business metrics and KPIs
- Technical performance metrics
- Custom Prometheus metrics
- Alerting thresholds
- Dashboard configuration
- Trend analysis
```

### **Distributed Tracing**

```
Configure tracing for [OPERATION]:
- OpenTelemetry span creation
- Cross-service correlation
- Performance profiling
- Error tracking and debugging
- Trace visualization
- Sampling strategies
```

## 🐳 **DevOps & Deployment Prompts**

### **Docker Optimization**

```
Optimize Dockerfile for [SERVICE]:
- Multi-stage builds for efficiency
- Security best practices
- Layer optimization and caching
- Health checks and monitoring
- Resource limits and constraints
- Production readiness
```

### **Kubernetes Configuration**

```
Create K8s manifests for [SERVICE]:
- Deployment configuration
- Service networking and discovery
- Resource management and limits
- Health checks and readiness
- Scaling policies and HPA
- Security and RBAC
```

### **CI/CD Pipeline**

```
Design pipeline for [SERVICE]:
- GitHub Actions workflow
- Build optimization and caching
- Testing strategy and coverage
- Security scanning and validation
- Deployment automation
- Rollback procedures
```

## 📚 **Documentation Prompts**

### **API Documentation**

```
Generate documentation for [ENDPOINT]:
- OpenAPI specification
- Request/response examples
- Error codes and handling
- Authentication requirements
- Rate limiting information
- Usage examples
```

### **Architecture Documentation**

```
Document architecture for [COMPONENT]:
- Component diagrams and relationships
- Data flow and dependencies
- Configuration and environment
- Deployment and scaling
- Monitoring and alerting
- Security considerations
```

### **Code Documentation**

```
Document code for [FUNCTION]:
- Function purpose and behavior
- Parameters and return values
- Usage examples and patterns
- Error conditions and handling
- Performance considerations
- Dependencies and requirements
```

## 🎨 **API Design Prompts**

### **RESTful API Design**

```
Design REST API for [FUNCTIONALITY]:
- Resource naming conventions
- HTTP methods and status codes
- Request/response formats
- Error handling patterns
- Pagination and filtering
- Versioning strategy
```

### **Response Format Design**

```
Define response format for [OPERATION]:
- Consistent structure and naming
- Error handling and codes
- Pagination and metadata
- Filtering and sorting options
- Data transformation
- Performance optimization
```

## 🔧 **Configuration & Management Prompts**

### **Environment Configuration**

```
Configure environment for [SERVICE]:
- Development vs production settings
- Configuration validation
- Security considerations
- Default values and overrides
- Environment-specific behavior
- Configuration management
```

### **Feature Flags**

```
Implement feature flags for [FUNCTIONALITY]:
- Configuration management
- Runtime toggling
- A/B testing support
- Rollout strategies
- Monitoring and metrics
- Rollback capabilities
```

## 💡 **Improvement & Refactoring Prompts**

### **Code Refactoring**

```
Refactor this code following:
- Clean code principles
- SOLID principles
- Go best practices and idioms
- Performance optimization
- Maintainability improvements
- Testing enhancements
```

### **Architecture Review**

```
Review architecture considering:
- Scalability requirements
- Performance bottlenecks
- Security vulnerabilities
- Maintainability issues
- Technology debt
- Future growth plans
```

### **Migration Strategy**

```
Plan migration for [COMPONENT]:
- Backward compatibility
- Data migration strategy
- Rollback procedures
- Testing approach
- Deployment timeline
- Risk mitigation
```

## 🎯 **Project-Specific Prompts**

### **Microservices Communication**

```
Design communication for [SERVICES]:
- Event-driven messaging
- Synchronous API calls
- Service discovery
- Load balancing
- Circuit breaker patterns
- Monitoring and alerting
```

### **Event Sourcing Implementation**

```
Implement event sourcing for [DOMAIN]:
- Event store configuration
- Event structure and versioning
- Projection updates
- Event replay capabilities
- Performance optimization
- Monitoring and debugging
```

### **CQRS Pattern Application**

```
Apply CQRS to [BUSINESS_OPERATION]:
- Command and query separation
- Read/write model optimization
- Event sourcing integration
- Performance considerations
- Testing strategy
- Monitoring and metrics
```

## 🚨 **Troubleshooting Prompts**

### **Error Investigation**

```
Investigate error in [COMPONENT]:
- Error context and stack traces
- Log analysis and correlation
- Performance impact assessment
- Root cause identification
- Prevention strategies
- Monitoring improvements
```

### **Performance Issues**

```
Analyze performance issue in [OPERATION]:
- Bottleneck identification
- Resource utilization analysis
- Database query optimization
- Caching strategy review
- Scaling considerations
- Monitoring and alerting
```

### **Debugging Support**

```
Debug issue in [SERVICE]:
- Log analysis and correlation
- Distributed tracing review
- Error context examination
- Performance profiling
- Root cause analysis
- Solution implementation
```
