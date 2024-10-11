#!/bin/bash

# Name of the Docker network
NETWORK_NAME="deblock-external-network"

# Check if the network already exists
if docker network ls --format '{{.Name}}' | grep -wq "$NETWORK_NAME"; then
    echo "Docker network '$NETWORK_NAME' already exists, skipping creation."
else
    # Create the Docker network
    docker network create "$NETWORK_NAME"
    echo "Docker network '$NETWORK_NAME' has been created."
fi