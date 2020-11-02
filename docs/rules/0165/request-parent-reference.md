---
rule:
  aip: 165
  name: [core, '0165', request-parent-reference]
  summary: |
    Purge requests should annotate the `parent` field with `google.api.resource_reference`.
permalink: /165/request-parent-reference
redirect_from:
  - /0165/request-parent-reference
---

# Purge requests: Parent resource reference

This rule enforces that all `Purge` requests have
`google.api.resource_reference` on their `string parent` field, as mandated in
[AIP-165][].

## Details

This rule looks at the `parent` field of any message matching `Purge*Request`
and complains if it does not have a `google.api.resource_reference` annotation.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message PurgeBooksRequest {
  // The `google.api.resource_reference` annotation should also be included.
  string parent = 1 [(google.api.field_behavior) = REQUIRED];

  string filter = 2 [(google.api.field_behavior) = REQUIRED];

  bool force = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.
message PurgeBooksRequest {
  string parent = 1 [
    (google.api.field_behavior) = REQUIRED,
    (google.api.resource_reference).child_type = "library.googleapis.com/Book"
  ];

  string filter = 2 [(google.api.field_behavior) = REQUIRED];

  bool force = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
message PurgeBooksRequest {
  // (-- api-linter: core::0165::request-parent-reference=disabled
  //     aip.dev/not-precedent: We need to do this because reasons. --)
  string parent = 1 [(google.api.field_behavior) = REQUIRED];

  string filter = 2 [(google.api.field_behavior) = REQUIRED];

  bool force = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-165]: https://aip.dev/165
[aip.dev/not-precedent]: https://aip.dev/not-precedent
