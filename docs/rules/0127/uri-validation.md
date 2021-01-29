---
rule:
  aip: 127
  name: [core, '0127', uri-validation]
  summary: HTTP URIs must have the proper format.
permalink: /127/uri-validation
redirect_from:
  - /0127/uri-validation
---

# HTTP URI format

This rule enforces that the URI in an HTTP annotation has the proper format.

## Details

This rule scans all methods that a `google.api.http` annotation is present on
and complains if the URI format is invalid.

## Examples

### Missing forward slash prefix

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetBook(GetBookRequest) returns (Book) {
  option (google.api.http) = {
    // Missing a forward slash at the beginning.
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

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0127::uri-validation=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
    get: "v1/{name=publishers/*/books/*}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-127]: https://aip.dev/127
[aip.dev/not-precedent]: https://aip.dev/not-precedent
