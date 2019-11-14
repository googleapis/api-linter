---
rule:
  aip: 158
  name: [core, '0158', response-next-page-token-field]
  summary: Paginated RPCs must have a `next_page_token` field in the response.
permalink: /158/response-next-page-token-field
redirect_from:
  - /0158/response-next-page-token-field
---

# Paginated methods: Next page token field

This rule enforces that all `List` and `Search` methods have a
`string next_page_token` field in the response message, as mandated in
[AIP-158][].

## Details

This rule looks at any message matching `List*Response` or `Search*Response`
and complains if either the `next_page_token` field is missing, or if it has
any type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page = 2;  // Field name should be `next_page_token`.
}
```

```proto
// Incorrect.
message ListBooksResponse {
  repeated Book books = 1;
  bytes next_page_token = 2;  // Field type should be `string`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `next_page_token` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0158::response-next-page-token-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message ListBooksResponse {
  repeated Book books = 1;
  string next_page = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-158]: https://aip.dev/158
[aip.dev/not-precedent]: https://aip.dev/not-precedent
