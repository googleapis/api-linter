---
rule:
  aip: 191
  name: [core, '0191', proto-package]
  summary: Proto package must match the directory structure.
permalink: /191/proto-package
redirect_from:
  - /0191/proto-package
---

# Protobuf Package

This rule attempts to enforce that the proto package and the directory structure
match, as mandated in [AIP-191][].

## Details

Accordig to the [Protobuf Style Guide][], the package name must correspond to
the directory structure.

This rule currently complains if the package and the directory structure do not
correspond.

## Examples

**Incorrect** directory structures and proto packages for this rule:

- `example/v1` `example.library.v1`
- `example/library/v1` `example.librarian.v1`

**Correct** directory structures and proto packages for this rule:

- `example/library/v1` `example.library.v1`
- `example/library/v1/types` `example.library.v1.types`

## Disabling

If you need to violate this rule, use a comment at the top of the file.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0191::proto-package=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
syntax = "proto3";
```

[aip-191]: https://aip.dev/191
[aip.dev/not-precedent]: https://aip.dev/not-precedent
[Protobuf Style Guide]: https://developers.google.com/protocol-buffers/docs/style#packages