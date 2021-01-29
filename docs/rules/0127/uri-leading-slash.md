---
rule:
  aip: 127
  name: [core, '0127', uri-leading-slash]
  summary: URIs should always begin with a leading slash.
permalink: /127/uri-leading-slash
redirect_from:
  - /0127/uri-leading-slash
---

# URI Forward Slashes

This rule enforces that URIs must begin with a forward slash, as mandated in
[AIP-127][].

## Details

This rule scans all methods and complains if it finds a URI that does not start
with `/`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    // Should be /v1/{name=publishers/*/books/*}
    get: "v1/{name=publishers/*/books/*}"
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

Do not violate this rule. This would create an invalid URL.

[aip-127]: https://aip.dev/127
