---
rule:
  aip: 123
  name: [core, '0123', resource-definition-type-name]
  summary: Resource type names must be of the form {Service Name}/{Type}.
permalink: /123/resource-definition-type-name
redirect_from:
  - /0123/resource-definition-type-name
---

# Resource type name

This rule enforces that files that define a resource with the
`google.api.resource_definition` annotation have a properly formatted `type`, as
described in [AIP-123][].

## Details

This rule scans files with `google.api.resource_definition` annotations, and
validates the format of the `type` field conforms to `{Service Name}/{Type}`.

## Examples

**Incorrect** code for this rule:

```proto
import "google/api/resource.proto";

// Incorrect.
option (google.api.resource_definition) = {
  // Should not have more than one separating '/'.
  type: "library.googleapis.com/Genre/Mystery/Book"
  pattern: "publishers/{publisher}/books/{book}"
};
```

**Correct** code for this rule:

```proto
import "google/api/resource.proto";

// Correct.
option (google.api.resource_definition) = {
  type: "library.googleapis.com/Book"
  pattern: "publishers/{publisher}/books/{book}"
};
```

## Disabling

If you need to violate this rule, use a leading comment above the annotation.

```proto
import "google/api/resource.proto";

// (-- api-linter: core::0123::resource-definition-type-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
option (google.api.resource_definition) = {
  type: "library.googleapis.com/Genre/Mystery/Book"
  pattern: "publishers/{publisher}/books/{book}"
};
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-123]: http://aip.dev/123
[aip.dev/not-precedent]: https://aip.dev/not-precedent
