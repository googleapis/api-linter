# Rules

* [input-only](#input-only)
* [output-only](#output-only)
* [immutable](#immutable)

# Input only

When a field on a resource is input only, this should be described using the
`google.api.field_behavior` annotation instead of a comment. See
[AIP-203](https://goto.google.com/aip/203).

## Rule details

This rule inspects the leading comments of each field and if anything looks
similar to `Input only.`, it throws a warning. This rule actively skips this for
internal-only comments (e.g., `// (-- internal comment --)`).

### Examples

**Incorrect** code for this rule:

```proto {.bad}
message Book {
  // Secrets to be stored in the book.
  // @InputOnly
  string secret = 1;
}
```

```proto {.bad}
message Book {
  // Input only. Secret to be stored in the book.
  string secret = 1;
}
```

**Correct** code for this rule:

```proto {.good}
message Book {
  // Secret to be stored in the book.
  string secret = 1 [(google.api.field_behavior) = INPUT_ONLY];
}
```

Or if you must use a special annotation:

```proto
message Book {
  // Secret to be stored in the book.
  // (-- @InputOnly --)
  string generated_uri = 1 [(google.api.field_behavior) = INPUT_ONLY];
}
```

## Disabling

If you need to violate this rule for a single field, use an in-line or leading
comment. Please try to use internal comments instead though, as that should be
acceptable for our internal tooling that would consume these comments.
Additionally, please still include the annotation if possible.

```proto
message Book {
  // (-- api-linter: input-only-format=disabled --)
  // [Input only] Secret to be stored in the book.
  string secret = 1;
}
```

If you need to violate this rule for an entire file, use a file-level comment.

```proto
// (-- api-linter: input-only-format=disabled --)
syntax = "proto3";

message Book {
  // @InputOnly
  // Secret to be stored in the book.
  string secret = 1;
}
```

## Known limitations

-   None

## When to disable it

-   You have tooling that looks for a specific format (e.g., `@InputOnly`) and
    it **cannot** be an internal-only comment.

<!--*
# Document freshness: For more information, see go/fresh-source.
freshness: { owner: 'jjg' reviewed: '2019-08-20' }
*-->

# Output only

When a field on a resource is output only, this should be described using the
`google.api.field_behavior` annotation instead of a comment. See
[AIP-203](https://goto.google.com/aip/203).

## Rule details

This rule inspects the leading comments of each field and if anything looks
similar to `Output only.`, it throws a warning. This rule actively skips this
for internal-only comments (e.g., `// (-- internal comment --)`).

### Examples

**Incorrect** code for this rule:

```proto {.bad}
message Book {
  // A generated URI for this book.
  // @OutputOnly
  string generated_uri = 1;
}
```

```proto {.bad}
message Book {
  // Output only. A generated URI for this book.
  string generated_uri = 1;
}
```

**Correct** code for this rule:

```proto {.good}
message Book {
  // A generated URI for this book.
  string generated_uri = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

Or if you must use a special annotation:

```proto
message Book {
  // A generated URI for this book.
  // (-- @OutputOnly --)
  string generated_uri = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

## Disabling

If you need to violate this rule for a single field, use an in-line or leading
comment. Please try to use internal comments instead though, as that should be
acceptable for our internal tooling that would consume these comments.
Additionally, please still include the annotation if possible.

```proto
message Book {
  // (-- api-linter: output-only-format=disabled --)
  // [Output only] A generated URI for this book.
  string generated_uri = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

If you need to violate this rule for an entire file, use a file-level comment.

```proto
// (-- api-linter: output-only-format=disabled --)
syntax = "proto3";

message Book {
  // @OutputOnly
  // A generated URI for this book.
  string generated_uri = 1 [(google.api.field_behavior) = OUTPUT_ONLY];
}
```

## Known limitations

-   None

## When to disable it

-   You have tooling that looks for a specific format (e.g., `@OutputOnly`) and
    it **cannot** be an internal-only comment.

# Immutable

When a field on a resource is required, this should be described using the
`google.api.field_behavior` annotation instead of a comment. See
[AIP-203](https://goto.google.com/aip/203).

## Rule details

This rule inspects the leading comments of each field and if anything looks
similar to `Immutable.`, it throws a warning. This rule actively skips this for
internal-only comments (e.g., `// (-- internal comment --)`).

### Examples

**Incorrect** code for this rule:

```proto {.bad}
message Book {
  // The title of the book.
  // @Immutable
  string title = 1;
}
```

```proto {.bad}
message Book {
  // Immutable. The title of the book.
  string title = 1;
}
```

**Correct** code for this rule:

```proto {.good}
message Book {
  // The title of the book.
  string title = 1 [(google.api.field_behavior) = IMMUTABLE];
}
```

Or if you must use a special annotation:

```proto
message Book {
  // The title of the book.
  // (-- @Immutable --)
  string title = 1 [(google.api.field_behavior) = IMMUTABLE];
}
```

## Disabling

If you need to violate this rule for a single field, use an in-line or leading
comment. Please try to use internal comments instead though, as that should be
acceptable for our internal tooling that would consume these comments.
Additionally, please still include the annotation if possible.

```proto
message Book {
  // (-- api-linter: required-format=disabled --)
  // [Immutable] The title of the book.
  string title = 1 [(google.api.field_behavior) = IMMUTABLE];
}
```

If you need to violate this rule for an entire file, use a file-level comment.

```proto
// (-- api-linter: required-format=disabled --)
syntax = "proto3";

message Book {
  // @Immutable
  // The title of the book.
  string title = 1 [(google.api.field_behavior) = IMMUTABLE];
}
```

## Known limitations

-   None

## When to disable it

-   You have tooling that looks for a specific format (e.g., `@Immutable`) and
    it **cannot** be an internal-only comment.