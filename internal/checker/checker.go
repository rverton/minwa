package checker

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"minwa/internal/database"
)

// CheckEndpoint makes a HTTP request to url and returns the status code and response time.
func CheckEndpoint(url string) (int, time.Duration, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, time.Since(start), nil
}

// ScheduleCheck runs CheckEndpoint every interval and persists the result in the database
func ScheduleCheck(ctx context.Context, db *sql.DB, interval time.Duration) {
	queries := database.New(db)
	for {
		slog.Info("checking endpoints")
		start := time.Now()

		endpoints, err := queries.EndpointsList(ctx)
		if err != nil {
			slog.Error("unable to list endpoints", "error", err)
			continue
		}

		for _, endpoint := range endpoints {
			status, responseTime, err := CheckEndpoint(endpoint.Url)
			if err != nil {
				slog.Error("unable to check endpoint", "error", err)
				continue
			}

			if err := queries.ChecksCreate(context.Background(), database.ChecksCreateParams{
				EndpointID:   endpoint.ID,
				Status:       int64(status),
				ResponseTime: int64(responseTime.Milliseconds()),
			}); err != nil {
				slog.Error("unable to create check", "error", err)
			}
			slog.Info("check done", "endpoint", endpoint.Url, "status", status, "response_time", responseTime)
		}

		slog.Info("checking endpoints done", "duration", time.Since(start))
		time.Sleep(interval)
	}
}
