# Terraform Provider Semvers

Implements a data-source `semvers_list` to make semver sorting easy in TF, and
provider functions (TF >= 1.8 is required):
- `data.semvers_list`: sorts a list of semver strings
- `provider::semvers::sort`: sorts a list of semver strings, returns sorted one
- `provider::semvers::pick`: filters a list of semver strings by constraint


See the [Terraform Registry provider page][1].

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
