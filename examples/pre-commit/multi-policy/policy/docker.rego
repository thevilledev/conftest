package docker

deny[msg] {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    endswith(container.image, ":latest")
    msg = sprintf("Container '%v' uses the 'latest' tag which is not allowed", [container.name])
}

deny[msg] {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    not startswith(container.image, "registry.company.com/")
    msg = sprintf("Container '%v' must use an approved registry (registry.company.com)", [container.name])
}

deny[msg] {
    input.kind == "Deployment"
    container := input.spec.template.spec.containers[_]
    container.securityContext.privileged == true
    msg = sprintf("Container '%v' must not run as privileged", [container.name])
} 