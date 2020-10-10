---
rule:
  aip: 154
  name: [core, '0154', no-duplicate-etag]
  summary: |
    Etag fields should not be set on request messages that include
    the resource.
permalink: /154/no-duplicate-etag
redirect_from:
  - /0154/no-duplicate-etag
---

# Required etags

This rule enforces that `etag` fields are set on resources and requests that
reference those resources by name, but not on requests that include the
resource directly, as mandated in [AIP-154][].

## Details

This rule looks at any field named `etag` and complains if it is part of a
request message including a resource that itself includes an etag.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask = 2;
  string etag = 3;  // The Book message already includes etag.
}
```

**Correct** code for this rule:

```proto
// Correct.
message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask = 2;
  // (-- api-linter: core::0154::no-duplicate-etag=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string etag = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-154]: https://aip.dev/154
[aip.dev/not-precedent]: https://aip.dev/not-precedent
