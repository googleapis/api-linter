---
rule:
  aip: 158
  name: [core, '0158', request-page-token-field]
  summary: Paginated RPCs must have a `page_token` field in the request.
permalink: /158/request-page-token-field
redirect_from:
  - /0158/request-page-token-field
---

# Paginated methods: Page token field

This rule enforces that Pagination requests have a `string page_token` field in
the request message, as mandated in [AIP-158][].

## Details

This rule looks at any message mthat has `page_size` field and complains if
either the `page_token` field is missing, or if it has any type other than
`string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  string name = 1;
  int32 page_size = 2;
  string offset = 3;  // Field name should be `page_token`.
}
```

```proto
// Incorrect.
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  bytes page_token = 3;  // Field type should be `string`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `page_token` field is missing) or above the field (if it is the wrong
type). Remember to also include an [aip.dev/not-precedent][] comment explaining
why.

```proto
// (-- api-linter: core::0158::request-page-token-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message ListBooksRequest {
  string name = 1;
  int32 page_size = 2;
  string offset = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-158]: https://aip.dev/158
[aip.dev/not-precedent]: https://aip.dev/not-precedent
