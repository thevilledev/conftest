# Multi-policy Pre-commit Example

This example demonstrates how to use Conftest with pre-commit to validate configurations against multiple policies simultaneously using the combine feature.

## Policy Description

The example includes three policy files that work together:

1. `policy/kubernetes.rego`: Kubernetes-specific policies
   - Requires `app` and `environment` labels
   - Enforces resource limits
   - Requires security context

2. `policy/docker.rego`: Container image policies
   - Prohibits use of `:latest` tag
   - Enforces approved registry usage
   - Prevents privileged containers

3. `policy/naming.rego`: Common naming conventions
   - Enforces lowercase alphanumeric names with hyphens
   - Validates namespace naming
   - Checks label naming format

## Files
- `policy/*.rego`: The Rego policy files containing various rules
- `configs/deployment.yaml`: An example Kubernetes Deployment that violates multiple policies

## Setup

1. Install pre-commit:
```bash
pip install pre-commit
```

2. Create `.pre-commit-config.yaml` in your repository:
```yaml
repos:
  - repo: https://github.com/open-policy-agent/conftest
    rev: v0.58.0  # Use the desired version
    hooks:
      - id: conftest
        args: 
          - --policy
          - examples/pre-commit/multi-policy/policy
          - --combine
        files: configs/.*\.(yaml|yml)$
```

3. Install the pre-commit hook:
```bash
pre-commit install
```

## Testing

The example deployment will fail multiple policy checks. To see this in action:

```bash
pre-commit run conftest --files configs/deployment.yaml
```

Expected violations:
1. Invalid resource name format (uppercase letters)
2. Invalid namespace name (underscore)
3. Invalid label name (uppercase letters)
4. Use of latest tag
5. Missing resource limits
6. Privileged container
7. Non-approved registry

To fix the deployment, you would need to:
1. Use lowercase hyphenated names
2. Add resource limits
3. Use a specific version tag from approved registry
4. Remove privileged access
5. Add required labels

Example of a compliant deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  namespace: prod-env
  labels:
    app: my-app
    environment: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
        environment: production
    spec:
      containers:
      - name: my-app
        image: registry.company.com/my-app:1.0.0
        ports:
        - containerPort: 80
        resources:
          limits:
            cpu: "1"
            memory: "1Gi"
          requests:
            cpu: "200m"
            memory: "256Mi"
        securityContext:
          privileged: false
          readOnlyRootFilesystem: true
``` 