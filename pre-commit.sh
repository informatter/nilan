#!/bin/bash

# Lints code and runs unit tests before allowing commits.
# To add this file as a pre-commit hook, run:
# `cp pre-commit.sh .git/hooks/pre-commit && chmod +x .git/hooks/pre-commit`

set -e

echo "🔍 Running pre-commit checks..."

# Get the project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$PROJECT_ROOT"

echo "📝 Checking code formatting..."
if ! go fmt ./...; then
    echo "❌ Code formatting failed"
    exit 1
fi

echo "🔍 Running go vet..."
if ! go vet ./...; then
    echo "❌ Code analysis (go vet) failed"
    exit 1
fi

echo "🧪 Running unit tests..."
if ! go test -v ./...; then
    echo "❌ Unit tests failed"
    exit 1
fi

echo "✅ All pre-commit checks passed!"
exit 0