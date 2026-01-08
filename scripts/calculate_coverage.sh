#!/bin/bash

# Calculate coverage only for business logic (what SonarCloud will analyze)
go tool cover -func=coverage.out | \
  grep -E "(controller_impl|persistence.*_impl|presenter_impl|use_case.*_impl|client_impl)" | \
  grep -v "New.*Impl" | \
  awk '{ total += $NF; count++ } END { if (count > 0) print "Business Logic Coverage: " total/count "%" }'
