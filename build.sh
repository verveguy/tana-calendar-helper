#!/bin/bash

# Ensure the script stops on the first error
set -e

BINARY_NAME="tana-calendar-helper"

# Clean up first
rm -rf bin
rm -rf ${BINARY_NAME}
rm -f ${BINARY_NAME}.zip

mkdir -p bin/scripts

# Compile for macOS
echo "Building for macOS arm64..."
ARCH=arm64
GOOS=darwin GOARCH=${ARCH} go build -o bin/${BINARY_NAME}.${ARCH} ${BINARY_NAME}.go
swiftc scripts/getcalendar.swift -o bin/scripts/getcalendar.${ARCH} -target arm64-apple-macosx10.15

echo "Building for macOS amd64..."
ARCH=amd64
GOOS=darwin GOARCH=${ARCH} go build -o bin/${BINARY_NAME}.${ARCH} ${BINARY_NAME}.go
swiftc scripts/getcalendar.swift -o bin/scripts/getcalendar.${ARCH} -target x86_64-apple-macosx10.15

echo "Packaging universal binary..."
lipo -create -output bin/${BINARY_NAME} bin/${BINARY_NAME}.arm64 bin/${BINARY_NAME}.amd64
lipo -create -output bin/scripts/getcalendar bin/scripts/getcalendar.arm64 bin/scripts/getcalendar.amd64

echo "Removing arch builds..."
rm bin/${BINARY_NAME}.arm64 bin/${BINARY_NAME}.amd64
rm bin/scripts/*.arm64 bin/scripts/*.amd64

# Compile for Windows
# echo "Building for Windows..."
# GOOS=windows GOARCH=386 go build -o bin/${BINARY_NAME}.exe ${BINARY_NAME}.go

# Compile for Linux
# echo "Building for Linux..."
# GOOS=linux GOARCH=386 go build -o bin/${BINARY_NAME}.linux ${BINARY_NAME}.go

#echo "Staging required scripts..."
cp scripts/calendar_auth.scpt bin/scripts
chmod +x bin/scripts/*

echo "Preparing zip"
cp -rp bin ${BINARY_NAME}
zip -r ${BINARY_NAME} ${BINARY_NAME}

# cleanup
rm -rf ${BINARY_NAME}

echo "Build complete! Binaries are in the 'bin' directory."
