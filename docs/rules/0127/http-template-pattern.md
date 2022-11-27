---
rule:
  aip: 127
  name: [core, '0127', http-template-pattern]
  summary: |
    HTTP annotations should match the patterns defined by their resources.
permalink: /127/http-template-pattern
redirect_from:
  - /127/http-template-pattern
---

# HTTP Pattern Variables

This rule enforces that any HTTP annotations that reference a resource must
match one of the pattern strings defined by that resource, as mandated in
[AIP-127][].

## Details

This rule checks methods for the `google.api.http` annotation. If one is present
and if it includes a variable, this rule will check that the variable's template
matches a pattern for the corresponding resource.

This rule first finds the field or sub-field in the method's request message
that matches the variable name. It then identifies the name of the resource by
reading that field's `google.api.resource_reference` annotation, and uses this
to find the resource message. Finally, it reads the list of patterns in the
resource message's `google.api.resource` annotation, and it confirms that the
variable's template matches one of those patterns.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// The template for the `name` variable in the `google.api.http` annotation
// is missing segments from the Book message's `pattern`.
rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
        get: "v1/{name=shelves/*}"
    };
}
message GetBookRequest {
    string name = 1 [
        (google.api.resource_reference).type = "library.googleapis.com/Book"
    ];
}
message Book {
    option (google.api.resource) = {
        type: "library.googleapis.com/Book"
        pattern: "shelves/{shelf}/books/{book}"
    };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
        get: "v1/{name=shelves/*/books/*}"
    };
}
message GetBookRequest {
    string name = 1 [
        (google.api.resource_reference).type = "library.googleapis.com/Book"
    ];
}
message Book {
    option (google.api.resource) = {
        type: "library.googleapis.com/Book"
        pattern: "shelves/{shelf}/books/{book}"
    };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0127::http-template-pattern=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc GetBook(GetBookRequest) returns (Book) {
    option (google.api.http) = {
        get: "v1/{name=shelves/*}"
    };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-127]: https://aip.dev/127
[aip.dev/not-precedent]: https://aip.dev/not-precedent