---
rule:
  aip: 132
  name: [core, '0132', request-show-deleted-required]
  summary: List requests must have a `show-deleted` field for resources
  supporting soft delete.
permalink: /132/request-show-deleted-required
redirect_from:
  - /0132/request-show-deleted-required
---

# List methods: `show_deleted` field

This rule enforces that all `List` standard methods have a `bool show_deleted`
field in the request message if the resource supports soft delete, as mandated
in [AIP-132](http://aip.dev/132).

## Details

This rule looks at any message matching `List*Request` and complains if the
`show_deleted` field is missing and the corresponding resource has an
`Undelete` method.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.

service Library {
  ...
  rpc UndeleteBook(UndeleteBookRequest) returns (Book) { ... }
}

// Missing the `bool show_deleted` field.
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
}
```

**Correct** code for this rule:

```proto
// Correct.

service Library {
  ...
  rpc UndeleteBook(UndeleteBookRequest) returns (Book) { ... }
}

message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
  bool show_deleted = 4;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0132::request-show-deleted-required=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message ListBooksRequest {
  string parent = 1;
  int32 page_size = 2;
  string page_token = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip.dev/not-precedent]: https://aip.dev/not-precedent
