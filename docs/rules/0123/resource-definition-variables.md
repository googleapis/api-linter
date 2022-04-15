---
rule:
  aip: 123
  name: [core, '0123', resource-definition-variables]
  summary: Resource patterns should use consistent variable naming.
permalink: /123/resource-definition-variables
redirect_from:
  - /0123/resource-definition-variables
---

# Resource pattern variables

This rule enforces that resource patterns use consistent variable naming
conventions, as described in [AIP-123][].

## Details

This rule scans all files with `google.api.resource_definition` annotations, and
complains if variables in a `pattern` use camel case, or end in `_id`.

## Examples

**Incorrect** code for this rule:

```proto
import "google/api/resource.proto";

// Incorrect.
option (google.api.resource_definition) = {
  type: "library.googleapis.com/Book"
  // Should be: publishers/{publisher}/books/{book}
  pattern: "publishers/{publisher_id}/books/{book_id}"
};
```

```proto
import "google/api/resource.proto";

// Incorrect.
option (google.api.resource_definition) = {
  type: "library.googleapis.com/ElectronicBook"
  // Should be: publishers/{publisher}/electronicBooks/{electronic_book}
  pattern: "publishers/{publisher}/electronicBooks/{electronicBook}"
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

```proto
import "google/api/resource.proto";

// Correct.
option (google.api.resource_definition) = {
  type: "library.googleapis.com/ElectronicBook"
  pattern: "publishers/{publisher}/electronicBooks/{electronic_book}"
};
```

## Disabling

If you need to violate this rule, use a leading comment above the annotation.

```proto
import "google/api/resource.proto";

// (-- api-linter: core::0123::resource-definition-variables=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
option (google.api.resource_definition) = {
  type: "library.googleapis.com/Book"
  pattern: "publishers/{publisher_id}/books/{book_id}"
};
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-123]: http://aip.dev/123
[aip.dev/not-precedent]: https://aip.dev/not-precedent
