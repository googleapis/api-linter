---
rule:
  aip: 133
  name: [core, '0133', response-message-name]
  summary: Create methods must return the resource.
---

# Create methods: Resource response message

This rule enforces that all `Create` RPCs have a response message of the
resource, as mandated in [AIP-133][].

## Details

This rule looks at any message matching beginning with `Create`, and complains
if the name of the corresponding output message does not match the name of the
RPC with the prefix `Create` removed.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `Book`.
rpc CreateBook(CreateBookRequest) returns (CreateBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*}/books"
    body: "book"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*}/books"
    body: "book"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CreateBook(CreateBookRequest) returns (CreateBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*}/books"
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
