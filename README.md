# ALB-Monitoring — Custom Prometheus Exporter for AWS Application Load Balancer

A custom Prometheus exporter written in Go that collects metrics from AWS Application Load Balancer via CloudWatch API.

## Flow

```
IAM User → AWS SDK v2 → CloudWatch.GetMetricStatistics() → Go Struct → prometheus.MustNewConstMetric() → /metrics
```

## Metrics Exposed

| Metric | Type | Description |
|--------|------|-------------|
| `alb_request_count_total` | Counter | Total number of requests |
| `alb_active_connections` | Counter | Active connections |
| `alb_target_response_time` | Gauge | Target response time in seconds |
| `alb_http_code_elb_5xx` | Counter | HTTP 5xx errors from ALB |
| `alb_http_code_target_4xx` | Counter | HTTP 4xx errors from targets |
| `alb_http_code_target_5xx` | Counter | HTTP 5xx errors from targets |
| `alb_healthy_host_count` | Gauge | Number of healthy targets |
| `alb_unhealthy_host_count` | Gauge | Number of unhealthy targets |

## Project Structure

```
ALB-Monitoring/
├── main.go           # HTTP server + AWS config + Prometheus registry
├── collector/
│   └── alb.go        # Prometheus Collector + CloudWatch API calls
└── go.mod
```

## Prerequisites

- Go 1.24+
- AWS CLI configured
- Prometheus

## 1. AWS IAM Setup

Create an IAM User with the following inline policy:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "cloudwatch:GetMetricStatistics",
        "cloudwatch:ListMetrics",
        "elasticloadbalancing:DescribeLoadBalancers",
        "elasticloadbalancing:DescribeTargetGroups",
        "elasticloadbalancing:DescribeTargetHealth"
      ],
      "Resource": "*"
    }
  ]
}
```

Create Access Key → **Application running outside AWS**

Configure named profile:

```bash
aws configure --profile alb-monitoring
```

Verify:

```bash
aws configure list --profile alb-monitoring
aws sts get-caller-identity --profile alb-monitoring
aws elbv2 describe-load-balancers --profile alb-monitoring
```

## 2. Installation

```bash
git clone https://github.com/<your-username>/ALB-Monitoring.git
cd ALB-Monitoring
go mod tidy
```

Edit `collector/alb.go` and set your ALB name:

```go
Names: []string{"your-alb-name"},
```

Build:

```bash
go build -o alb-exporter .
```

## 3. Run

```bash
./alb-exporter
```

Test:

```bash
curl http://localhost:2111/metrics
```

## 4. Systemd Service

```bash
sudo vim /etc/systemd/system/alb-exporter.service
```

```ini
[Unit]
Description=ALB Prometheus Exporter
After=network.target

[Service]
Type=simple
User=<your-user>
ExecStart=/path/to/alb-exporter
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable alb-exporter
sudo systemctl start alb-exporter
sudo systemctl status alb-exporter
```

## 5. Prometheus Configuration

Add to `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: "alb-exporter"
    scrape_interval: 60s
    static_configs:
      - targets: ["localhost:2111"]
        labels:
          app: "alb-exporter"
```

```bash
sudo systemctl restart prometheus
```

Verify at: `http://localhost:9090/targets`

## References

- [AWS CloudWatch Metrics for ALB](https://docs.aws.amazon.com/elasticloadbalancing/latest/application/load-balancer-cloudwatch-metrics.html)
- [AWS SDK v2 for Go](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2)
- [Prometheus Go Client](https://github.com/prometheus/client_golang)
- [Prometheus Naming Conventions](https://prometheus.io/docs/practices/naming/)
- [Writing Exporters](https://prometheus.io/docs/instrumenting/writing_exporters/)