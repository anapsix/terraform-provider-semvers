## 0.4.0

### Added:
- The  `provider::semvers::pick` function which takes list of semver strings,
  and semver constraint, and returns a list of filtered semver strings, sorted and deduped,
  matching the constraint. See [Masterminds/semver](https://github.com/Masterminds/semver/tree/master?tab=readme-ov-file#checking-version-constraints) for constraint syntax.

## 0.3.1

### Fixed

- The changelog formatting
- The provider and function docs
- The function examples

## 0.3.0

### Added:
- The  `provider::semvers::sort` function which takes list of semver strings,
  and returns a list of semver strings, sorted and deduped

## 0.2.1

### Added:
- Deduplication based on parsed semver version string

## 0.2.0

### Added:
- The `original` attribute to `sorted_versions` object with passed semver string
- The descriptions to data-source attributes

## 0.1.0

- Initial release
