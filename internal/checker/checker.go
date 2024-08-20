package checker

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"minwa/internal/database"
)

const timeout = 5 * time.Second

// CheckEndpoint makes a HTTP request to url and returns the status code and response time.
func CheckEndpoint(url string) (int, time.Duration, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}

	client := http.Client{
		Timeout: timeout,
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, time.Since(start), nil
}

// ScheduleCheck runs CheckEndpoint every interval and persists the result in the database
func ScheduleCheck(ctx context.Context, db *sql.DB, interval time.Duration) {
	queries := database.New(db)
	var mu sync.Mutex // Create a mutex for synchronization

	for {
		slog.Info("checking endpoints")
		start := time.Now()

		endpoints, err := queries.EndpointsList(ctx)
		if err != nil {
			slog.Error("unable to list endpoints", "error", err)
			continue
		}

		var wg sync.WaitGroup

		for _, endpoint := range endpoints {
			wg.Add(1)
			go func(endpoint database.Endpoint) {
				defer wg.Done()

				status, responseTime, err := CheckEndpoint(endpoint.Url)
				if err != nil {
					slog.Error("unable to check endpoint", "error", err)
					return
				}

				// Lock the mutex before writing to the database
				mu.Lock()
				defer mu.Unlock()

				if err := queries.ChecksCreate(context.Background(), database.ChecksCreateParams{
					EndpointID:   endpoint.ID,
					Status:       int64(status),
					ResponseTime: int64(responseTime.Milliseconds()),
				}); err != nil {
					slog.Error("unable to create check", "error", err)
				}
				slog.Info("check done", "endpoint", endpoint.Url, "status", status, "response_time", responseTime)
			}(endpoint) // Pass the endpoint to the goroutine
		}

		wg.Wait() // Wait for all goroutines to finish

		slog.Info("checking endpoints done", "duration", time.Since(start))
		time.Sleep(interval)
	}
}
