---
rule:
  aip: 134
  name: [core, '0134', request-mask-field]
  summary: Update RPCs must have a field mask in the request.
permalink: /134/request-mask-field
redirect_from:
  - /0134/request-mask-field
---

# Update methods: Mask field

This rule enforces that all `Update` standard methods have a field in the
request message for the field mask, as mandated in [AIP-134][].

## Details

This rule looks at any message matching `Update*Request` and complains if it
can not find a field named `update_mask`, or if that field has any type other
than `google.protobuf.FieldMask`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// `google.protobuf.FieldMask update_mask` is missing.
message UpdateBookRequest {
  Book book = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
message UpdateBookRequest {
  Book book = 1;
  google.protobuf.FieldMask update_mask = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the resource field is missing) or above the field (if it is improperly named).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0134::request-mask-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message UpdateBookRequest {
  Book book = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-134]: https://aip.dev/134
[aip.dev/not-precedent]: https://aip.dev/not-precedent
