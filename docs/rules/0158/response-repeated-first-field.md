---
rule:
  aip: 158
  name: [core, '0158', response-repeated-first-field]
  summary:
    First field (by both position and field number) of Paginated RPCs' response
    should be repeated.
permalink: /158/response-repeated-first-field
redirect_from:
  - /0158/response-repeated-first-field
---

# Paginated methods: Page token field

This rule enforces that all `List` and `Search` methods have a repeatable field
as a first field in the response message, as mandated in [AIP-158][].

## Details

This rule looks at any message matching `List*Response` or `Search*Response`
that has `next_page_token` field and complains if the first field (by both
position and field number) is not repeated.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect
message ListStudentProfilesResponse {
    Profile studentProfiles = 1; // Field should be repeated.
    string next_page_token = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListStudentProfilesResponse {
  repeated Profile studentProfiles = 1;
  string next_page_token = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message or
above the field. Remember to also include an [aip.dev/not-precedent][] comment
explaining why.

```proto
// (-- api-linter: core::0158::response-repeated-first-field
//     aip.dev/not-precedent: We need to do this because reasons. --)
message ListStudentProfilesResponse {
    Profile studentProfiles = 1; // Field should be repeated.
    string next_page_token = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-158]: https://aip.dev/158
[aip.dev/not-precedent]: https://aip.dev/not-precedent
