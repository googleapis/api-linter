---
rule:
  aip: 135
  name: [core, '0135', force-field]
  summary: Delete RPCs for resources with child collections should have a `force` field in the request.
permalink: /135/force-field
redirect_from:
  - /0135/force-field
---

# Delete methods: `force` field

This rule enforces that the standard `Delete` method for a resource that parents
other resources in the service have a `bool force` field in the request message,
as mandated in [AIP-135][].

## Details

This rule looks at any message matching `Delete*Request` for a resource with
child resources in the same service and complains if the `force` field is
missing.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message DeletePublisherRequest {
  // Where Publisher parents the Book resource.
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Publisher"]; 

  // Missing `bool force` field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message DeletePublisherRequest {
  // Where Publisher parents the Book resource.
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Publisher"]; 

  // If set to true, any books from this publisher will also be deleted.
  // (Otherwise, the request will only work if the publisher has no books.)
  bool force = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message (if
the `name` field is missing) or above the field (if it is the wrong type).
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0135::force-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message DeletePublisherRequest {
  // Where Publisher parents the Book resource.
  string name = 1 [
    (google.api.resource_reference).type = "library.googleapis.com/Publisher"]; 
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-135]: https://aip.dev/135#cascading-delete
[aip.dev/not-precedent]: https://aip.dev/not-precedent
