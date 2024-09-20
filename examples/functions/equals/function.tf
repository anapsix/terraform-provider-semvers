terraform {
  required_providers {
    semvers = {
      source = "anapsix/semvers"
    }
  }
}

provider "semvers" {}

output "semvers_equals_results" {
  value = [
    {
      expected  = "true"
      arguments = "0.1.1, 0.1.1"
      result    = provider::semvers::equals("0.1.1", "0.1.1")
    },
    {
      expected  = "true"
      arguments = "0.1.0, 0.1.0+test"
      result    = provider::semvers::equals("0.1.0", "0.1.0+test")
    },
    {
      expected  = "false"
      arguments = "0.1.0, 0.1.1"
      result    = provider::semvers::equals("0.1.0", "0.1.1")
    },
    {
      expected  = "false"
      arguments = "0.1.1, 0.1.0"
      result    = provider::semvers::equals("0.1.1", "0.1.0")
    },
  ]
}
