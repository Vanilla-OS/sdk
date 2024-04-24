#!/bin/bash

IMAGE_NAME="vanilla-sdk-test"

if command -v docker &> /dev/null; then
    CMD=docker
elif command -v podman &> /dev/null; then
    CMD=podman
else
    echo "Neither Docker nor Podman is installed. Install one and try again."
    exit 1
fi

echo "Building image with $CMD..."
$CMD build -t $IMAGE_NAME -f Containerfile .

if [ $? -ne 0 ]; then
    echo "Image build failed."
    exit 1
fi

echo "Running tests..."
$CMD run --name go_vos_sdk_test_container $IMAGE_NAME

echo "Removing the container..."
$CMD rm go_vos_sdk_test_container

if [ "$1" == "--remove-image" ]; then
    echo "Removing the image..."
    $CMD rmi $IMAGE_NAME
fi

echo "Operations completed."
