#!/bin/sh

# Start the dcron service
crond

# Get the environment variables
CRONOMETER_USERNAME=$CRONOMETER_USERNAME
CRONOMETER_PASSWORD=$CRONOMETER_PASSWORD
INFLUXDB_URL=$INFLUXDB_URL
INFLUXDB_USERNAME=$INFLUXDB_USERNAME
INFLUXDB_PASSWORD=$INFLUXDB_PASSWORD
INFLUXDB_DATABASE=$INFLUXDB_DATABASE

# Run your Go application
./go-app -cronometer_username "$CRONOMETER_USERNAME" -cronometer_password "$CRONOMETER_PASSWORD" -influxdb_url "$INFLUXDB_URL" -influxdb_username "$INFLUXDB_USERNAME" -influxdb_password "$INFLUXDB_PASSWORD" -influxdb_database "$INFLUXDB_DATABASE"
