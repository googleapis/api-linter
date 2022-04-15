---
rule:
  aip: 123
  name: [core, '0123', resource-definition-pattern]
  summary: Resource annotations should define a pattern.
permalink: /123/resource-definition-pattern
redirect_from:
  - /0123/resource-definition-pattern
---

# Resource patterns

This rule enforces that files that define a resource with the
`google.api.resource_definition` annotation have a `pattern` defined, as
described in [AIP-123][].

## Details

This rule scans all `google.api.resource_definition` annotations in all files,
and complains if `pattern` is not provided at least once. It also complains if
the segments outside of variable names contain underscores.

## Examples

**Incorrect** code for this rule:

```proto
import "google/api/resources.proto";

// Incorrect.
option (google.api.resource_definition) = {
  type: "library.googleapis.com/Book"
  // pattern should be here
};
```

```proto
import "google/api/resources.proto";

// Incorrect.
option (google.api.resource_definition) = {
  type: "library.googleapis.com/ElectronicBook"
  // Should be: publishers/{publisher}/electronicBooks/{electronic_book}
  pattern: "publishers/{publisher}/electronic_books/{electronic_book}"
};
```

**Correct** code for this rule:

```proto
import "google/api/resources.proto";

// Correct.
option (google.api.resource_definition) = {
  type: "library.googleapis.com/Book"
  pattern: "publishers/{publisher}/books/{book}"
};
```

```proto
import "google/api/resource.proto";

// Correct.
option (google.api.resource_definition) = {
  type: "library.googleapis.com/ElectronicBook"
  pattern: "publishers/{publisher}/electronicBooks/{electronic_book}"
};
```

## Disabling

If you need to violate this rule, use a comment on the annotation.

```proto
import "google/api/resource.proto";

// (-- api-linter: core::0123::resource-definition-pattern=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
option (google.api.resource_definition) = {
  type: "library.googleapis.com/Book"
};
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-123]: http://aip.dev/123
[aip.dev/not-precedent]: https://aip.dev/not-precedent
