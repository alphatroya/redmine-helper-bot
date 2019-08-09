workflow "Test" {
  on = "push"
  resolves = ["continuous-integration-workflow"]
}

action "continuous-integration-workflow" {
  uses = "./workflows/continuous-integration-workflow.yml"
}
