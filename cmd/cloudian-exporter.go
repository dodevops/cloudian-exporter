package main

import (
	"cloudian-exporter/internal"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"time"
)
import "github.com/sirupsen/logrus"

func main() {
	logLevel := os.Getenv("EXPORTER_LOGLEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	if l, err := logrus.ParseLevel(logLevel); err != nil {
		log.Fatalf("Can not parse loglevel %s: %v", logLevel, err)
	} else {
		logrus.SetLevel(l)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})

	baseURL := os.Getenv("CLOUDIAN_URL")
	username := os.Getenv("CLOUDIAN_USERNAME")
	password := os.Getenv("CLOUDIAN_PASSWORD")

	if baseURL == "" {
		logrus.Fatal("Missing base url to Cloudian. Please set CLOUDIAN_URL")
	}
	if username == "" {
		logrus.Fatal("Missing username to Cloudian. Please set CLOUDIAN_USERNAME")
	}
	if password == "" {
		logrus.Fatal("Missing username to Cloudian. Please set CLOUDIAN_PASSWORD")
	}
	api := internal.NewCloudianAPI(baseURL, username, password)

	refresh := 5 * time.Minute
	if r, ok := os.LookupEnv("EXPORTER_REFRESH"); ok {
		if pr, err := time.ParseDuration(fmt.Sprintf("%sm", r)); err != nil {
			logrus.Fatal("Can not parse refresh value from EXPORTER_REFRESH (%s): %v", r, err)
		} else {
			refresh = pr
		}
	}

	reg := prometheus.NewRegistry()
	exit := make(chan bool)

	exporter := internal.NewCloudianExporter(refresh, api, exit, reg)
	logrus.Info("Starting cloudian exporter")
	exporter.Run()

	listen := os.Getenv("EXPORTER_LISTEN")
	if listen == "" {
		listen = ":8080"
	}

	logrus.Info("Starting metrics server")
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	if err := http.ListenAndServe(listen, nil); err != nil {
		logrus.Fatal("Can not start metrics server: %v", err)
	}
}
