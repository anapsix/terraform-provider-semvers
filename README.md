# Terraform Provider Semvers

Implements a data-source `semvers_list` to make semver sorting easy in TF, and
provider functions (TF >= 1.8 is required):
- `data.semvers_list`: sorts a list of semver strings
- `provider::semvers::compare`: compares two semver strings
- `provider::semvers::constrained`: checks whether semver is within constrains
- `provider::semvers::equals`: checks two semver strings for equality
- `provider::semvers::pick`: filters a list of semver strings by constraint
- `provider::semvers::sort`: sorts a list of semver strings, returns sorted one

See the [Terraform Registry provider page][1].

## Provider configuration

```hcl
provider "semvers" {
  ignore_invalid_tags = false # default; set true to skip invalid semver strings in list inputs
}
```

When `ignore_invalid_tags` is `true`, invalid semver strings in list inputs are skipped instead of returning an error. This applies to `semvers_list`, `sort`, and `pick`. Skipped tags are logged at **INFO** when `TF_LOG` is enabled:

```sh
TF_LOG=INFO terraform plan
```

Example — pick the latest matching dev tag from a mixed image tag list:

```hcl
provider "semvers" {
  ignore_invalid_tags = true
}

reverse(provider::semvers::pick(
  ["develop-latest", "1.2.3-dev.1", "1.2.3-dev.2", "0.1.0"],
  "~> 1.2.3-dev",
))[0]
# => "1.2.3-dev.2"
```

Single-version functions (`compare`, `equals`, `constrained`) are unchanged and still return an error when a version string is invalid.

### Provider aliases and provider functions

Provider functions are called as `provider::semvers::<function>(...)`. In **Terraform**, that syntax always uses the **default** (non-aliased) provider configuration. There is no way to select an aliased provider for a function call — aliases only apply to resources and data sources via the `provider` meta-argument:

```hcl
provider "semvers" {
  ignore_invalid_tags = false # default
}

provider "semvers" {
  alias               = "lenient"
  ignore_invalid_tags = true
}

# Uses semvers.lenient — alias selection works
data "semvers_list" "example" {
  provider = semvers.lenient
  list     = ["1.0.0", "develop-latest"]
}

# Always uses the default semvers provider — alias is ignored
provider::semvers::pick(["1.0.0", "develop-latest"], ">= 1.0")
```

If you need different provider settings (such as `ignore_invalid_tags`) in the same module under Terraform, use a data source with `provider = semvers.<alias>`, or split the work into child modules with different default providers passed via `providers`.

**OpenTofu** (1.7+) supports alias-aware provider functions:

```hcl
provider::semvers::lenient::pick(["1.0.0", "develop-latest"], ">= 1.0")
```

[1]: https://registry.terraform.io/providers/anapsix/semvers

## Development

```sh
# install dev version of the provider
go install

# test dev version of the provider
(cd ./test; terraform plan)

# run acceptance testing
TF_ACC=1 go test -v ./...

# check and fix formatting
go fmt ./...

# update go modules
go get -u

# generate docs
go generate ./...
```
