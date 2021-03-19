---
rule:
  aip: 164
  name: [core, '0164', resource-expire-time-field]
  summary: Resources supporting soft delete must have an `expire_time` field.
permalink: /164/resource-expire-time-field
redirect_from:
  - /0164/resource-expire-time-field
---

# Resources supporting soft delete: `expire_time` field required

This rule enforces that all resources supporting soft delete have an
`google.protobuf.Timestamp expire_time` field, as mandated in [AIP-164][].

## Details

This rule looks at any resource with a corresponding `Undelete*` method, and
complains if it does not have a `google.protobuf.Timestamp expire_time` field.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
service Library {
  rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
    option (google.api.http) = {
      post: "/v1/{name=publishers/*/books/*}:undelete"
      body: "*"
    };
  }
}

message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;

  // Should have an `expire_time` field.
}
```

**Correct** code for this rule:

```proto
// Correct.
service Library {
  rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
    option (google.api.http) = {
      post: "/v1/{name=publishers/*/books/*}:undelete"
      body: "*"
    };
  };
}

message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;

  google.protobuf.Timestamp expire_time = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the resource.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
service Library {
  rpc UndeleteBook(UndeleteBookRequest) returns (Book) {
    option (google.api.http) = {
      post: "/v1/{name=publishers/*/books/*}:undelete"
      body: "*"
    };
  };
}

// (-- api-linter: core::0164::resource-expire-time-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
  };

  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent
