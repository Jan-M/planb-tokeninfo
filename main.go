package main

import (
	"flag"
	"fmt"
	gometrics "github.com/rcrowley/go-metrics"
	"github.com/zalando/planb-tokeninfo/handlers/healthcheck"
	"github.com/zalando/planb-tokeninfo/handlers/metrics"
	"github.com/zalando/planb-tokeninfo/handlers/tokeninfo"
	"github.com/zalando/planb-tokeninfo/handlers/tokeninfo/jwt"
	"github.com/zalando/planb-tokeninfo/handlers/tokeninfo/proxy"
	"github.com/zalando/planb-tokeninfo/keys"
	"github.com/zalando/planb-tokeninfo/options"
	"log"
	"net/http"
	"time"
)

const (
	defaultListenAddr        = ":9021"
	defaultMetricsListenAddr = ":9020"
)

var (
	Version string = "0.0.1"
)

func init() {
	flag.Parse()
}

func setupMetrics() {
	gometrics.RegisterRuntimeMemStats(gometrics.DefaultRegistry)
	go gometrics.CaptureRuntimeMemStats(gometrics.DefaultRegistry, 60*time.Second)
	http.Handle("/metrics", metrics.Default)
	go http.ListenAndServe(defaultMetricsListenAddr, nil)
}

func main() {
	log.Printf("Started server at %v, /metrics endpoint at %v\n", defaultListenAddr, defaultMetricsListenAddr)

	setupMetrics()

	ph := tokeninfoproxy.NewTokenInfoProxyHandler(options.UpstreamTokeninfoUrl)
	kl := keys.NewCachingOpenIdProviderLoader(options.OpenIdProviderConfigurationUrl)
	jh := jwthandler.NewJwtHandler(kl)

	mux := http.NewServeMux()
	mux.Handle("/health", healthcheck.Handler(fmt.Sprintf("OK\n%s", Version)))
	mux.Handle("/oauth2/tokeninfo", tokeninfo.Handler(ph, jh))
	log.Fatal(http.ListenAndServe(defaultListenAddr, mux))
}
