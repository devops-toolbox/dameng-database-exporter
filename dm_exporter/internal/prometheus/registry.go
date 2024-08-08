package prometheus

import (
	"dm_exporter/global"
	"errors"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	client "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

func Registry() *client.Registry {
	log.Infoln("Starting dm_exporter " + global.Version)
	dsn := os.Getenv("DATA_SOURCE_NAME")
	// Load default metrics
	if _, err := toml.DecodeFile(*global.DefaultFileMetrics, &MetricsToScrap); err != nil {
		log.Errorln(err)
		panic(errors.New("Error while loading " + *global.DefaultFileMetrics))
	} else {
		log.Infoln("Successfully loaded default metrics from: " + *global.DefaultFileMetrics)
	}

	// If custom metrics, load it
	if strings.Compare(*global.CustomMetrics, "") != 0 {
		if _, err := toml.DecodeFile(*global.CustomMetrics, &AdditionalMetrics); err != nil {
			log.Errorln(err)
			panic(errors.New("Error while loading " + *global.CustomMetrics))
		} else {
			log.Infoln("Successfully loaded custom metrics from: " + *global.CustomMetrics)
		}

		MetricsToScrap.Metric = append(MetricsToScrap.Metric, AdditionalMetrics.Metric...)
	} else {
		log.Infoln("No custom metrics defined.")
	}

	exporter := NewExporter(dsn)
	log.Infoln("No custom metrics defined.")
	//prometheus.MustRegister(exporter)
	registry := client.NewRegistry()
	registry.MustRegister(exporter)
	return registry
}
