---
rule:
  aip: 158
  name: [core, '0158', response-plural-first-field]
  summary: First field of Paginated RPCs' response should be plural.
permalink: /158/response-plural-first-field
redirect_from:
  - /0158/response-plural-first-field
---

# Paginated methods: Page token field

This rule enforces that all paginated response messages have a plural
name repeatable field as a first field, as mandated in [AIP-158][].

## Details

This rule looks at any response message  that has `next_page_token`
field and complains if the first field's name is not plural.


## Examples

**Incorrect** code for this rule:

```proto
// Incorrect
message ListBooksResponse {
  // Field name should be `books`.
  repeated Book book = 1;
  string next_page_token = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
}
```

**NOTE:** If the field is a resource, the plural form of the resource type is expected.

**Incorrect** code for this rule (field is a resource)

```proto
import "google/api/resource.proto";

message LibraryBook {
    option (google.api.resource) = {
        type: "example.com/LibraryBook"
        pattern: "libraryBooks/{libraryBook}"
        singular: "libraryBook"
        plural: "libraryBooks"
    };
    string name = 1;
}

// Incorrect.
message ListLibraryBooksResponse {
    repeated LibraryBook books = 1;
    string next_page_token = 2;
}
```

**Correct** code for this rule (resource type):

```proto
import "google/api/resource.proto";

message LibraryBook {
    option (google.api.resource) = {
        type: "example.com/LibraryBook"
        pattern: "libraryBooks/{libraryBook}"
        singular: "libraryBook"
        plural: "libraryBooks"
    };
    string name = 1;
}

// Correct.
message ListLibraryBooksResponse {
    repeated LibraryBook library_books = 1;
    string next_page_token = 2;
}
```


## Disabling

If you need to violate this rule, use a leading comment above the message or
above the field. Remember to also include an [aip.dev/not-precedent][] comment
explaining why.

```proto
// (-- api-linter: core::0158::response-plural-first-field=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
message ListBooksResponse {
  repeated Book book = 1;
  string next_page_token = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-158]: https://aip.dev/158
[aip.dev/not-precedent]: https://aip.dev/not-precedent
