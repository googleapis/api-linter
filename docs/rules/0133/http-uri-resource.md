---
rule:
  aip: 133
  name: [core, '0133', http-uri-resource]
  summary: The collection where the resource is added should map to the URI path.
permalink: /133/http-uri-resource
redirect_from:
  - /0133/http-uri-resource
---

# Create methods: HTTP URI resource

This rule enforces that the collection identifier used in the URI path is
provided in the definition for the resource being created, as mandated in
[AIP-133][].

## Details

This rule looks at any method beginning with `Create`, and retrieves the URI
path from the `google.api.http` annotation on the method. The final segment of
the URI is extracted as the `collection_identifier`.

This rule then ensures that each `google.api.http` annotation on the method's
return type contains the string `"{collection_identifier}/"` in each `pattern`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    // There collection identifier should appear after the final `/` in the URI.
    post: "/v1/"
    body: "book"
  };
}

message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };
}
```

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
  };
}

message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    // The pattern does not contain the collection identifier `books`.
    pattern: "publishers/{publisher}"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{parent=publishers/*}/books"
    body: "book"
  };
}

message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::http-uri-resource=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/"
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
