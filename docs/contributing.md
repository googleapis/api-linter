---
---

# Contributing

We are thrilled that you are interested in contributing to the API linter. This
software is fully open-sourced, licensed under the Apache license, and we do
accept contributions.

The most common way to contribute is by writing a new linter rule.

## Setup

The API linter is written in [Go][], so you will need the Go language
installed. The version of Go you need is officially documented in our
[`go.mod`][] file. Most of the time, we will most likely support the most
recent two versions.

Once you have Go installed, you can clone the repository the usual way, and
then follow up by running the tests:

```bash
$ git clone https://github.com/googleapis/api-linter
$ cd api-linter
$ go test ./...
```

**Note:** Unless you have commit bit, you will likely need to make your own
fork in GitHub in order to send us pull requests (in which case you clone your
fork instead).

## Writing rules

One of the best ways to contribute is by writing a new lint rule. Rules are
located in the `rules/` directory. Rules are grouped into packages based on the
[AIP][] that mandates the behavior.

**Important:** **All** linter rules **must** have a corresponding AIP that
mandates the behavior. There are no exceptions to this.

Additionally, we observe the following guidelines around rules:

- One rule per file.
  - The filename **must** correspond to the last segment of the rule name (with
    hyphens converted to underscores).
  - The rule must have a corresponding test (or tests). Each rule package
    **must** maintain 100% coverage of statements.
- All rules **must** have a three part name. The first part is the "rule group"
  (such as `core`), the second part is the AIP number, zero-padded, and the
  third part is a unique name for the rule. The name only has to be unique
  within the scope of the AIP.
  - If word separation is needed, `kebab-case` **must** be used.
- Every rule **must** have corresponding documentation. Documentation lives in
  the `docs/` directory, and powers this documentation site using GitHub Pages.

We have a CI lint to remind you to do so.

Writing a rule is straightforward: the linter employs a [visitor pattern][]
that goes to every descriptor in the proto file and runs each lint rule against
it. Most rules are run against only a certain _type_ of descriptor (for
example, a "message rule" is run against all of the message descriptors, but
not services or fields).

Consider a bare-bones message rule:

```go
var myRule = &lint.MessageRule{
  Name: lint.NewRuleName(0, "my-rule"),
  LintMessage: func(m *desc.MessageDescriptor) []lint.Problem {
    // This lint rule does nothing and always passes.
    return nil
  },
}
```

The actual lint function takes a [protoreflect][] descriptor. Beyond this, the
function is free-form; the developer can check anything desired and return a
slice of [`Problem`][] objects.

## Registering rules

Once a rule is written, it must be _registered_ with the rule registry, which
is defined in [`rules.go`][].

There are two steps:

1. In the registry in `rules.go`, ensure that the corresponding AIP package is
   imported, and that its `AddRules` function is called.
2. In the AIP package, ensure that the new rule is included in the `AddRules`
   function.

We have a CI lint to remind you to do so.

## Documentation

Rule documentation is the primary purpose of this site, and it is important
that all rules are documented. This documentation is written in Markdown, and
goes in the `docs/rules/` directory. The naming convention is
`{aip}-{rule_name}.md`:

- `{aip}` is the _four-digit_ AIP number (zero-padded if needed!)
- `{rule_name}` is the final component of the rule name in the rule itself.

The actual Markdown document is fairly boilerplate, and copy and paste from
another file is reasonable.

The top of the file **must** include the proper "front matter" for GitHub
Pages. The format is:

    ---
    rule:
      aip: 0
      name: [core, '0000', my-rule]
    ---

- The AIP number **must** be included as an integer, and **must** be set as a
  string in the `name` array (quotes are required to keep the YAML parser from
  interpreting it as an integer and dropping leading zeroes).
- The `name` field **must** be the name of the rule (as passed to
  `lint.NewRuleName` in array form).

In addition to that, when providing protobuf examples, it is often useful to
mark one as being explicitly "incorrect" (or "correct"). Do this by beginning
the code block with a special comment:

```
// Incorrect.
message BadThing {
  // ...
}
```

If a proto block _begins with_ a comment that says only `Incorrect.` or
`Correct.`, it picks up different styling when viewing in GitHub Pages.

<!-- prettier-ignore-start -->
[aip]: https://aip.dev/
[go]: https://golang.org/
[`go.mod`]: https://github.com/googleapis/api-linter/blob/main/go.mod
[`problem`]: https://godoc.org/github.com/googleapis/api-linter/lint#Problem
[protoreflect]: https://godoc.org/github.com/jhump/protoreflect
[`rules.go`]: https://github.com/googleapis/api-linter/blob/main/rules/rules.go
[visitor pattern]: https://en.wikipedia.org/wiki/Visitor_pattern
<!-- prettier-ignore-end -->
