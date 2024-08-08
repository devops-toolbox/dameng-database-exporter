package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"dm_exporter/global"
	"dm_exporter/internal/prometheus"

	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
	//Required for debugging
	//_ "net/http/pprof"
)

func main() {
	log.AddFlags(kingpin.CommandLine)
	kingpin.Version("dm_exporter " + global.Version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	//http.Handle(*global.MetricPath,  promhttp.Handler())
	registry := prometheus.Registry()
	http.Handle(*global.MetricPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(global.LandingPage)
	})
	log.Infoln("Listening on", *global.ListenAddress)
	log.Fatal(http.ListenAndServe(*global.ListenAddress, nil))
}
