---
rule:
  aip: 122
  name: [core, '0122', embedded-resource]
  summary: Resource references should not be embedded resources.
permalink: /122/embedded-resource
redirect_from:
  - /0122/embedded-resource
---

# Resource reference type

This rule enforces that references to resource are via
`google.api.resource_reference`, not by embedding the resource message, as
mandated in [AIP-122][].

## Details

This rule complains if it sees a resource field of type `message` that is also
annotated as a `google.api.resource`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  string name = 1;

  // Incorrect. Resource references should not be embedded resource messages.
  Author author = 2;
}

message Author {
    option (google.api.resource) = {
        type: "library.googleapis.com/Author"
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
    pattern: "books/{book}"
  };
  string name = 1;

  string author = 2 [(google.api.resource_reference) = {
    type: "library.googleapis.com/Author"
  }];
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

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
  string name = 1;

  // (-- api-linter: core::0122::embedded-resource=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  Author author = 2;
}

message Author {
    option (google.api.resource) = {
        type: "library.googleapis.com/Author"
        pattern: "authors/{author}"
    };

    string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-122]: https://aip.dev/122
[aip.dev/not-precedent]: https://aip.dev/not-precedent
