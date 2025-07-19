#!/bin/bash
# build-all.sh - Build api-linter and its plugin together

set -e

# Go version info
GO_VERSION=$(go version)
echo "Building with $GO_VERSION"

# Build api-linter
echo "Building api-linter..."
go build -o api-linter-with-plugins ./cmd/api-linter

# Save the build info for verification
BUILT_WITH=$(go version -m ./api-linter-with-plugins)
echo "Built api-linter with:"
echo "$BUILT_WITH"

# Build the plugin
echo -e "\nBuilding plugin..."
cd julioz-api-rules

# Create a clean go.mod that matches exactly
cat > go.mod << EOF
module github.com/julioz/julioz-api-rules

go 1.23

require github.com/googleapis/api-linter v0.0.0-local

replace github.com/googleapis/api-linter => ../
EOF

# Build plugin
go mod tidy
go build -buildmode=plugin -o julioz-api-rules.so

# Verify plugin build info
PLUGIN_INFO=$(go version -m julioz-api-rules.so)
echo "Built plugin with:"
echo "$PLUGIN_INFO"

# Return to main directory
cd ..

# Test with a sample proto file
if [ -f "julioz-api-rules/test.proto" ]; then
  echo -e "\nTesting plugin with sample proto file..."
  ./api-linter-with-plugins --rule-plugin=julioz-api-rules/julioz-api-rules.so julioz-api-rules/test.proto
  TEST_EXIT=$?
  if [ $TEST_EXIT -eq 0 ]; then
    echo "✓ Plugin test successful"
  else
    echo "✗ Plugin test failed with exit code $TEST_EXIT"
  fi
fi

echo -e "\nBuild complete"
echo "API Linter binary with plugin support: ./api-linter-with-plugins"
echo "Custom rule plugin: julioz-api-rules/julioz-api-rules.so"
echo
echo "Run linter with: ./api-linter-with-plugins --rule-plugin=julioz-api-rules/julioz-api-rules.so YOUR-PROTO-FILE.proto" 