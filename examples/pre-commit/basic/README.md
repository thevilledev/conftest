# Basic Pre-commit Example

This example demonstrates how to use Conftest with pre-commit to validate Kubernetes manifests against a simple policy.

## Policy Description

The policy in `policy/kubernetes.rego` enforces two rules for Kubernetes Deployments:
1. All Deployments must have an `app` label
2. All Deployments must have an `environment` label

## Files
- `policy/kubernetes.rego`: The Rego policy file containing the rules
- `deployment.yaml`: An example Kubernetes Deployment that fails the policy checks

## Setup

1. Install pre-commit:
```bash
pip install pre-commit
```

2. Create `.pre-commit-config.yaml` in your repository:
```yaml
repos:
  - repo: https://github.com/open-policy-agent/conftest
    rev: v0.45.0  # Use the desired version
    hooks:
      - id: conftest
        args: [--policy, examples/pre-commit/basic/policy]
        files: deployment\.yaml$
```

3. Install the pre-commit hook:
```bash
pre-commit install
```

## Testing

The example deployment will fail the policy check because it's missing the required `environment` label. To see this in action:

```bash
pre-commit run conftest --files deployment.yaml
```

To fix the deployment, add the required labels to the metadata:

```yaml
metadata:
  name: example-app
  labels:
    app: example-app
    environment: production
``` 