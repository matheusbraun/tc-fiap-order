#!/bin/bash

# Script to generate test coverage reports locally
# Usage: ./scripts/coverage.sh

set -e

echo "🧪 Running tests with coverage..."
go test ./... -coverprofile=coverage.out -covermode=atomic

echo ""
echo "📊 Coverage summary:"
go tool cover -func=coverage.out

echo ""
echo "📈 Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "✅ Coverage reports generated:"
echo "   - coverage.out (for SonarCloud)"
echo "   - coverage.html (for local viewing)"
echo ""
echo "💡 Open coverage.html in your browser to see detailed coverage"
