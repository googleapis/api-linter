---
rule:
  aip: 121
  name: [core, '0121', no-cycles]
  summary: Resources must not form a resource reference cycle.
permalink: /121/no-cycles
redirect_from:
  - /0121/no-cycles
---

# Resource must support get

This rule enforces that resources do not create reference cycles as mandated in
[AIP-121][].

## Details

This rule scans the fields of every resource and ensures that any references to
other resources do not create a cycle between them.

## Examples

**Incorrect** code for this rule:

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1;

  // Incorrect. Creates potential reference cycle.
  string author = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Author"
  ];
}

message Author {
  option (google.api.resource) = {
    type: "library.googleapis.com/Author"
    pattern: "authors/{author}"
  };

  string name = 1;

  // Incorrect. Creates potential reference cycle.
  string book = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

**Correct** code for this rule:

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1;

  // Correct because the other reference is OUTPUT_ONLY.
  string author = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Author"
  ];
}

message Author {
  option (google.api.resource) = {
    type: "library.googleapis.com/Author"
    pattern: "authors/{author}"
  };

  string name = 1;

  // Correct because an OUTPUT_ONLY reference breaks the mutation cycle.
  string book = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Book",
    (google.api.field_behavior) = OUTPUT_ONLY
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the service.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1;

  // (-- api-linter: core::0121::no-cycles=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string author = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Author"
  ];
}

message Author {
  option (google.api.resource) = {
    type: "library.googleapis.com/Author"
    pattern: "authors/{author}"
  };

  string name = 1;

  // (-- api-linter: core::0121::no-cycles=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string book = 2 [
    (google.api.resource_reference).type = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-121]: https://aip.dev/121
[aip.dev/not-precedent]: https://aip.dev/not-precedent