package marketnexus.security

deny[msg] {
  input.kind == "Deployment"
  some i
  c := input.spec.template.spec.containers[i]
  not c.securityContext.runAsNonRoot
  msg := sprintf("container %s must set runAsNonRoot", [c.name])
}

deny[msg] {
  input.kind == "Deployment"
  some i
  c := input.spec.template.spec.containers[i]
  not c.securityContext.readOnlyRootFilesystem
  msg := sprintf("container %s must set readOnlyRootFilesystem", [c.name])
}

deny[msg] {
  input.kind == "Deployment"
  some i
  c := input.spec.template.spec.containers[i]
  not c.securityContext.capabilities.drop[_] == "ALL"
  msg := sprintf("container %s must drop all capabilities", [c.name])
}
