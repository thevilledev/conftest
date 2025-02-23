package kubernetes

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

deny[msg] {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    not container.resources.limits
    msg = sprintf("Container '%v' must have resource limits", [container.name])
}

deny[msg] {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    not container.securityContext
    msg = sprintf("Container '%v' must have a security context defined", [container.name])
} 