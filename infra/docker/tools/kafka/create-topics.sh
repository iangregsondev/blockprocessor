#!/bin/bash

# Wait for Kafka to be ready
echo "Waiting for Kafka to be ready..."
while ! kafka-topics.sh --bootstrap-server kafka:9092 --list >/dev/null 2>&1; do
  sleep 1
done
echo "Kafka is up and running"

# Define topics to create
TOPICS=(
  # Block topics
  "bitcoin-blocks"
  "ethereum-blocks"
  "solana-blocks"
  
  # Transaction topics
  "bitcoin-transactions"
  "ethereum-transactions"
  "solana-transactions"
  
  # Filtered transaction topics
  "bitcoin-transactions-filtered"
  "ethereum-transactions-filtered"
  "solana-transactions-filtered"
)

for TOPIC in "${TOPICS[@]}"; do
  # Check if the topic already exists
  if kafka-topics.sh --bootstrap-server kafka:9092 --list | grep -q "^${TOPIC}$"; then
    echo "Topic ${TOPIC} already exists, skipping creation"
  else
    # Create the topic if it doesn't exist
    kafka-topics.sh --create --bootstrap-server kafka:9092 --replication-factor 1 --partitions 1 --topic "${TOPIC}"
    echo "Topic ${TOPIC} does not exist, created topic"
  fi
done