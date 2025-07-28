#!/bin/bash
#
# Start script for emergency-auth-code-api
PORT="20188"

# Read brokers and topics from environment and split on comma
IFS=',' read -ra BROKERS <<< "${KAFKA_BROKER_ADDR}"

echo Kafka broker addresses: "${KAFKA_BROKER_ADDR}"
echo "Kafka brokers: ${BROKERS[@]}"

# Ensure we only populate the broker address and topic via application arguments
unset KAFKA_BROKER_ADDR

exec ./emergency-auth-code-api "-bind-addr=:${PORT}" $(for broker in "${BROKERS[@]}"; do echo -n "-broker-addr=${broker} "; done)