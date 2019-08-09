workflow "Test" {
  on = "push"
  resolves = ["continuous-integration-workflow"]
}

action "continuous-integration-workflow" {
  uses = "./continuous-integration-workflow.yml"
}
