#!/bin/bash
set -e

coverageOut='.coverage/coverage.cov'
coverageReport='.coverage/index.html'

# Run test and generate coverage output an report
go test -v -coverprofile "$coverageOut" ./...
echo "=== RUN   Coverage threshold"
go tool cover -html "$coverageOut" -o "$coverageReport"

# Coverage threshold check, this can be set in env var
coverageThreshold="${COVERAGE_THRESHOLD:=70.0}"
totalCoverage=$(go tool cover -func="$coverageOut" | grep total | grep -Eo '[0-9]+\.[0-9]+')

if (( $(echo "$totalCoverage $coverageThreshold" | awk '{print ($1 < $2)}') )); then
    echo "--- FAIL: Coverage at $totalCoverage%, out of $coverageThreshold%" >&2
    echo "FAIL" >&2
    exit 1
else
    echo "--- PASS: Coverage at $totalCoverage%, out of $coverageThreshold%"
    echo "PASS"
fi
