---
rule:
  aip: 134
  name: [core, '0134', method-signature]
  summary: |
    Update RPCs should annotate an appropriate method signature.
permalink: /134/method-signature
redirect_from:
  - /0134/method-signature
---

# Update methods: Method signature

This rule enforces that all `Update` standard methods have a
`google.api.method_signature` annotation with an appropriate value, as mandated
in [AIP-134][].

## Details

This rule looks at any method beginning with `Update`, and complains if the
`google.api.method_signature` annotation is missing, or if it is set to an
value other than `"{resource},update_mask`. Additional method signatures, if
present, are ignored.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  // A google.api.method_signature annotation should be present.
}
```

```proto
// Incorrect.
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  // Should be "book,update_mask".
  option (google.api.method_signature) = "book";
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.method_signature) = "book,update_mask";
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0134::method-signature=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc UpdateBook(UpdateBookRequest) returns (Book) {
  option (google.api.method_signature) = "book";
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip.dev/not-precedent]: https://aip.dev/not-precedent
