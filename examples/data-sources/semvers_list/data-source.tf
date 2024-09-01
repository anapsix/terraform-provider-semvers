data "semvers_list" "example" {
  provider = semvers
  list = [
    "0.90.1",
    "0.9.0",
    "0.80.0",
    "0.91.0",
    "0.9.10",
    "1.0.1",
    "2.0.0",
    "2.0.0-rc1",
    "2.0.1-rc1",
    "v0.1",
    "5.0.1-rc1+dead"
  ]
}
