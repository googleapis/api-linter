---
rule:
  aip: 132
  name: [core, '0132', method-signature]
  summary: |
    List RPCs should annotate a method signature of "parent".
permalink: /132/method-signature
redirect_from:
  - /0132/method-signature
---

# List methods: Method signature

This rule enforces that all `List` standard methods have a
`google.api.method_signature` annotation with a value of `"parent"`, as
mandated in [AIP-132][].

## Details

This rule looks at any method beginning with `List`, and complains if the
`google.api.method_signature` annotation is missing, or if it is set to any
value other than `"parent"`. Additional method signatures, if present, are
ignored.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc ListBooks(ListBooksRequest) returns (Book) {
  // A google.api.method_signature annotation should be present.
}
```

```proto
// Incorrect.
rpc ListBooks(ListBooksRequest) returns (Book) {
  option (google.api.method_signature) = "publisher";  // Should be "parent".
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc ListBooks(ListBooksRequest) returns (Book) {
  option (google.api.method_signature) = "parent";
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::method-signature=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc ListBooks(ListBooksRequest) returns (Book);
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-132]: https://aip.dev/132
[aip.dev/not-precedent]: https://aip.dev/not-precedent
