#!/bin/sh

# Start the cron service
service cron start

# Get the environment variables
CRONOMETER_USERNAME=$CRONOMETER_USERNAME
CRONOMETER_PASSWORD=$CRONOMETER_PASSWORD
INFLUXDB_USERNAME=$INFLUXDB_USERNAME
INFLUXDB_PASSWORD=$INFLUXDB_PASSWORD
INFLUXDB_DATABASE=$INFLUXDB_DATABASE # New variable for InfluxDB database

# Run your Go application
./go-app -cronometer_username "$CRONOMETER_USERNAME" -cronometer_password "$CRONOMETER_PASSWORD" -influxdb_username "$INFLUXDB_USERNAME" -influxdb_password "$INFLUXDB_PASSWORD" -influxdb_database "$INFLUXDB_DATABASE"
