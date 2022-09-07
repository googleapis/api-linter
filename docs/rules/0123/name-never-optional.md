---
rule:
  aip: 123
  name: [core, '0123', name-never-optional]
  summary: Resource name fields must never be labeled with proto3_optional.
permalink: /123/name-never-optional
redirect_from:
  - /0123/name-never-optional
---

# Resource name field never optional

This rule enforces that the name field of a resource message is not labeled with
proto3_optional.

## Details

This rule scans for messages with a `google.api.resource` annotation and ensures
that the configured name field (either `name` or whichever field specified via
`name_field`) is not labeled as `optional`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  // The name field should not be labeled as optional.
  optional string name = 1;
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
```

## Disabling

If you need to violate this rule, use a leading comment above the message.

```proto
// (-- api-linter: core::0123::name-never-optional=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  optional string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip.dev/not-precedent]: https://aip.dev/not-precedent
