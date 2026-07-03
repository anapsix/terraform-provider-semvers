## Unreleased

### Changed
- Migrate `.goreleaser.yaml` from deprecated `archives.format` to `formats`

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
