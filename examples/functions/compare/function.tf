terraform {
  required_providers {
    semvers = {
      source = "anapsix/semvers"
    }
  }
}

provider "semvers" {}

output "semvers_compare_results" {
  value = [
    {
      expected  = "equals (0)"
      arguments = "0.1.1, 0.1.1"
      result    = provider::semvers::compare("0.1.1", "0.1.1")
    },
    {
      expected  = "equals (0)"
      arguments = "0.1.0, 0.1.0+test"
      result    = provider::semvers::compare("0.1.0", "0.1.0+test")
    },
    {
      expected  = "lesser (-1)"
      arguments = "0.1.0, 0.1.1"
      result    = provider::semvers::compare("0.1.0", "0.1.1")
    },
    {
      expected  = "greater (1)"
      arguments = "0.1.1, 0.1.0"
      result    = provider::semvers::compare("0.1.1", "0.1.0")
    },
  ]
}
