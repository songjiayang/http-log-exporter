package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/songjiayang/nginx-log-exporter/collector"
	"github.com/songjiayang/nginx-log-exporter/config"
)

func main() {
	var listenAddress, configFile string
	var placeholderReplace bool

	flag.StringVar(&listenAddress, `web.listen-address`, `:9999`, `Address to listen on for the web interface and API.`)
	flag.StringVar(&configFile, `config.file`, `config.yml`, `Nginx log exporter configuration file name.`)
	flag.BoolVar(&placeholderReplace, `placeholder.replace`, false, `Enable placeholder replacement when rewriting the request path.`)
	flag.Parse()

	cfg, err := config.LoadFile(configFile)
	if err != nil {
		log.Panic(err)
	}

	var options config.Options
	options.SetPlaceholderReplace(placeholderReplace)

	for _, app := range cfg.App {
		go collector.NewCollector(app, options).Run()
	}

	fmt.Printf("running HTTP server on address %s\n", listenAddress)

	http.Handle("/metrics", promhttp.HandlerFor(
		prometheus.DefaultGatherer,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		log.Fatalf("start server with error: %v\n", err)
	}
}
