---
rule:
  aip: 136
  name: [core, '0136', http-name-variable]
  summary:
    Custom methods should only use `name` if the RPC noun matches the resource.
permalink: /136/http-name-variable
redirect_from:
  - /0136/http-name-variable
---

# Custom methods: HTTP name variable

This rule enforces that custom methods only use the `name` variable if the RPC
noun matches the resource, as mandated in [AIP-136][].

## Details

This rule looks at custom methods and, if the URI contains a `name` variable,
it ensures that the RPC name ends with the same text as the final component in
the URI (after adjusting for case).

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// The variable should be "book", or the RPC name should change.
rpc WritePage(WritePageRequest) return (WritePageResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:writePage"
    body: "*"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc WritePage(WritePageRequest) return (WritePageResponse) {
  option (google.api.http) = {
    post: "/v1/{book=publishers/*/books/*}:writePage"
    body: "*"
  };
}
```

```proto
// Correct.
rpc WriteBook(WriteBookRequest) return (WriteBookResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:write"
    body: "*"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0136::http-name-variable=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc WritePage(WritePageRequest) return (WritePageResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}:writePage"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-136]: https://aip.dev/136
[aip.dev/not-precedent]: https://aip.dev/not-precedent
