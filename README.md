# Terraform Provider Semvers

Implements a data-source `semvers_list` to make semver sorting easy in TF.


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
