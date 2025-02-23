#!/usr/bin/env bats

DIR="$( cd "$( dirname "${BATS_TEST_FILENAME}" )" >/dev/null 2>&1 && pwd )"

setup_file() {
    cd "$DIR/../.."
    # Verify pre-commit is installed and in PATH
    which pre-commit
    pre-commit --version
    cat > .pre-commit-config.yaml << 'EOF'
repos:
- repo: .
  rev: HEAD
  hooks:
    - id: conftest
      args:
        - --policy
        - examples/pre-commit/basic/policy
EOF
}

@test "Can install pre-commit hook" {
    cd "$DIR/../.."
    run pre-commit try-repo . conftest --verbose
    [ "$status" -eq 0 ]
}

@test "pre-commit: basic policy validation fails as expected" {
    cd "$DIR/../.."
    run pre-commit run conftest --files examples/pre-commit/basic/deployment.yaml
    [ "$status" -eq 1 ]
    [[ "$output" =~ "Deployments must have an 'app' label" ]]
    [[ "$output" =~ "Deployments must have an 'environment' label" ]]
}