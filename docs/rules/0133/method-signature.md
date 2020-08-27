---
rule:
  aip: 133
  name: [core, '0133', method-signature]
  summary: |
    Create RPCs should annotate an appropriate method signature.
permalink: /133/method-signature
redirect_from:
  - /0133/method-signature
---

# Create methods: Method signature

This rule enforces that all `Create` standard methods have a
`google.api.method_signature` annotation with an appropriate value, as mandated
in [AIP-133][].

## Details

This rule looks at any method beginning with `Create`, and complains if the
`google.api.method_signature` annotation is missing, or if it is set to an
incorrect value. Additional method signatures, if present, are ignored.

The correct value is `"parent,{resource},{resource}_id"` if the `{resource}_id`
field exists, and `"parent,{resource}"` otherwise.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  // A google.api.method_signature annotation should be present.
}
```

```proto
// Incorrect.
rpc CreateBook(CreateBookRequest) returns (Book) {
  // Should be "parent,book" or "parent,book,book_id", depending on whether
  // a "book_id" field exists.
  option (google.api.method_signature) = "publisher,book";
}
```

**Correct** code for this rule:

If the `book_id` field does not exist:

```proto
// Correct.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.method_signature) = "parent,book";
}
```

If the `book_id` field exists:

```proto
// Correct.
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.method_signature) = "parent,book,book_id";
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::method-signature=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc CreateBook(CreateBookRequest) returns (Book) {
  option (google.api.method_signature) = "publisher,book";
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-133]: https://aip.dev/133
[aip.dev/not-precedent]: https://aip.dev/not-precedent
