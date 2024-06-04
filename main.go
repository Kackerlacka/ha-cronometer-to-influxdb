package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"io"

	"github.com/influxdata/influxdb1-client/v2"
	"github.com/jrmycanady/gocronometer"
)

func main() {
	// Set up InfluxDB client
	influxdbURL := os.Getenv("INFLUXDB_URL")
	influxdbUsername := os.Getenv("INFLUXDB_USERNAME")
	influxdbPassword := os.Getenv("INFLUXDB_PASSWORD")
	influxdbDatabase := os.Getenv("INFLUXDB_DATABASE")

	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxdbURL,
		Username: influxdbUsername,
		Password: influxdbPassword,
	})
	if err != nil {
		fmt.Printf("failed to create InfluxDB client: %s\n", err)
		return
	}
	defer influxClient.Close()
	
    // Define start time for data deletion (midnight one week ago)
    startDate := time.Now().AddDate(0, 0, -7) // Subtract 7 days
    startDate = startDate.Truncate(24 * time.Hour) // Truncate time to midnight
    endDate := time.Now()

    // Delete data within the specified time range
    err = deleteDataInRange(influxClient, influxdbDatabase, startDate, endDate)
    if err != nil {
        fmt.Printf("failed to delete data: %s\n", err)
        return
    }


	// Set up Cronometer client
	c := gocronometer.NewClient(nil)
	cronometerUsername := os.Getenv("CRONOMETER_USERNAME")
	cronometerPassword := os.Getenv("CRONOMETER_PASSWORD")


	// Login to Cronometer
	err = c.Login(context.Background(), cronometerUsername, cronometerPassword)
	if err != nil {
		fmt.Printf("failed to login to Cronometer: %s\n", err)
		return
	}

	// Retrieve biometrics data from Cronometer
	biometricsData, err := c.ExportBiometrics(context.Background(), startDate, endDate)
	if err != nil {
		fmt.Printf("failed to retrieve biometrics data: %s\n", err)
		return
	}

	// Export raw biometrics data to CSV
	err = exportRawBiometricsToCSV(biometricsData, "raw_biometrics.csv")
	if err != nil {
		fmt.Printf("failed to export raw biometrics data to CSV: %s\n", err)
		return
	}

	// Format biometrics data for InfluxDB
	formattedBiometricsData, err := formatBiometricsForInfluxDB(biometricsData)
	if err != nil {
		fmt.Printf("failed to format biometrics data: %s\n", err)
		return
	}

	// Export formatted biometrics data to CSV
	err = exportFormattedBiometricsToCSV(formattedBiometricsData, "formatted_biometrics.csv")
	if err != nil {
		fmt.Printf("failed to export formatted biometrics data to CSV: %s\n", err)
		return
	}

	// Write biometrics data to InfluxDB
	err = writeDataToInfluxDB(influxClient, influxdbDatabase, formattedBiometricsData)
	if err != nil {
		fmt.Printf("failed to write data to InfluxDB: %s\n", err)
		return
	}

	// Retrieve daily nutrition data from Cronometer
	dailyNutritionData, err := c.ExportDailyNutrition(context.Background(), startDate, endDate)
	if err != nil {
		fmt.Printf("failed to retrieve daily nutrition data: %s\n", err)
		return
	}

	// Export raw daily nutrition data to CSV
	err = exportRawDailyNutritionToCSV(dailyNutritionData, "raw_daily_nutrition.csv")
	if err != nil {
		fmt.Printf("failed to export raw daily nutrition data to CSV: %s\n", err)
		return
	}

	// Format daily nutrition data for InfluxDB
	formattedNutritionData, err := formatDailyNutritionForInfluxDB(dailyNutritionData)
	if err != nil {
		fmt.Printf("failed to format daily nutrition data: %s\n", err)
		return
	}

	// Export formatted daily nutrition data to CSV
	err = exportFormattedDailyNutritionToCSV(formattedNutritionData, "formatted_daily_nutrition.csv")
	if err != nil {
		fmt.Printf("failed to export formatted daily nutrition data to CSV: %s\n", err)
		return
	}

	// Write daily nutrition data to InfluxDB
	err = writeDataToInfluxDB(influxClient, influxdbDatabase, formattedNutritionData)
	if err != nil {
		fmt.Printf("failed to write data to InfluxDB: %s\n", err)
		return
	}

	fmt.Println("Biometrics and daily nutrition data written to InfluxDB successfully.")
}

