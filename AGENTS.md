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
- **List arguments can contain `null`.** Read list elements as `types.String` and convert them with `versionStrings`, which turns `null` into an invalid version. Reading into `[]string` fails with an internal "Value Conversion Error".
- **Functions never see unknown values.** Terraform waits until every argument, including each list element, is known before calling a function, so function code does not handle unknowns. The data source and provider configuration do have to.

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
- Versions pinned in more than one place must be bumped together: golangci-lint (`.pre-commit-config.yaml` and the `lint` job), and the Terraform version of the `generate` job, which must match the newest version in the `test` matrix. Run it with the Go version from `go.mod` (e.g. `GOTOOLCHAIN=go1.26.10 golangci-lint run ./...`). Older golangci-lint builds fail with `export data version ... is greater than maximum supported version` under a newer local Go.
- Manual testing: `go install`, add a `dev_overrides` block for `anapsix-dev/semvers` pointing at `~/go/bin` (see the comment at the top of `test/main.tf`), then `cd test && terraform plan`. Check with both `terraform` and `tofu` when changing function or provider-config behavior.
- Dependabot locally: to see which GitHub Action updates Dependabot would propose, without waiting for the weekly run, use the [Dependabot CLI](https://github.com/dependabot/cli) (`brew install dependabot`; it is not a `gh` extension). It needs a running Docker daemon.
  ```sh
  LOCAL_GITHUB_ACCESS_TOKEN=$(gh auth token) \
    dependabot update github_actions anapsix/terraform-provider-semvers --local .
  ```
  It prints the proposed changes in its log and does not modify files or open PRs, so apply the new SHAs and version comments yourself. It ignores the `groups` setting and proposes one update per action. With Docker Desktop, if `docker` is not on `PATH`, set `DOCKER_HOST=unix://$HOME/.docker/run/docker.sock`.

## Code conventions

- **License headers:** every Go file needs one; see [Licensing](#licensing).
- Go identifiers are camelCase; Terraform-facing names (attributes, parameters) are snake_case.
- Parse with `semverutil.ParseVersion` / `semverutil.ParseConstraint`, not `semver.New*` directly, so errors carry the offending value.
- Function errors use `function.NewArgumentFuncError(position, err.Error())`, with each argument parsed separately so the position is known. Don't add generic "Error performing operation" wrappers.
- Leave framework defaults unset (e.g. no `AllowNullValue: false`).
- Exported identifiers get doc comments.

### Adding or changing a provider function

1. Add `internal/provider/function_<name>.go` and register it in `semversProvider.Functions`.
2. Add `internal/provider/function_<name>_test.go` covering known values, invalid input, null/wrong-type arguments, `null` list elements, and `ignore_invalid` passed as `true` and `false` where the function accepts it.
3. Add `examples/functions/<name>/function.tf`.
4. Update the function list in `README.md` and the provider `Description` in `provider.go`.
5. Run `go generate ./...` and commit `docs/`.

### Tests

- Acceptance tests use `ProtoV6ProviderFactories: testAccProtoV6ProviderFactories` and configure the provider with **real `provider "semvers" {}` blocks** in the test config. Never inject config into the provider struct; that is how the functions bug in #3 went unnoticed.
- Function tests start with `tfversion.SkipBelow(tfversion.Version1_8_0)`.
- Terraform wraps long error lines, so `ExpectError` regexes should use `\s+` between words, e.g. `` `"version_b"\s+parameter:\s+invalid\s+version` ``.
- For a bug fix, write the failing test first.

## Licensing

The provider is licensed under the [Mozilla Public License 2.0](LICENSE). Large parts of it started as HashiCorp's [provider scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding-framework), which is also MPL-2.0.

- **Every Go file starts with a license header.** Files written from scratch use:
  ```go
  // Copyright (c) Anastas Dancha
  // SPDX-License-Identifier: MPL-2.0
  ```
  Files that started from the scaffolding, or from another file that did, also keep HashiCorp's line first:
  ```go
  // Copyright (c) HashiCorp, Inc.
  // Copyright (c) Anastas Dancha
  // SPDX-License-Identifier: MPL-2.0
  ```
- **Never remove or alter an existing copyright notice**, however much the file has changed. MPL-2.0 section 3.4 forbids it. When splitting, renaming or copying a file, carry its notices over to the new file.
- **When unsure whether a file is derived, keep the notice.** Check its history with `git log --follow` and its first version's header before deciding a file is original.
- **Copying code from elsewhere:** only from sources under MIT, BSD, Apache-2.0 or MPL-2.0. Keep the original copyright and license notice next to the copied code, and name the source in the PR. Never copy from GPL, AGPL, unlicensed or unknown-license code; ask the maintainer instead.
- **New dependencies** must use one of the licenses above. Name the dependency and its license in the PR. `go run github.com/google/go-licenses@latest report ./...` lists the licenses of everything the provider links.
- **Don't edit `LICENSE`** or change the `SPDX-License-Identifier`. Relicensing is the maintainer's decision.
- Markdown, HCL examples, generated `docs/` and config files need no license header.

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
- Note when a PR depends on another one, unless the PRs are in a GitHub stack, which already shows the dependency.
- CI (`.github/workflows/test.yaml`) must pass: build, lint, generate drift check, and acceptance tests on Terraform 1.8, 1.12 and the latest minor.
- PRs are squash-merged (the only merge method enabled): the PR title becomes the commit subject on `main` and the PR description its body. Write the title in the commit subject format (`[type] short summary`); individual branch commit messages matter less.
- Keep history linear: update a branch by rebasing onto its base, not by merging.
- Force-push rewritten branches only with a lease pinned to the commit you last saw on the remote: `git push --force-with-lease=<branch>:<old-sha> origin <branch>`. A bare `--force-with-lease` can still overwrite commits that a fetch updated without you noticing; never use plain `--force`.
- When a branch lower in a stack is rewritten, rebase every PR above it onto it and force-push each the same way. Before pushing, check the result with `git diff <old-tip> HEAD`: it should show only the changes you meant to make.
- When pushing follow-up commits, update the PR description and the changelog in the same push. Tick the "CI green" checkbox only after the checks pass. In a stack of PRs, each PR's changelog entries belong on that PR's own branch.
- Workflow security: trigger on `pull_request`, never `pull_request_target`; grant least-privilege `permissions`; never interpolate event data (PR titles, branch names, bodies) into `run:` steps.
- Pin every GitHub Action, in every workflow, by full commit SHA with the version in a trailing comment (`uses: actions/checkout@<sha> # v7.0.1`). Dependabot (`.github/dependabot.yml`) opens a weekly grouped PR that bumps the SHAs and comments.
- Changing `.github/workflows/` over HTTPS requires a GitHub token with the `workflow` scope.

### PR descriptions

- **Start with a bullet list of the high-level changes.** No heading or intro sentence before it, and no implementation details; the diff has those.
- Keep each bullet to one short, plain sentence about the outcome ("`sort` and `pick` handle `null` list elements"), not the mechanism ("read list elements as `types.String`").
- Put user-visible changes first and the most important first. Fold internal refactors, tests and naming into one `Internal cleanup: ...` bullet.
- Aim for at most about 7 bullets; merge related changes rather than listing each one.
- After the bullets, add short sections only when they apply:
  - **Breaking change:** what users must change, in one or two sentences.
  - **Notes:** anything reviewers can't see in CI, e.g. changes that only run at release time.
  - **Test plan:** a checklist of what was run.
- Describe the PR as it is now, not how it got there. When pushing follow-up commits, rewrite the bullets instead of appending "review follow-ups" sections.
- Avoid file paths, function names and code blocks unless they are what users touch (attributes, function arguments). One short usage example is fine for a new feature or a breaking change.

### Commit messages

- Subject: `[type] short summary`, lowercase, imperative, no trailing period, at most about 72 characters. Types used here: `[feat]`, `[fix]`, `[chore]`, `[docs]`, `[ci]`. Release commits stay `version X.Y.Z`.
- Body (optional): a few short bullets of the high-level changes, plus the why when it isn't obvious. Don't list every file or restate the diff.
- One logical change per commit; split unrelated changes rather than describing them all in one subject.

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

CI does not exercise `release.yaml` or `.goreleaser.yaml`; they only run when a tag is pushed. After changing either, run `goreleaser check` and the snapshot dry-run, and note in the PR that the next release is the first real test.
