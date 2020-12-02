---
rule:
  aip: 152
  name: [core, '0152', request-message-name]
  summary: Run methods must have standardized request message names.
permalink: /152/request-message-name
redirect_from:
  - /0152/request-message-name
---

# Run methods: Request message

This rule enforces that all `Run` RPCs have a request message name
of `Run*JobRequest`, as mandated in [AIP-152][].

## Details

This rule looks at any message beginning with `Run`, and complains
if the name of the corresponding input message does not match the name of the
RPC with the suffix `Request` appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// Should be `RunWriteBookJobRequest`.
rpc RunWriteBookJob(WriteBookJob) returns (google.longrunning.Operation) {
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
// (-- api-linter: core::0152::request-message-name=disabled
//     aip.dev/not-precedent: We need to do this because reasons. --)
rpc RunWriteBookJob(WriteBookJob) returns (google.longrunning.Operation) {
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

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aip-152]: https://aip.dev/152
[aip.dev/not-precedent]: https://aip.dev/not-precedent
