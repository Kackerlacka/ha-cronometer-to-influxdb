#!/bin/sh

# Start the dcron service
crond

# List the contents of the directory
ls -l

# Print environment variables to verify they are set
echo "CRONOMETER_USERNAME: $CRONOMETER_USERNAME"
echo "CRONOMETER_PASSWORD: $CRONOMETER_PASSWORD"
echo "INFLUXDB_URL: $INFLUXDB_URL"
echo "INFLUXDB_USERNAME: $INFLUXDB_USERNAME"
echo "INFLUXDB_PASSWORD: $INFLUXDB_PASSWORD"
echo "INFLUXDB_DATABASE: $INFLUXDB_DATABASE"

# Run your Go application
./cronapp -cronometer_username "$CRONOMETER_USERNAME" -cronometer_password "$CRONOMETER_PASSWORD" -influxdb_url "$INFLUXDB_URL" -influxdb_username "$INFLUXDB_USERNAME" -influxdb_password "$INFLUXDB_PASSWORD" -influxdb_database "$INFLUXDB_DATABASE"
