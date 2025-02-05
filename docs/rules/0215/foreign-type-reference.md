---
rule:
  aip: 215
  name: [core, '0215', foreign-type-reference]
  summary: API-specific protos should be in versioned packages.
permalink: /215/versioned-packages
redirect_from:
  - /0215/versioned-packages
---

# Versioned packages

This rule enforces that none of the fields in an API reference message types in a different
proto package namespace other than well-known common packages.

## Details

This rule examines all fields in an API's messages and complains if the package of the
referenced message type is not the same or from one of the common packages such as
`google.api`, `google.protobuf`, etc.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
package foo.bar;
import "some/other.proto";

message SomeMessage {
    some.OtherMessage other_message = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
package foo.bar;
import "other.proto";

message SomeMessage {
    OtherMessage other_message = 1;
}
```

## Known issues

This package disallows subpackaging, and allows any package that ends with `.type` to
avoid flagging possible legitimate uses of AIP-213 component packages.

## Disabling

If you need to violate this rule, place the comment above the package statement.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0215::foreign-type-reference=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
package foo.bar;
```

[aip-215]: https://aip.dev/215
[aip.dev/not-precedent]: https://aip.dev/not-precedent
