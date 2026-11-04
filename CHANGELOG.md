## Unreleased

### Breaking
- The provider `ignore_invalid_tags` setting no longer claims to apply to `provider::semvers::sort` and `provider::semvers::pick`; pass the new `ignore_invalid` argument instead, e.g. `provider::semvers::pick(list, "~> 1.2", true)` ([#3])

### Added
- Optional trailing `ignore_invalid` boolean argument (default `false`) for `provider::semvers::sort` and `provider::semvers::pick` ([#3])
- CI workflow running fmt, vet, build, and acceptance tests against Terraform 1.8, 1.12, and 1.17 on pull requests and pushes to `main` ([#3])
- CI runs golangci-lint (same version as pre-commit) and fails when `go generate` output (docs, formatted examples) is not committed ([#3])

### Fixed
- `provider::semvers::sort` and `provider::semvers::pick` failed with an internal "Value Conversion Error" when the list contained a `null` element. Null elements are now treated as invalid versions: an error by default, skipped with `ignore_invalid`, matching `semvers_list` ([#3])
- `semvers_list` silently parsed strictly when the provider `ignore_invalid_tags` value was not known until apply (e.g. derived from a resource attribute), failing plan with a misleading parse error. It now returns an error explaining that the setting must be known at plan time; lists containing only valid strings are unaffected ([#3])
- `provider::semvers::sort` and `provider::semvers::pick` ignored the provider `ignore_invalid_tags` setting under real Terraform/OpenTofu, because provider-defined functions never receive provider configuration ([#3])
- Duplicate versions in `semvers_list`, `sort`, and `pick` now keep the first occurrence from the input (stable sort); previously the kept `original` could be any of the duplicates ([#3])
- `provider::semvers::compare` docs no longer claim it returns `99` on error; it returns a function error ([#3])

### Changed
- Function errors now point at the offending argument and include the bad value, e.g. `Invalid value for "version_b" parameter: invalid version "blah": invalid semantic version` ([#3])
- Provider server address for `-debug` mode is now `registry.terraform.io/anapsix/semvers` ([#3])
- Rename internal `helpers` package to `semverutil`; camelCase variable names; drop redundant parameter defaults and dead code ([#3])
- Remove HashiCorp copyright headers from files not derived from the provider scaffolding template ([#3])
- Tests use real `provider` blocks instead of injecting provider config ([#3])

### Documentation
- `ignore_invalid_tags` documented as applying to the `semvers_list` data source only; README section on provider aliases and functions rewritten ([#3])
- Document lenient version parsing (`v` prefix, missing minor/patch parts) ([#3])
- Add `ignore_invalid` examples for `sort` and `pick`; fix description typos ([#3])

[#3]: https://github.com/anapsix/terraform-provider-semvers/pull/3

## 0.8.1

### Changed
- Migrate `.goreleaser.yaml` from deprecated `archives.format` to `formats`
- Bump Go directive 1.25.8 -> 1.26.0 (required by updated dependencies)
- Upgrade terraform-plugin-docs v0.19.4 -> v0.25.0
- Upgrade terraform-plugin-log v0.10.0 -> v0.11.0
- Upgrade terraform-plugin-testing v1.10.0 -> v1.16.0
- Update indirect dependencies (`golang.org/x/*`, grpc, protobuf, and others)

### Documentation
- Regenerate docs with terraform-plugin-docs v0.25.0

## 0.8.0

### Added
- Optional `ignore_invalid_tags` provider attribute (default `false`) to skip invalid semver strings in list inputs ([#2])
- INFO-level logging of skipped tags when `TF_LOG` is enabled

### Changed
- `semvers_list`, `provider::semvers::pick`, and `provider::semvers::sort` respect `ignore_invalid_tags`
- Data source returns an error when no valid semver strings remain after parsing
- Bump Go toolchain to 1.25.8
- Upgrade terraform-plugin-framework, terraform-plugin-go, terraform-plugin-sdk/v2, and related HashiCorp dependencies
- Upgrade Masterminds/semver to v3.5.0
- Bump pre-commit-hooks v5.0.0 -> v6.0.0
- Bump golangci-lint v1.61.0 -> v2.12.2
- Bump release workflow actions (SHA-pinned): `actions/setup-go` v5 -> v6.5.0, `crazy-max/ghaction-import-gpg` v6 -> v7.0.0, `goreleaser/goreleaser-action` v6.2.1 -> v7.2.3

### Documentation
- Provider configuration, provider alias limitations with functions, and OpenTofu differences

[#2]: https://github.com/anapsix/terraform-provider-semvers/issues/2

## 0.7.1
- updating go modules
- bumping `goreleaser/goreleaser-action` v6.0.0 -> v6.2.1
- bumping `go` 1.23.4 -> 1.24.1

## 0.7.0

### Changed
- bumping the `terraform-plugin-framework` v1.13.0 -> v1.14.0, supports TF 1.11
- udpating other go modules

## 0.6.2

### Changed
- udpating go modules

## 0.6.1

## Changed
- updating golang.org/x/net v0.32.0 -> v0.33.0 ([CVE-2024-45338])

[CVE-2024-45338]: https://github.com/advisories/GHSA-w32m-9786-jp63

## 0.6.0

### Added
- The `provider::semvers::constrained` function for checking semver constrains

### Changed
- updating google.golang.org/grpc v1.68.1 -> v1.69.0

## 0.5.5

### Changed
- updating golang.org/x/crypto v0.30.0 -> v0.31.0

### Changed
- updated go modules

## 0.5.4

### Changed
- updated go modules

## 0.5.3

### Changed
- updating go modules
- bumping go to 1.23.2
- bumping pre-commit repos

## 0.5.2

### Fixed
- mentioning `equals` function in README

## 0.5.1

### Added
- The `provider::semvers::equals` function for checking equality of two semver

### Fixed
- spacing consistency
- go modules update
- improving functions docs

## 0.5.0

### Added
- The `provider::semvers::compare` function for comparing semver strings
- Tests covering invalid values

## 0.4.2

### Fixed
- provider function examples and docs

## 0.4.1

### Added
- The `provider::semvers::pick` function examples

## 0.4.0

### Added
- The `provider::semvers::pick` function which takes list of semver strings,
  and semver constraint, and returns a list of filtered semver strings, sorted and deduped,
  matching the constraint. See [Masterminds/semver](https://github.com/Masterminds/semver/tree/master?tab=readme-ov-file#checking-version-constraints) for constraint syntax.

## 0.3.1

### Fixed

- The changelog formatting
- The provider and function docs
- The function examples

## 0.3.0

### Added
- The `provider::semvers::sort` function which takes list of semver strings,
  and returns a list of semver strings, sorted and deduped

## 0.2.1

### Added
- Deduplication based on parsed semver version string

## 0.2.0

### Added
- The `original` attribute to `sorted_versions` object with passed semver string
- The descriptions to data-source attributes

## 0.1.0

- Initial release
