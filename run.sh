#!/bin/bash

# Start the cron service
service cron start

# Load bashio functions
source /usr/lib/hassio-addons/bashio

# Get the environment variables
CRONOMETER_USERNAME=$(bashio::config 'cronometer_username')
CRONOMETER_PASSWORD=$(bashio::config 'cronometer_password')
INFLUXDB_USERNAME=$(bashio::config 'influxdb_username')
INFLUXDB_PASSWORD=$(bashio::config 'influxdb_password')

# Run your Go application
./go-app -cronometer_username "$CRONOMETER_USERNAME" -cronometer_password "$CRONOMETER_PASSWORD" -influxdb_username "$INFLUXDB_USERNAME" -influxdb_password "$INFLUXDB_PASSWORD"
