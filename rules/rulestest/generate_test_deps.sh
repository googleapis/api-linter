#!/bin/bash

# Add imports to test_deps.proto and run this script. The output will be a
# file: test_deps.protoset, which will be a portable protoset containing
# all of the dependencies.

protoc --include_source_info --include_imports -otest_deps.protoset \
  -I$(p4 g4d -f testclient_golint_deps 2>/dev/null) -I$(pwd) test_deps.proto

g4 citc -d testclient_golint_deps >/dev/null