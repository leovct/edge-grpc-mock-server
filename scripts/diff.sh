#!/bin/bash

# Check diffs in git repo.

git --no-pager diff

if [[ $(git status --porcelain=v1 2>/dev/null) ]]; then
  echo "❌ Error: gRPC service is not up to date. Please run 'make gen' and commit your changes."
  #exit 1
else
  echo "✅ gRPC service is up to date."
fi
