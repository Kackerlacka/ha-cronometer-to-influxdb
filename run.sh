#!/bin/sh

# Start the dcron service
crond

# List the contents of the directory
ls -l

# Run your Go application
./cronapp -cronometer_username "$CRONOMETER_USERNAME" -cronometer_password "$CRONOMETER_PASSWORD" -influxdb_url "$INFLUXDB_URL" -influxdb_username "$INFLUXDB_USERNAME" -influxdb_password "$INFLUXDB_PASSWORD" -influxdb_database "$INFLUXDB_DATABASE"
