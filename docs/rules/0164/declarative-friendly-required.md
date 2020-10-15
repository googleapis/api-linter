---
rule:
  aip: 164
  name: [core, '0164', declarative-friendly-required]
  summary: Declarative-friendly resources should have an Undelete method.
permalink: /164/declarative-friendly-required
redirect_from:
  - /0164/declarative-friendly-required
---

# Declarative-friendly resources: Undelete method required

This rule enforces that all declarative-friendly resources have an Undelete
method, as mandated in [AIP-164][].

## Details

This rule looks at any resource with a `google.api.resource` annotation that
includes `style: DECLARATIVE_FRIENDLY`, and complains if it does not have a
corresponding Undelete method.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should have an `UndeleteBook` method.
service Library {
}

message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;
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
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the resource.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
service Library {
}

// (-- api-linter: core::0164::declarative-friendly-required=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message Book {
  option (google.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "publishers/{publisher}/books/{book}"
    style: DECLARATIVE_FRIENDLY
  };

  string name = 1;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-164]: https://aip.dev/164
[aip.dev/not-precedent]: https://aip.dev/not-precedent
