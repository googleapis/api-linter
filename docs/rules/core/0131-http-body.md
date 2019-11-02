---
rule:
  aip: 131
  name: [core, '0131', http-body]
  summary: Get methods must not have an HTTP body.
---

# Get methods: Request message

This rule enforces that all `Get` RPCs omit the HTTP `body`, as mandated in
[AIP-131][].

## Details

This rule looks at any message matching beginning with `Get`, and complains if
the HTTP `body` field is set.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
    body: "*"  // This should be absent.
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0131::http-body=disabled
//     api-linter: core::0131::http-method=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}"
    body: "*"
  };
}
```

**Important:** HTTP `GET` requests are unable to have an HTTP body, due to the
nature of the protocol. The only valid way to include a body is to also use a
different HTTP method (as depicted above).

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-131]: https://aip.dev/131
[aip.dev/not-precedent]: https://aip.dev/not-precedent
