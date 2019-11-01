---
rule:
  aip: 131
  name: [core, '0131', synonyms]
  summary: Get methods must be named starting with "Get".
---

# Get methods: Synonym check

This rule enforces that single-resource lookup methods have names starting with
`Get`, as mandated in [AIP-131][].

## Details

This rule looks at any message with names similar to `Get`, and suggests using
`Get` instead. It complains if it sees the following synonyms:

- Acquire
- Fetch
- Lookup
- Read
- Retrieve

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc FetchBook(FetchBookRequest) returns (Book) {  // Should be `GetBook`.
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
// (-- api-linter: core::0131::synonyms=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc FetchBook(GetBookReq) returns (Book) {
  option (google.api.http) = {
    get: "/v1/{name=publishers/*/books/*}"
  }
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-131]: https://aip.dev/131
[aip.dev/not-precedent]: https://aip.dev/not-precedent
