# grafana-exporter

`grafana-exporter` exports all of the dashboards in a remote Grafana instance to JSON files.

## Usage

        go run main.go -path ./dashboards -token <your-grafana-api-token> -url http://your-grafana-server:3000
