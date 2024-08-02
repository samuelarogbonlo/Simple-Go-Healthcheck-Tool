# Go-HealthCheck-Aggregator

This repository contains code that is designed to perform health checks on a list of server endpoints, aggregate the results, and generate a comprehensive report.

The tool queries each server's /healthz endpoint to collect data on the server's application name, version, uptime, request counts, error counts, and success counts. It then calculates success rates for each application version, prints a detailed report to the console, and saves the report in a structured JSON format.

## Prerequisities
- Golang
- Server.txt file

## Features
- Concurrent health checks on multiple servers.
- Aggregation of health check results by application and version.
- Calculation of success rates.
- Detailed report generation in both console output and JSON file.
- Makefile to automate building and running processes

## Installation And Testing
- Clone the repository
```
git clone https://github.com/samuelarogbonlo/cloud-.git
cd health-check-aggregator
```
- Build the project
```
make run
```
> **_Note:_**
If you wish to clean up your stack to re-run the processes then use the command below:
```
make clean
```

# Maintainers

[@samuelarogbonlo](https://github.com/samuelarogbonlo). Contributions are welcome! Please open an issue or submit a pull request with your improvements.

# License

Fully open source and dual-licensed under MIT and Apache 2.
