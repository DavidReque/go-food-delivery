module.exports = {
  extends: ["@commitlint/config-conventional"],
  rules: {
    // Type validation
    "type-enum": [
      2,
      "always",
      [
        "feat", // New features
        "fix", // Bug fixes
        "docs", // Documentation changes
        "style", // Code style changes (formatting, missing semicolons, etc)
        "refactor", // Code refactoring
        "perf", // Performance improvements
        "test", // Adding or updating tests
        "build", // Build system or external dependencies
        "ci", // CI/CD changes
        "chore", // Other changes that don't modify src or test files
        "revert", // Revert previous commits
        "breaking", // Breaking changes
      ],
    ],

    // Scope validation for our Go Food Delivery project
    "scope-enum": [
      2,
      "always",
      [
        // Services
        "orders", // Orders service
        "catalogs", // Catalogs service
        "customers", // Customers service
        "payments", // Payments service
        "delivery", // Delivery service
        "auth", // Authentication service

        // Technical layers
        "api", // API layer
        "core", // Core domain
        "infrastructure", // Infrastructure layer
        "testing", // Testing framework

        // Documentation and tools
        "docs", // Documentation
        "ci", // CI/CD
        "docker", // Docker configuration
        "k8s", // Kubernetes configuration

        // General
        "go-food-delivery", // General project changes
      ],
    ],

    // Subject validation
    "subject-case": [2, "always", "lower-case"],
    "subject-empty": [2, "never"],
    "subject-full-stop": [2, "never", "."],
    "header-max-length": [2, "always", 72],

    // Body validation
    "body-leading-blank": [2, "always"],
    "body-max-line-length": [2, "always", 100],

    // Footer validation
    "footer-leading-blank": [2, "always"],
    "footer-max-line-length": [2, "always", 100],
  },
};
