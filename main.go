package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

// the structure of the health check result from a server.
type HealthCheck struct {
	Application  string `json:"application"`
	Version      string `json:"version"`
	Uptime       int64  `json:"uptime"`
	RequestCount int64  `json:"requestCount"`
	ErrorCount   int64  `json:"errorCount"`
	SuccessCount int64  `json:"successCount"`
}

// represents the aggregated health check report for an application version.
type Report struct {
	Application string  `json:"application"`
	Version     string  `json:"version"`
	SuccessRate float64 `json:"successRate"`
}

const workerCount = 10 // Number of concurrent workers

func main() {
	servers, err := readServerList("server.txt")
	if err != nil {
		log.Fatalf("Error reading server list: %v", err)
	}

	client := &http.Client{}
	results := fetchHealthChecks(client, servers)
	aggregated := aggregateResults(results)

	printReport(aggregated)
	if err := saveReport(aggregated, "report.json"); err != nil {
		log.Fatalf("Error saving report: %v", err)
	}
}

// To read list of server addresses file.
func readServerList(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return strings.Split(strings.TrimSpace(string(data)), "\n"), nil
}

// Concurrently fetches health check data from the servers.
func fetchHealthChecks(client *http.Client, servers []string) []HealthCheck {
	var wg sync.WaitGroup
	results := make(chan HealthCheck, len(servers))
	serverChan := make(chan string, len(servers))

	for _, server := range servers {
		serverChan <- server
	}
	close(serverChan)

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for server := range serverChan {
				if hc, err := checkHealth(client, "http://"+server+"/healthz"); err == nil {
					results <- hc
				} else {
					log.Printf("Error querying server %s: %v", server, err)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var healthChecks []HealthCheck
	for hc := range results {
		healthChecks = append(healthChecks, hc)
	}
	return healthChecks
}

// performs a health check on the given server URL.
func checkHealth(client *http.Client, url string) (HealthCheck, error) {
	resp, err := client.Get(url)
	if err != nil {
		return HealthCheck{}, fmt.Errorf("failed to GET %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HealthCheck{}, fmt.Errorf("non-200 response from %s: %d", url, resp.StatusCode)
	}

	var hc HealthCheck
	if err := json.NewDecoder(resp.Body).Decode(&hc); err != nil {
		return HealthCheck{}, fmt.Errorf("failed to decode response from %s: %w", url, err)
	}
	return hc, nil
}

// aggregates the HC results by application and version.
func aggregateResults(results []HealthCheck) map[string]map[string]Report {
	agg := make(map[string]map[string]Report)
	for _, hc := range results {
		appVersions := agg[hc.Application]
		if appVersions == nil {
			appVersions = make(map[string]Report)
			agg[hc.Application] = appVersions
		}
		r := appVersions[hc.Version]
		r.Application = hc.Application
		r.Version = hc.Version
		r.SuccessRate = float64(hc.SuccessCount) / float64(hc.RequestCount)
		appVersions[hc.Version] = r
	}
	return agg
}

// prints the aggregated HC report
func printReport(aggregated map[string]map[string]Report) {
	for app, versions := range aggregated {
		for ver, report := range versions {
			fmt.Printf("Application: %s, Version: %s, Success Rate: %.2f%%\n",
				app, ver, report.SuccessRate*100)
		}
	}
}

// saves the aggregated health check report to a JSON file.
func saveReport(aggregated map[string]map[string]Report, filename string) error {
	data, err := json.MarshalIndent(aggregated, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write report to file %s: %w", filename, err)
	}
	return nil
}
