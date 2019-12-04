---
rule:
  aip: 215
  name: [core, '0215', versioned-packages]
  summary: API-specific protos should be in versioned packages.
permalink: /215/versioned-packages
redirect_from:
  - /0215/versioned-packages
---

# Versioned packages

This rule enforces that API-specific protos use versioned packages, as mandated
in [AIP-215][].

## Details

This rule examines the proto package and complains if it does not end something
that appears to be a version component, such as `v1` or `v1beta`. It also
permits proto packages ending in `type`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
package foo.bar;
```

```proto
// Incorrect.
package foo.bar.v1.resources;
```

**Correct** code for this rule:

```proto
// Correct.
package foo.bar.v1;
```

## Known issues

This rule blindly assumes that the proto it is looking at is for an individual
API (unless it ends in `type`). It will therefore complain about protos where
the lack of a version is proper (such as `google.longrunning`, for example).
This rule **should** be disabled in this situation.

## Disabling

If you need to violate this rule, use a leading comment above the enum.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0215::versioned-packages=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
package foo.bar;
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-215]: https://aip.dev/215
[aip.dev/not-precedent]: https://aip.dev/not-precedent
