#!/bin/bash
# build.sh - Build script for the Julio API Linter Rules plugin

set -e

# Check Go version
GO_VERSION=$(go version)
echo "Building plugin with $GO_VERSION"

# Initialize the module if it doesn't exist
if [ ! -f "go.mod" ]; then
    echo "Initializing Go module"
    go mod init github.com/julioz/julioz-api-rules
fi

# Get the exact versions of required dependencies from api-linter
# You should use the same versions as your target api-linter installation
API_LINTER_VERSION="v1.69.2"  # Current version from internal/version.go
echo "Using api-linter version $API_LINTER_VERSION"

# Ensure dependencies are in sync with api-linter
go get github.com/googleapis/api-linter@$API_LINTER_VERSION

# Tidy the dependencies
go mod tidy

# Set the output file name
OUTPUT_FILE="julioz-api-rules.so"

# Build the plugin
echo "Building plugin..."
go build -buildmode=plugin -o $OUTPUT_FILE

# Get the absolute path of the built plugin
PLUGIN_PATH=$(realpath $OUTPUT_FILE)

echo "Plugin built: $PLUGIN_PATH" 