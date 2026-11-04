# AGENTS.md

Instructions for AI coding agents (and humans) working on `terraform-provider-semvers`.

## Project overview

A Terraform provider built on [terraform-plugin-framework](https://developer.hashicorp.com/terraform/plugin/framework), published as [`anapsix/semvers`](https://registry.terraform.io/providers/anapsix/semvers). All semver logic comes from [`Masterminds/semver/v3`](https://github.com/Masterminds/semver).

- `semvers_list` data source: parses, sorts and dedups a list of versions
- Provider functions (Terraform >= 1.8): `compare`, `constrained`, `equals`, `pick`, `sort`
- Provider attribute `ignore_invalid_tags`: applies to the `semvers_list` data source only

| Path | Contents |
|------|----------|
| `main.go` | Provider server entry point, `go:generate` directives |
| `internal/provider/` | Provider, data source, functions (`function_<name>.go`) and their acceptance tests |
| `internal/semverutil/` | Parsing, sorting, filtering helpers with typed errors (`InvalidVersionError`, `InvalidConstraintError`) |
| `examples/` | HCL examples, embedded into the generated docs |
| `docs/` | **Generated** registry docs. Never edit by hand |
| `test/main.tf` | Manual end-to-end config for a locally installed build |
| `tools/tools.go` | Pins `tfplugindocs` for `go generate` |

## Things that are easy to get wrong

- **Provider functions never receive provider configuration.** Terraform and OpenTofu call functions on an unconfigured provider instance, so function behavior must come only from arguments. Never read provider state from a function. (`sort` and `pick` take an optional trailing `ignore_invalid` argument for this reason.)
- **Provider config can be unknown during plan.** `Configure` records whether `ignore_invalid_tags` is unknown, and `semvers_list` errors out instead of guessing. Handle unknowns the same way for any new attribute.
- **Parsing is lenient.** `semver.NewVersion` accepts a `v` prefix and missing minor/patch parts (`v1.2` becomes `1.2.0`). Dedup uses the normalized string, so `1.0.0+a` and `1.0.0+b` are both kept even though they compare equal.
- **Sorting must stay stable** (`sort.Stable`), so dedup keeps the first occurrence from the input.

## Development commands

```sh
go build ./...                  # build
go vet ./...                    # vet
gofmt -l .                      # must print nothing
go test ./...                   # unit tests (acceptance tests skip without TF_ACC)
TF_ACC=1 go test ./...          # acceptance tests; needs terraform >= 1.8 in PATH
TF_ACC=1 TF_ACC_TERRAFORM_VERSION=1.8.5 go test ./internal/provider/   # oldest supported Terraform
golangci-lint run ./...         # lint (see version note below)
go generate ./...               # format examples/ and test/, regenerate docs/
```

- **Run `go generate ./...` after changing any schema, description, parameter or example, and commit the result.** CI fails when generated files are out of date.
- golangci-lint is pinned to the same version in `.pre-commit-config.yaml` and `.github/workflows/test.yaml`; bump both together. Run it with the Go version from `go.mod` (e.g. `GOTOOLCHAIN=go1.26.10 golangci-lint run ./...`). Older golangci-lint builds fail with `export data version ... is greater than maximum supported version` under a newer local Go.
- Manual testing: `go install`, add a `dev_overrides` block for `anapsix-dev/semvers` pointing at `~/go/bin` (see the comment at the top of `test/main.tf`), then `cd test && terraform plan`. Check with both `terraform` and `tofu` when changing function or provider-config behavior.

## Code conventions

- **License headers:** every Go file starts with
  ```go
  // Copyright (c) Anastas Dancha
  // SPDX-License-Identifier: MPL-2.0
  ```
  Keep the existing `Copyright (c) HashiCorp, Inc.` line in files that are close copies of the HashiCorp scaffolding template (`tools/tools.go`, `internal/provider/provider_test.go`). MPL-2.0 does not allow removing copyright notices from copied code.
- Go identifiers are camelCase; Terraform-facing names (attributes, parameters) are snake_case.
- Parse with `semverutil.ParseVersion` / `semverutil.ParseConstraint`, not `semver.New*` directly, so errors carry the offending value.
- Function errors use `function.NewArgumentFuncError(position, err.Error())`, with each argument parsed separately so the position is known. Don't add generic "Error performing operation" wrappers.
- Leave framework defaults unset (e.g. no `AllowNullValue: false`).
- Exported identifiers get doc comments.

### Adding or changing a provider function

1. Add `internal/provider/function_<name>.go` and register it in `semversProvider.Functions`.
2. Add `internal/provider/function_<name>_test.go` covering known values, invalid input and null/wrong-type arguments.
3. Add `examples/functions/<name>/function.tf`.
4. Update the function list in `README.md` and the provider `Description` in `provider.go`.
5. Run `go generate ./...` and commit `docs/`.

### Tests

- Acceptance tests use `ProtoV6ProviderFactories: testAccProtoV6ProviderFactories` and configure the provider with **real `provider "semvers" {}` blocks** in the test config. Never inject config into the provider struct; that is how the functions bug in #3 went unnoticed.
- Function tests start with `tfversion.SkipBelow(tfversion.Version1_8_0)`.
- Terraform wraps long error lines, so `ExpectError` regexes should use `\s+` between words, e.g. `` `"version_b"\s+parameter:\s+invalid\s+version` ``.
- For a bug fix, write the failing test first.

## CHANGELOG

Update `CHANGELOG.md` in the same PR for every user-visible change (behavior, schema, docs, dependencies, CI that affects contributors).

- Add entries under `## Unreleased` at the top. Never edit a released section, except to fix broken links.
- Use these subsections, in this order, and only the ones needed: `### Breaking`, `### Added`, `### Fixed`, `### Changed`, `### Documentation`.
- One short line per change, written for provider users: what changed and how it affects them. Don't restate the code change (no implementation details, examples or error texts); the PR has those. Wrap code, attributes and function names in backticks.
- Group internal refactors, tests and naming changes into a single `Internal cleanup: ...` line per PR.
- End each entry with the PR or issue reference, e.g. `([#3])`, and define the link once below the section it first appears in:
  ```markdown
  [#3]: https://github.com/anapsix/terraform-provider-semvers/pull/3
  ```
- Dependency bumps use `old -> new`, e.g. `Upgrade terraform-plugin-log v0.10.0 -> v0.11.0`.

## Pull requests

- Branch from `main` (`fix/...`, `feat/...`, `chore/...`). Never commit directly to `main`; the only exception is the `version X.Y.Z` release commit (see [Releases](#releases)).
- PR description: a summary of user-visible changes, a **Breaking change** section when relevant, and a test plan checklist. Note when a PR depends on another one, unless the PRs are in a GitHub stack, which already shows the dependency.
- CI (`.github/workflows/test.yaml`) must pass: build, lint, generate drift check, and acceptance tests on Terraform 1.8, 1.12 and the latest minor.
- Pin GitHub Actions by full commit SHA with the version in a trailing comment (`uses: actions/checkout@<sha> # v7.0.1`). Dependabot (`.github/dependabot.yml`) opens a weekly grouped PR that bumps the SHAs and comments.
- Changing `.github/workflows/` over HTTPS requires a GitHub token with the `workflow` scope.

## Releases

Versions follow [SemVer](https://semver.org). While the provider is `0.x`, breaking changes bump the **minor** version (`0.8.1` -> `0.9.0`) and everything else bumps the patch version. Releases are cut from `main`, and only by the maintainer or at the maintainer's explicit request. Never tag or push a release on your own initiative.

1. Make sure `main` is up to date and CI is green.
2. In `CHANGELOG.md`, add `## X.Y.Z` directly below `## Unreleased`, so the pending entries move under the new version and `## Unreleased` is left empty. Fix any missing link definitions.
3. Commit on `main` with the message `version X.Y.Z` and push.
4. Create an annotated tag and push it:
   ```sh
   git tag -a vX.Y.Z -m vX.Y.Z
   git push origin vX.Y.Z
   ```
5. The `release` workflow (`.github/workflows/release.yaml`) runs GoReleaser: it builds all platforms, signs the checksums with the GPG key from repository secrets, and publishes the GitHub release. The Terraform Registry picks up the release from there.
6. GoReleaser's changelog is disabled, so once the release is published, set its notes from the CHANGELOG section:
   ```sh
   awk '/^## X\.Y\.Z$/{f=1;next} /^## /{f=0} f' CHANGELOG.md > /tmp/notes.md
   gh release edit vX.Y.Z --notes-file /tmp/notes.md
   ```
7. Check that the new version appears on the [registry page](https://registry.terraform.io/providers/anapsix/semvers).

To dry-run the release locally: `goreleaser release --snapshot --clean --skip=sign`.
