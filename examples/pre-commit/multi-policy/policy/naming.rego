package naming

deny[msg] {
    not regex.match("^[a-z0-9][a-z0-9-]*[a-z0-9]$", input.metadata.name)
    msg = sprintf("Resource name '%v' must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character", [input.metadata.name])
}

deny[msg] {
    not regex.match("^[a-z0-9][a-z0-9-]*[a-z0-9]$", input.metadata.namespace)
    msg = sprintf("Namespace '%v' must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character", [input.metadata.namespace])
}

deny[msg] {
    label_name := [name | name := input.metadata.labels[_]; not regex.match("^[a-z0-9][a-z0-9-\\.]*[a-z0-9]$", name)][_]
    msg = sprintf("Label name '%v' must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character", [label_name])
} 