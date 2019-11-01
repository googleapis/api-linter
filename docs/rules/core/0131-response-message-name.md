---
rule:
  aip: 131
  name: [core, '0131', response-message-name]
  summary: Get methods must return the resource.
---

# Upper snake case values

This rule enforces that all `Get*` RPCs have a response message of the
resource, as mandated in [AIP-131][].

## Details

This rule looks at any message matching beginning with `Get`, and complains if
the name of the corresponding output message does not match the name of the RPC
with the prefix `Get` removed.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (GetBookResponse) {  // Should be `Book`.
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  }
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  }
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0131::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (GetBookResponse) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  }
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-131]: https://aip.dev/131
[aip.dev/not-precedent]: https://aip.dev/not-precedent
