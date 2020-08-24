---
rule:
  aip: 158
  name: [core, '0158', response-unary]
  summary: Paginated responses must not use streaming.
permalink: /158/response-unary
redirect_from:
  - /0158/response-unary
---

# Paginated methods: Unary responses

This rule enforces that all paginated methods (`List` and `Search` methods, or
methods with pagination fields) use unary responses, as mandated in
[AIP-158][].

## Details

This rule looks at any message matching `List*Response` or `Search*Response`,
or any response message that has `next_page_token` field, and complains if the
method uses gRPC server streaming (the `stream` keyword).

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Streaming is prohibited on paginated responses.
rpc ListBooks(ListBooksRequest) returns (stream ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc ListBooks(ListBooksRequest) returns (ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message or
above the field. Remember to also include an [aip.dev/not-precedent][] comment
explaining why.

```proto
// (-- api-linter: core::0158::response-unary
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBooks(ListBooksRequest) returns (stream ListBooksResponse) {
  option (google.api.http) = {
    get: "/v1/{parent=publishers/*}/books"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-158]: https://aip.dev/158
[aip.dev/not-precedent]: https://aip.dev/not-precedent
