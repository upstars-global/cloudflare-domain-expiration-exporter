package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/upstars-global/domains-expiration-exporter/internal/config"
	"github.com/upstars-global/domains-expiration-exporter/internal/exporter"
	"github.com/upstars-global/domains-expiration-exporter/internal/registrator"
	"github.com/upstars-global/domains-expiration-exporter/internal/registrator/domainsList"
	"github.com/upstars-global/domains-expiration-exporter/internal/registrator/godaddy"
	pn "github.com/upstars-global/domains-expiration-exporter/internal/registrator/pananames"
	"go.uber.org/zap"
	"net/http"
)

const namespace = "domain-expiration-checker"

func main() {
	log, _ := zap.NewProduction()
	defer func(log *zap.Logger) {
		_ = log.Sync()
	}(log)

	var (
		api registrator.Registrator
		err error
	)

	regAuth, err := config.New()

	if err != nil {
		log.Fatal("failed to get registrator auth info", zap.Error(err))
	}

	if regAuth.PananameAuth.Token != "" {
		api, err = pn.New(regAuth.PananameAuth.Token)
	}

	if regAuth.GodaddyAuth.Secret != "" && regAuth.GodaddyAuth.Token != "" {
		api, err = godaddy.New(regAuth.GodaddyAuth.Secret, regAuth.GodaddyAuth.Secret, "dev")
	}

	if len(regAuth.DomainsList.List) > 0 {
		api, err = domainsList.New(regAuth.DomainsList.List)
	}

	if err != nil {
		log.Fatal("Unable to create api", zap.Error(err))
	}

	if api == nil {
		log.Fatal("Unable to get registrator's auth info")
	}

	domains, err := api.GetDomains(context.Background())
	if err != nil {
		log.Fatal("Unable to get domains", zap.Error(err))
	}

	log.Info("registering exporter")
	prometheus.MustRegister(exporter.New(&domains))

	log.Info("starting http server", zap.String("address", ":8080"))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal("failed to start http server", zap.Error(http.ListenAndServe(":8080", nil)))
}
