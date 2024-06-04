#!/bin/sh

# Start the dcron service
crond

# Run your Go application
./go-app -cronometer_username "$CRONOMETER_USERNAME" -cronometer_password "$CRONOMETER_PASSWORD" -influxdb_url "$INFLUXDB_URL" -influxdb_username "$INFLUXDB_USERNAME" -influxdb_password "$INFLUXDB_PASSWORD" -influxdb_database "$INFLUXDB_DATABASE"
