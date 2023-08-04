---
rule:
  aip: 121
  name: [core, '0121', resource-must-support-get]
  summary: All resources must have a Standard Get method.
permalink: /121/resource-must-support-get
redirect_from:
  - /0121/resource-must-support-get
---

# Resource must support get

This rule enforces that all resources support the Get operation as mandated in
[AIP-121][].

## Details

This rule scans a service for Create, Update, and List methods for resources,
and ensures each one has a Get method.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
service Foo {
  // Book has a create, but no Get method.
  rpc CreateBook(CreateBookRequest) returns (Book) {};
}

message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
service Foo {
  rpc CreateBook(CreateBookRequest) returns (Book) {};
  rpc GetBook(GetBookRequest) returns (Book) {};
}

message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the service.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0121::resource-must-support-get=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
service Foo {
  // Book has a create, but no Get method.
  rpc CreateBook(CreateBookRequest) returns (Book) {};
}

message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-121]: https://aip.dev/121
[aip.dev/not-precedent]: https://aip.dev/not-precedent