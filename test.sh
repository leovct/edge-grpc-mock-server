#!/bin/bash

git update-index --really-refresh
git --no-pager diff

if $(git diff-index --quiet HEAD); then 
  echo "✅ gRPC service is up to date."
else
  echo "❌ Error: gRPC service is not up to date. Please run 'make gen' and commit your changes."
  exit 1
fi