// Modify the formatBiometricsForInfluxDB function to handle splitting Blood Pressure
func formatBiometricsForInfluxDB(biometricsData string) ([]*client.Point, error) {
    var points []*client.Point
    reader := csv.NewReader(strings.NewReader(biometricsData))

    // Read CSV header
    header, err := reader.Read()
    if err != nil {
        return nil, fmt.Errorf("failed to read CSV header: %w", err)
    }

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("failed to read CSV record: %w", err)
        }

        // Map the CSV record to the header
        data := make(map[string]string)
        for i, value := range record {
            data[header[i]] = value
        }

        // Extract metric and unit from the record
        metric := data["Metric"]
        unit := data["Unit"]

        // Handle blood pressure data
        if metric == "Blood Pressure" {
            // Parse blood pressure into systolic and diastolic values
            systolic, diastolic, err := parseBloodPressure(data["Amount"])
            if err != nil {
                return nil, fmt.Errorf("failed to parse blood pressure: %w", err)
            }

            // Create InfluxDB point for systolic value
            point, err := createInfluxDBPoint(data, "Systolic Blood Pressure", float64(systolic), unit, false, "biometrics")
            if err != nil {
                return nil, fmt.Errorf("failed to create InfluxDB point for systolic value: %w", err)
            }
            points = append(points, point)

            // Create InfluxDB point for diastolic value
            point, err = createInfluxDBPoint(data, "Diastolic Blood Pressure", float64(diastolic), unit, false, "biometrics")
            if err != nil {
                return nil, fmt.Errorf("failed to create InfluxDB point for diastolic value: %w", err)
            }
            points = append(points, point)
        } else {
            // Create InfluxDB point for other biometrics
            amount, err := strconv.ParseFloat(data["Amount"], 64)
            if err != nil {
                return nil, fmt.Errorf("failed to parse amount as float64: %w", err)
            }
            point, err := createInfluxDBPoint(data, metric, amount, unit, false, "biometrics")
            if err != nil {
                return nil, fmt.Errorf("failed to create InfluxDB point: %w", err)
            }
            points = append(points, point)
        }
    }

    return points, nil
}


func formatDailyNutritionForInfluxDB(nutritionData string) ([]*client.Point, error) {
    var points []*client.Point
    reader := csv.NewReader(strings.NewReader(nutritionData))

    // Read CSV header
    header, err := reader.Read()
    if err != nil {
        return nil, fmt.Errorf("failed to read CSV header: %w", err)
    }
    fmt.Println("Header:", header)

    // Extract metric names and units from headers
    var metrics []string
    var units []string
    for _, column := range header {
        parts := strings.Split(column, " (")
        metric := parts[0] // Extract metric name before "("
        metrics = append(metrics, metric)
        if len(parts) > 1 {
            unit := strings.TrimSuffix(parts[1], ")") // Extract unit inside parentheses
            units = append(units, unit)
        } else {
            units = append(units, "") // No unit specified
        }
    }

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, fmt.Errorf("failed to read CSV record: %w", err)
        }
        fmt.Println("Record:", record)

        // Skip the record if it's the "Completed" header
        if record[0] == "Completed" {
            fmt.Println("Skipping Completed header")
            continue
        }

        // Map the CSV record to the header
        data := make(map[string]string)
        for i, value := range record {
            data[metrics[i]] = value // Use extracted metric names
        }
        fmt.Println("Mapped data:", data)

        // Create InfluxDB point for each metric in daily nutrition
        date := data["Date"]
        for i, metric := range metrics {
            if metric == "Date" || metric == "Completed" {
                continue
            }
            // Skip processing if the value is empty
            if data[metric] == "" {
                continue
            }
            amount, err := strconv.ParseFloat(data[metric], 64)
            if err != nil {
                return nil, fmt.Errorf("failed to parse amount for %s as float64: %w", metric, err)
            }
            data["Day"] = date // Add date to data map for createInfluxDBPoint
            unit := units[i]  // Get corresponding unit for the metric
            point, err := createInfluxDBPoint(data, metric, amount, unit, false, "nutrition")
            if err != nil {
                return nil, fmt.Errorf("failed to create InfluxDB point for %s: %w", metric, err)
            }
            points = append(points, point)
        }
    }

    return points, nil
}


func getSource(metric string) string {
	if strings.HasSuffix(metric, " (Health Connect)") {
		return "automatic"
	}
	return "manual"
}

func getFinalMetric(metric string) string {
	trimmed := strings.TrimSpace(metric)
	if strings.HasSuffix(trimmed, " (Health Connect)") {
		return strings.TrimSuffix(trimmed, " (Health Connect)")
	}
	return trimmed
}
 
