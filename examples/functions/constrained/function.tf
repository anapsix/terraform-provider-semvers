terraform {
  required_providers {
    semvers = {
      source = "anapsix/semvers"
    }
  }
}

provider "semvers" {}

output "semvers_constrained_results" {
  value = [
    {
      expected  = true
      arguments = "0.1.1, >= 0.1"
      result    = provider::semvers::constrained("0.1.1", ">= 0.1")
    },
    {
      expected  = false
      arguments = "0.1, >= 0.1.1"
      result    = provider::semvers::constrained("0.1", ">= 0.1.1")
    },
    {
      expected  = true
      arguments = "0.1.2, ~> 0.1"
      result    = provider::semvers::constrained("0.1.2", "~> 0.1")
    },
    {
      expected  = false
      arguments = "0.2, ~> 0.1"
      result    = provider::semvers::constrained("0.2", "~> 0.1")
    },
    {
      expected  = true
      arguments = "0.2, 0.1.1 - 0.2.0"
      result    = provider::semvers::constrained("0.2", "0.1.1 - 0.2.0")
    },
  ]
}
