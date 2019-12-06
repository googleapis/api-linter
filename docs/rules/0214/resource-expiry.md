---
rule:
  aip: 214
  name: [core, '0214', resource-expiry]
  summary: Resources with user-set expiry times should offer a ttl field.
permalink: /214/resource-expiry
redirect_from:
  - /0214/resource-expiry
---

# Resource expiry

This rule enforces that resources that have an expire time that can be set by
users also have an input-only `ttl` field, as recommended in [AIP-214][].

## Details

This rule looks at fields named `expire_time`, and complains if **none** of the
following conditions are met:

- A field named `ttl` also exists in the message.
- The `expire_time` field is annotated as `OUTPUT_ONLY`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  google.protobuf.Timestamp expire_time = 2;  // Should have `ttl` also.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  google.protobuf.Timestamp expire_time = 2
    [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

```proto
// Correct.
message Book {
  string name = 1;
  oneof expiration {
    google.protobuf.Timestamp expire_time = 2;
    google.protobuf.Duration ttl = 3;
  }
}
```

## Disabling

If you need to violate this rule, use a leading comment above the enum.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string name = 1;
  // (-- api-linter: core::0214::resource-expiry=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  google.protobuf.Timestamp expire_time = 2
    [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-214]: https://aip.dev/214
[aip.dev/not-precedent]: https://aip.dev/not-precedent
