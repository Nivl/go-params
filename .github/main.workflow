workflow "Check code" {
  resolves = [
    "lint",
    "test",
  ]
  on = "push"
}

action "lint" {
  uses = "cedrickring/golang-action@1.3.0"
  args = "./tools/lint.sh"
  env = {
    GO111MODULE = "on"
    GOFLAGS = "-mod=readonly"
  }
}

action "test" {
  uses = "cedrickring/golang-action@1.3.0"
  args = "./tools/test.sh"
  env = {
    GO111MODULE = "on"
    GOFLAGS = "-mod=readonly"
    CI = "on"
  }
  secrets = ["CODECOV_TOKEN"]
}
