# Terraform Provider Semvers

Implements a data-source `semvers_list` to make semver sorting easy in TF, and
provider functions (TF >= 1.8 is required):
- `data.semvers_list`: sorts a list of semver strings
- `provider::semvers::compare`: compares two semver strings
- `provider::semvers::constrained`: checks whether semver is within constrains
- `provider::semvers::equals`: checks two semver strings for equality
- `provider::semvers::pick`: filters a list of semver strings by constraint
- `provider::semvers::sort`: sorts a list of semver strings, returns sorted one

Version strings are parsed leniently: a leading `v` and missing minor or patch parts are accepted, so `v1.2` is read as `1.2.0`. Build metadata is ignored when comparing and filtering. `sort` and `pick` return normalized versions (`"1.0.0"`, not `"v1.0.0"`) and drop duplicates (e.g. `v1.0.0`, `1.0`, and `1.0.0`); `semvers_list` keeps the first occurrence's input string in `original`. Versions that differ only in build metadata (e.g. `1.0.0+a` and `1.0.0+b`) are both kept, in input order.

See the [Terraform Registry provider page][1].

## Provider configuration

```hcl
provider "semvers" {
  ignore_invalid_tags = false # default; set true to skip invalid semver strings in `semvers_list`
}
```

When `ignore_invalid_tags` is `true`, the `semvers_list` data source skips invalid semver strings instead of returning an error. Skipped tags are logged at **INFO** when `TF_LOG` is enabled:

```sh
TF_LOG=INFO terraform plan
```

Use `provider = semvers.<alias>` on the data source to pick a strict or lenient configuration.

### Skipping invalid versions in provider functions

Terraform and OpenTofu never pass provider configuration to provider functions, so `ignore_invalid_tags` has **no effect** on them. Instead, `sort` and `pick` take an optional trailing `ignore_invalid` argument (default `false`):

```hcl
# pick the latest matching dev tag from a mixed image tag list;
# try() returns null instead of failing when nothing matches
try(reverse(provider::semvers::pick(
  ["develop-latest", "1.2.3-dev.1", "1.2.3-dev.2", "0.1.0"],
  "~> 1.2.3-dev",
  true,
))[0], null)
# => "1.2.3-dev.2"

provider::semvers::sort(["develop-latest", "0.2.1", "0.1.0"], true)
# => ["0.1.0", "0.2.1"]
```

Single-version functions (`compare`, `equals`, `constrained`) always return an error when a version string is invalid.

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
