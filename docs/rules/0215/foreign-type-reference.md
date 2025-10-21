---
rule:
  aip: 215
  name: [core, '0215', foreign-type-reference]
  summary: API should not reference foreign types outside of the API scope.
permalink: /215/foreign-type-reference
redirect_from:
  - /0215/foreign-type-reference
---

# Foreign type reference

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

This check only allows subpackage usage within a versioned path, but generates warnings for unversioned subpackage usage.
It also ignores if the referenced package is a "common" package like `google.api`, or if the package path ends in `.type` indicating
the package is an AIP-213 component package.

Examples of foreign type references and their expected results:

| Calling Package | Referenced Package | Result       |
| --------------- | ------------------ | ------------ |
| foo.bar         | foo.xyz            | lint warning |
| foo.v2.bar      | foo.v2.xyz         | ok           |
| foo.bar         | foo.type           | ok           |
| foo.bar         | google.api         | ok           |

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
