#!/bin/bash

# Set target operating systems and architectures
TARGETS="darwin/amd64 darwin/arm64 linux/amd64 linux/arm linux/arm64 windows/amd64 windows/386"

# Loop through each target and build
for target in $TARGETS; do
    GOOS=${target%/*}
    GOARCH=${target#*/}
    OUTPUT="fsr_${GOOS}_${GOARCH}"
    
    # Adjust output filename for Windows
    if [ "$GOOS" = "windows" ]; then
        OUTPUT="$OUTPUT.exe"
    fi
    
    echo "Building for $GOOS/$GOARCH..."
    GOOS=$GOOS GOARCH=$GOARCH go build -o $OUTPUT
    pwd
done
