# Zenjob API Linter (Fork of Google API Linter)

## About this Fork

This repository is a Zenjob-maintained fork of the [Google API Linter](https://github.com/googleapis/api-linter). Its purpose is to:

- Add and maintain Zenjob-specific lint rules and organizational customizations.
- Provide a Docker image and distribution pipeline under the Zenjob namespace.
- Stay closely in sync with upstream improvements and bugfixes.

> **Note:** We are tracking the upstream discussion for a plugin mechanism in [googleapis/api-linter#1485](https://github.com/googleapis/api-linter/issues/1485). If/when official plugin support is implemented, this fork may be retired in favor of upstream extensibility.

## Keeping the Fork Up to Date

To update this fork with the latest changes from upstream:

1. Add the upstream remote if not already present:

   ```sh
   git remote add upstream https://github.com/googleapis/api-linter.git
   ```

2. Fetch upstream changes:

   ```sh
   git fetch upstream
   ```

3. Rebase or merge upstream/main into your working branch:

   ```sh
   git checkout main
   git merge upstream/main
   # or: git rebase upstream/main
   ```

4. Resolve any conflicts, test, and push changes to the Zenjob fork.

---

# Google API Linter

[![ci](https://github.com/googleapis/api-linter/actions/workflows/ci.yaml/badge.svg)](https://github.com/googleapis/api-linter/actions/workflows/ci.yaml)
![latest release](https://img.shields.io/github/v/release/googleapis/api-linter)
![go version](https://img.shields.io/github/go-mod/go-version/googleapis/api-linter)

The API linter provides real-time checks for compliance with many of Google's
API standards, documented using [API Improvement Proposals][]. It operates on
API surfaces defined in [protocol buffers][].

It identifies common mistakes and inconsistencies in API surfaces:

```proto
// Incorrect.
message GetBookRequest {
  // This is wrong; it should be spelled `name`.
  string book = 1;
}
```

When able, it also offers a suggestion for the correct fix.

[_Read more ≫_](https://linter.aip.dev/)

## Versioning

The Google API linter does **not** follow semantic versioning. Semantic
versioning is challenging for a tool like a linter because the addition or
correction of virtually any rule is "breaking" (in the sense that a file that
previously reported no problems may now do so).

Therefore, the version numbers refer to the linter's core interface. In
general:

- Releases with only documentation, chores, dependency upgrades, and/or
  bugfixes are patch releases.
- Releases with new rules (or potentially removed rules) are minor releases.
- Releases with core interface alterations are major releases. This could
  include changes to the internal Go interface or the CLI user interface.

**Note:** Releases that increment the Go version will be considered minor.

This is an attempt to follow the spirit of semantic versioning while still
being useful.

## Contributing

If you are interested in contributing to the API linter, please review the [contributing guide](https://linter.aip.dev/contributing) to learn more.

## License

This software is made available under the [Apache 2.0][] license.

[apache 2.0]: https://www.apache.org/licenses/LICENSE-2.0
[api improvement proposals]: https://aip.dev/
[protocol buffers]: https://developers.google.com/protocol-buffers
[rule documentation]: ./rules/index.md
