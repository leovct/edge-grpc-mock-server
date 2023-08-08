#!/bin/bash
git --no-pager diff

if [[ $(git status --porcelain) ]]; then
  echo "❌ Error: gRPC service is not up to date. Please run 'make gen' and commit your changes."
  #exit 1
else
  echo "✅ gRPC service is up to date."
fi

if [[ $(git diff --quiet && git diff --cached --quiet) ]]; then
  echo "❌ Error: gRPC service is not up to date. Please run 'make gen' and commit your changes."
  #exit 1
else
  echo "✅ gRPC service is up to date."
fi
