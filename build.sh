#!/bin/bash

# Ensure the script stops on the first error
set -e

BINARY_NAME="tana-calendar-helper"

# Clean up first
rm -rf bin
rm -rf ${BINARY_NAME}
rm -f ${BINARY_NAME}.zip

# Compile for macOS
echo "Building for macOS arm64..."
GOOS=darwin GOARCH=arm64 go build -o bin/${BINARY_NAME}.arm64 ${BINARY_NAME}.go
echo "Building for macOS amd64..."
GOOS=darwin GOARCH=amd64 go build -o bin/${BINARY_NAME}.amd64 ${BINARY_NAME}.go

echo "Packaging universal binary..."
lipo -create -output bin/${BINARY_NAME} bin/${BINARY_NAME}.arm64 ./bin/${BINARY_NAME}.amd64

echo "Removing arch builds..."
rm bin/${BINARY_NAME}.arm64 ./bin/${BINARY_NAME}.amd64

# Compile for Windows
# echo "Building for Windows..."
# GOOS=windows GOARCH=386 go build -o bin/${BINARY_NAME}.exe ${BINARY_NAME}.go

# Compile for Linux
# echo "Building for Linux..."
# GOOS=linux GOARCH=386 go build -o bin/${BINARY_NAME}.linux ${BINARY_NAME}.go

echo "Staging required scripts..."
cp -r scripts bin
chmod +x bin/scripts/*

echo "Preparing zip"
mv bin ${BINARY_NAME}
zip -r ${BINARY_NAME} ${BINARY_NAME}

# cleanup
rm -rf ${BINARY_NAME}

echo "Build complete! Binaries are in the 'bin' directory."
