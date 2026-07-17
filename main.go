package main

import (
	"ALB-Monitoring/collector"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 1 - Load Config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("alb-monitoring"),
	)
	if err != nil {
		panic(fmt.Sprintf("failed loading config, %v", err))
	}

	// 2 - Create CloudWatch Clients from 'cfg'
	cwClient := cloudwatch.NewFromConfig(cfg)
	elbClient := elasticloadbalancingv2.NewFromConfig(cfg)

	// 3 - Prometheus Registry
	reg := prometheus.NewRegistry()
	// Collector from collector/ folder
	// reg.MustRegister( /* collector here */ )
	reg.MustRegister(collector.NewALBCollector(cwClient, elbClient))

	// 4 - HTTP server + /metric endpoint
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	// http.ListenAndServe(":2111", nil)
	log.Fatal(http.ListenAndServe(":2111", nil))
}
