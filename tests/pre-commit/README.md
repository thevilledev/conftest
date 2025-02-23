# Pre-commit Hook Tests

These tests verify the functionality of Conftest's pre-commit hook integration.

## Test Cases

1. Hook Installation
   - Verifies that the pre-commit hook can be installed successfully

2. Basic Policy Validation
   - Tests single policy validation using the basic example
   - Verifies that missing required labels are detected

3. Multi-policy Validation
   - Tests the combine feature with multiple policies
   - Verifies that all policy violations are detected:
     - Naming conventions
     - Docker security
     - Kubernetes requirements

4. Compliant Configuration
   - Tests that a fully compliant configuration passes all checks

## Running Tests

The tests are automatically run as part of the project's CI pipeline. To run them locally:

```bash
bats tests/pre-commit/test.bats
```

Note: Requires pre-commit to be installed (`pip install pre-commit`) 