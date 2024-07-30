package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"

	"github.com/upstars-global/domains-expiration-exporter/internal/cf"
	"github.com/upstars-global/domains-expiration-exporter/internal/checker"
	"github.com/upstars-global/domains-expiration-exporter/internal/expiration"
	"github.com/upstars-global/domains-expiration-exporter/internal/exporter"
)

const namespace = "domain-expiration-checker"

func main() {
	log, _ := zap.NewProduction()

	apiKeys := argvGetApiKeys()
	apis := make([]cf.CF, 0)
	for _, apiKey := range apiKeys {
		c, err := cf.New(apiKey)
		if err != nil {
			log.Fatal("could not create cloudflare api", zap.Error(err))
		}
		apis = append(apis, c)
	}

	manualExpirations, err := parseManualExpirations()
	if err != nil {
		log.Fatal("could not parse manual expirations", zap.Error(err))
	}

	ch := checker.New(log, apis, expiration.New(manualExpirations))
	go ch.Start()

	log.Info("registering exporter")
	prometheus.MustRegister(exporter.New(ch))

	log.Info("starting http server", zap.String("address", ":8080"))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal("failed to start http server", zap.Error(http.ListenAndServe(":8080", nil)))
}
