#!/usr/bin/env bash

# This script downloads ONNX Runtime binaries and verifies their SHA256 checksums.

set -e

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

download_and_verify() {
  local url=$1
  local expected_sha256=$2
  local filename=$(basename "$url")

  if [ ! -f "$filename" ]; then
    echo "downloading $filename..."
    curl -fL --progress-bar -o "$filename" "$url"
  else
    echo "$filename already exists, skipping download."
  fi

  echo "verifying $filename..."
  echo "$expected_sha256  $filename" | sha256sum -c -
  echo ""
}

# List of files and their SHA256 sums
FILES=(
  "https://github.com/microsoft/onnxruntime/releases/download/v1.23.2/onnxruntime-linux-aarch64-1.23.2.tgz 7c63c73560ed76b1fac6cff8204ffe34fe180e70d6582b5332ec094810241e5c"
  "https://github.com/microsoft/onnxruntime/releases/download/v1.23.2/onnxruntime-linux-x64-1.23.2.tgz 1fa4dcaef22f6f7d5cd81b28c2800414350c10116f5fdd46a2160082551c5f9b"
  "https://github.com/microsoft/onnxruntime/releases/download/v1.23.2/onnxruntime-osx-arm64-1.23.2.tgz b4d513ab2b26f088c66891dbbc1408166708773d7cc4163de7bdca0e9bbb7856"
  "https://github.com/microsoft/onnxruntime/releases/download/v1.23.2/onnxruntime-osx-x86_64-1.23.2.tgz d10359e16347b57d9959f7e80a225a5b4a66ed7d7e007274a15cae86836485a6"
  "https://github.com/microsoft/onnxruntime/releases/download/v1.23.2/onnxruntime-win-arm64-1.23.2.zip 1cfe88b6435df3b5fb0e9f6bd7d6f5df1e887b6174de7f6e2a47bab956f3f168"
  "https://github.com/microsoft/onnxruntime/releases/download/v1.23.2/onnxruntime-win-x64-1.23.2.zip 0b38df9af21834e41e73d602d90db5cb06dbd1ca618948b8f1d66d607ac9f3cd"
)

for entry in "${FILES[@]}"; do
  read -r url sha256 <<< "$entry"
  download_and_verify "$url" "$sha256"
done

echo "all files downloaded and verified successfully."
