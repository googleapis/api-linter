---
rule:
  aip: 123
  name: [core, '0123', duplicate-resource]
  summary: Resource types should not be defined more than once.
permalink: /123/duplicate-resource
redirect_from:
  - /0123/duplicate-resource
---

# Resource annotation presence

This rule enforces that the same resource type doesn't appear in more than one
`google.api.resource` annotation, as described in [AIP-123][].

## Details

This rule complains about messages that have the same `type` for the
`google.api.resource` annotation, which frequently occur due to copy-paste
errors and messages spread across multiple files and/or packages. Duplicate
resource definitions can cause compilation problems in generated client code.

## Examples

**Incorrect** code for this rule:

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;
}

message Author {
  option (google.api.resource) = {
    // Incorrect: should be "library.googleapis.com/Author".
    type: "library.googleapis.com/Book"
    pattern: "authors/{author}"
  };

  string name = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;
}

message Author {
  option (google.api.resource) = {
    type: "library.googleapis.com/Author"
    pattern: "authors/{author}"
  };

  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a comment at the top of the file.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0123::duplicate-resource=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
syntax = "proto3";

message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;
}

message Author {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "authors/{author}"
  };

  string name = 1;
}
```

[aip-123]: http://aip.dev/123
[aip.dev/not-precedent]: https://aip.dev/not-precedent