func writeDataToInfluxDB(influxClient client.Client, influxdbDatabase string, points []*client.Point) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  influxdbDatabase,
		Precision: "s",
	})
	if err != nil {
		return fmt.Errorf("failed to create batch points: %w", err)
	}

	bp.AddPoints(points)

	if err := influxClient.Write(bp); err != nil {
		return fmt.Errorf("failed to write batch points: %w", err)
	}

	return nil
}

func deleteDataInRange(influxClient client.Client, influxdbDatabase string, startTime, endTime time.Time) error {
    // Construct the InfluxDB query to delete data within the specified time range
    query := fmt.Sprintf("DELETE FROM /./ WHERE time >= '%s' AND time <= '%s'", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
    
    // Create the query object
    q := client.NewQuery(query, influxdbDatabase, "")
    
    // Execute the query
    response, err := influxClient.Query(q)
    if err != nil || response.Error() != nil {
        return fmt.Errorf("failed to delete data: %w", err)
    }
    
    return nil
}


// Helper function to create an InfluxDB point
func createInfluxDBPoint(data map[string]string, metric string, amount float64, unit string, hasTime bool, measurement string) (*client.Point, error) {
	var timestamp time.Time
	var err error

	if hasTime {
		if data["Time"] != "" {
			timestamp, err = parseTimeWithAMPM(data["Day"], data["Time"])
		} else {
			timestamp, err = time.Parse("2006-01-02", data["Day"])
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp with time: %w", err)
		}
	} else {
		timestamp, err = time.Parse("2006-01-02", data["Day"])
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp without time: %w", err)
		}
	}

	fields := map[string]interface{}{
		"Amount": amount,
	}

	tags := map[string]string{
		"Metric": metric,
		"Unit":   unit, // Include unit in tags
		"Source": getSource(metric),
	}

	return client.NewPoint(
		measurement, // Measurement name
		tags,
		fields,
		timestamp,
	)
}


// Helper function to parse time with AM/PM indicator
func parseTimeWithAMPM(date string, timeStr string) (time.Time, error) {
	// Try to parse with AM/PM format
	timestamp, err := time.Parse("2006-01-02 3:04 PM", date+" "+timeStr)
	if err == nil {
		return timestamp, nil
	}

	// Try to parse with 24-hour format
	timestamp, err = time.Parse("2006-01-02 15:04:05", date+" "+timeStr)
	if err == nil {
		return timestamp, nil
	}

	return time.Time{}, fmt.Errorf("failed to parse time with AM/PM: %w", err)
}

// Helper function to parse Blood Pressure into Systolic and Diastolic values
func parseBloodPressure(bp string) (int, int, error) {
	parts := strings.Split(bp, "/")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid blood pressure format")
	}

	systolicFloat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse systolic value: %w", err)
	}
	systolic := int(math.Round(systolicFloat))

	diastolicFloat, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse diastolic value: %w", err)
	}
	diastolic := int(math.Round(diastolicFloat))

	return systolic, diastolic, nil
}

// Export raw biometrics data to CSV before formatting
func exportRawBiometricsToCSV(data, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

// Export formatted biometrics data to CSV after formatting
func exportFormattedBiometricsToCSV(data []*client.Point, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write data points to CSV in InfluxDB line protocol format
    for _, point := range data {
        timestamp := strconv.FormatInt(point.Time().UnixNano(), 10)
        tags := point.Tags()
        fields, _ := point.Fields()

        line := fmt.Sprintf("biometrics,Metric=%s,Unit=%s,Source=%s Amount=%v %s",
            tags["Metric"],
            tags["Unit"],
            tags["Source"],
            fields["Amount"],
            timestamp,
        )

        err := writer.Write([]string{line})
        if err != nil {
            return err
        }
    }

    return nil
}

// Export raw daily nutrition data to CSV before formatting
func exportRawDailyNutritionToCSV(data, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

// Export formatted daily nutrition data to CSV after formatting
func exportFormattedDailyNutritionToCSV(data []*client.Point, filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Write data points to CSV in InfluxDB line protocol format
    for _, point := range data {
        timestamp := strconv.FormatInt(point.Time().UnixNano(), 10)
        tags := point.Tags()
        fields, _ := point.Fields()

        line := fmt.Sprintf("nutrition,Metric=%s,Unit=%s,Source=%s Amount=%v %s",
            tags["Metric"],
            tags["Unit"],
            tags["Source"],
            fields["Amount"],
            timestamp,
        )

        err := writer.Write([]string{line})
        if err != nil {
            return err
        }
    }

    return nil
}
