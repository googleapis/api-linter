---
rule:
  aip: 191
  name: [core, '0191', filenames]
  summary: Proto files must use reasonable filenames.
permalink: /191/filenames
redirect_from:
  - /0191/filenames
---

# Filenames

This rule attempts to enforce reasonable filenames for proto files, as mandated
in [AIP-191][].

## Details

Because proto filenames show up in client libraries (for example, as import
paths), it is important not to have odd paths.

This rule currently complains if the filename:

- ...is set to the proto version.

## Examples

**Incorrect** filenames for this rule:

- `v1.proto`
- `v1beta1.proto`

**Correct** filenames for this rule:

- `library.proto`
- `library_service.proto`

## Disabling

If you need to violate this rule, use a comment at the top of the file.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0191::filenames=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
syntax = "proto3";
```

[aip-191]: https://aip.dev/191
[aip.dev/not-precedent]: https://aip.dev/not-precedent
