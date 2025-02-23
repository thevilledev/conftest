package main

deny[msg] {
    input.kind == "Deployment"
    not input.metadata.labels.app
    msg = "Deployments must have an 'app' label"
}

deny[msg] {
    input.kind == "Deployment"
    not input.metadata.labels.environment
    msg = "Deployments must have an 'environment' label"
} 