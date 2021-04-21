---
rule:
  aip: 152
  name: [core, '0152', response-message-name]
  summary: Run methods must return a long-running operation.
permalink: /152/response-message-name
redirect_from:
  - /0152/response-message-name
---

# Run methods: Response message

This rule enforces that all `Run` RPCs return a long-running operation that
resolves to a response with a corresponding name, as mandated in [AIP-152][].

## Details

This rule looks at any method beginning with `Run`, and complains if the
response is not a long-running operation that resolves to a response matching
the name of the method with a `Response` suffix.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `google.longrunning.Operation`.
rpc RunWriteBookJob(RunWriteBookJobRequest) returns (RunWriteBookJobResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/writeBookJobs/*}:run"
    body: "*"
  };
}
```

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc RunWriteBookJob(RunWriteBookJobRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/writeBookJobs/*}:run"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    // Should be "RunWriteBookJobResponse".
    response_type: "WriteBookJob"
    metadata_type: "RunWriteBookJobMetadata"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc RunWriteBookJob(RunWriteBookJobRequest) returns (google.longrunning.Operation) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/writeBookJobs/*}:run"
    body: "*"
  };
  option (google.longrunning.operation_info) = {
    response_type: "RunWriteBookJobResponse"
    metadata_type: "RunWriteBookJobMetadata"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aip.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0152::response-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc RunWriteBookJob(RunWriteBookJobRequest) returns (RunWriteBookJobResponse) {
  option (google.api.http) = {
    post: "/v1/{name=publishers/*/writeBookJobs/*}:run"
    body: "*"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-152]: https://aip.dev/152
[aip.dev/not-precedent]: https://aip.dev/not-precedent
