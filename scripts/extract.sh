#!/bin/bash

# Extract block and trace mock files from an archive.
archive=$1

echo "🔍 Verifying the data set..."
blocks=$(tar -tf $archive | grep -c block_)
traces=$(tar -tf $archive | grep -c trace_)
if [ "$blocks" -eq "$traces" ]; then
  echo "✅ Data set is correct"
else
  echo "❌ Aborting: different number of blocks ($blocks) and traces ($traces)."
  #exit 1
fi

echo "📦 Extracting mock files..."
rm -rf ./mocks
mkdir -p ./mocks ./mocks/blocks ./mocks/traces
tar -xf $archive -C ./mocks/blocks --strip-components=1 $(tar -tf $archive | grep block_)
tar -xf $archive -C ./mocks/traces --strip-components=1 $(tar -tf $archive | grep trace_)

echo "🪄  Formatting mock files (it might take a while)..."
for block in ./mocks/blocks/*.json; do
  jq '.' "$block" > "$block.tmp"
  mv "$block.tmp" "$block"
done

for trace in ./mocks/traces/*.json; do
  jq '.' "$trace" > "$trace.tmp"
  mv "$trace.tmp" "$trace"
done
