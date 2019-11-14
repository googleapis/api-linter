---
rule:
  aip: 133
  name: [core, '0133', http-method]
  summary: Create methods must use the POST HTTP verb.
permalink: /133/http-method
redirect_from:
  - /0133/http-method
---

# Create methods: POST HTTP verb

This rule enforces that all `Create` RPCs use the `POST` HTTP verb, as mandated
in [AIP-133][].

## Details

This rule looks at any message matching beginning with `Create`, and complains
if the HTTP verb is anything other than `POST`. It _does_ check additional
bindings if they are present.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    put: "/v1/{parent=publishers/*}/books"  // Should be `post:`.
    body: "book"
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
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    put: "/v1/{parent=publishers/*}/books"
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
