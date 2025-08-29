#!/bin/bash

# =============================================================================
# CONFIGURATION - CHANGE THESE VARIABLES AS NEEDED
# =============================================================================

# Name of your proto files (without the .proto extension)
PROTO_FILE_1="cluster_leader"    # Change this to your first proto filename
PROTO_FILE_2="cluster_follower"   # Change this to your second proto filename

# Output directory for generated Go files
OUTPUT_DIR="./clustering/grpc"                 # Change this to your desired output directory

# Proto file directory (where your .proto files are located)
PROTO_DIR="./proto"                      # Change this if your proto files are elsewhere

# =============================================================================
# END OF CONFIGURATION - USUALLY NO CHANGES NEEDED BELOW THIS LINE
# =============================================================================

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc (Protocol Buffers compiler) is not installed."
    echo "Please install it from: https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# Check if Go plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Error: protoc-gen-go is not installed."
    echo "Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Error: protoc-gen-go-grpc is not installed."
    echo "Install with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Check if proto files exist
if [ ! -f "$PROTO_DIR/$PROTO_FILE_1.proto" ]; then
    echo "Error: Proto file $PROTO_DIR/$PROTO_FILE_1.proto not found!"
    exit 1
fi

if [ ! -f "$PROTO_DIR/$PROTO_FILE_2.proto" ]; then
    echo "Error: Proto file $PROTO_DIR/$PROTO_FILE_2.proto not found!"
    exit 1
fi

echo "Compiling proto files..."
echo "Proto directory: $PROTO_DIR"
echo "Output directory: $OUTPUT_DIR"

# Compile first proto file
echo "Compiling $PROTO_FILE_1.proto..."
protoc --proto_path="$PROTO_DIR" \
       --go_out="$OUTPUT_DIR" --go_opt=paths=source_relative \
       --go-grpc_out="$OUTPUT_DIR" --go-grpc_opt=paths=source_relative \
       "$PROTO_FILE_1.proto"

# Check if compilation was successful
if [ $? -eq 0 ]; then
    echo "✓ Successfully compiled $PROTO_FILE_1.proto"
else
    echo "✗ Failed to compile $PROTO_FILE_1.proto"
    exit 1
fi

# Compile second proto file
echo "Compiling $PROTO_FILE_2.proto..."
protoc --proto_path="$PROTO_DIR" \
       --go_out="$OUTPUT_DIR" --go_opt=paths=source_relative \
       --go-grpc_out="$OUTPUT_DIR" --go-grpc_opt=paths=source_relative \
       "$PROTO_FILE_2.proto"

# Check if compilation was successful
if [ $? -eq 0 ]; then
    echo "✓ Successfully compiled $PROTO_FILE_2.proto"
else
    echo "✗ Failed to compile $PROTO_FILE_2.proto"
    exit 1
fi

echo ""
echo "Compilation completed successfully!"
echo "Generated files are in: $OUTPUT_DIR"