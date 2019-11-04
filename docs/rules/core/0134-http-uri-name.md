---
rule:
  aip: 134
  name: [core, '0134', http-uri-name]
  summary: Update methods must map the resource's name field to the URI.
---

# Update methods: HTTP URI name field

This rule enforces that all `Update` RPCs map the `name` field from the
resource object to the HTTP URI, as mandated in [AIP-134][].

## Details

This rule looks at any message matching beginning with `Update`, and complains
if the `name` variable from the resource (not the request message) is not
included in the URI. It _does_ check additional bindings if they are present.

Additionally, if the resource uses multiple words, it ensures that word
separation uses `snake_case`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UpdateBookRequest(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}"  // Should be `book.name`.
    body: "book"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UpdateBookRequest(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{book.name=publishers/*/books/*}"
    body: "book"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0134::http-uri-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc UpdateBookRequest(UpdateBookRequest) returns (Book) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/books/*}"
    body: "book"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip.dev/not-precedent]: https://aip.dev/not-precedent
