package prometheus

import (
	"database/sql"
	"dm_exporter/global"
	"dm_exporter/internal/dameng"
	"strings"
	"sync"
	"time"

	client "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

// Exporter collects DmService DB metrics. It implements client.Collector.
type Exporter struct {
	dsn             string
	duration, error client.Gauge
	totalScrapes    client.Counter
	scrapeErrors    *client.CounterVec
	up              client.Gauge
	db              *sql.DB
}

// NewExporter returns a new DmService DB global.Exporter for the provided DSN.
func NewExporter(dsn string) *Exporter {
	db := dameng.Connect(dsn)
	return &Exporter{
		dsn: dsn,
		duration: client.NewGauge(client.GaugeOpts{
			Namespace: global.Namespace,
			Subsystem: global.Exporter,
			Name:      "last_scrape_duration_seconds",
			Help:      "Duration of the last scrape of metrics from DM DB.",
		}),
		totalScrapes: client.NewCounter(client.CounterOpts{
			Namespace: global.Namespace,
			Subsystem: global.Exporter,
			Name:      "scrapes_total",
			Help:      "Total number of times DM DB was scraped for metrics.",
		}),
		scrapeErrors: client.NewCounterVec(client.CounterOpts{
			Namespace: global.Namespace,
			Subsystem: global.Exporter,
			Name:      "scrape_errors_total",
			Help:      "Total number of times an error occured scraping a DM database.",
		}, []string{"collector"}),
		error: client.NewGauge(client.GaugeOpts{
			Namespace: global.Namespace,
			Subsystem: global.Exporter,
			Name:      "last_scrape_error",
			Help:      "Whether the last scrape of metrics from DM DB resulted in an error (1 for error, 0 for success).",
		}),
		up: client.NewGauge(client.GaugeOpts{
			Namespace: global.Namespace,
			Name:      "up",
			Help:      "Whether the DM database server is up.",
		}),
		db: db,
	}
}

// Describe describes all the metrics exported by the DM DB global.Exporter.
func (e *Exporter) Describe(ch chan<- *client.Desc) {
	metricCh := make(chan client.Metric)
	doneCh := make(chan struct{})

	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()

	e.Collect(metricCh)
	close(metricCh)
	<-doneCh

}

// Collect implements client.Collector.
func (e *Exporter) Collect(ch chan<- client.Metric) {
	e.scrape(ch)
	ch <- e.duration
	ch <- e.totalScrapes
	ch <- e.error
	e.scrapeErrors.Collect(ch)
	ch <- e.up
}

func (e *Exporter) scrape(ch chan<- client.Metric) {
	e.totalScrapes.Inc()
	var err error
	defer func(begun time.Time) {
		e.duration.Set(time.Since(begun).Seconds())
		if err == nil {
			e.error.Set(0)
		} else {
			e.error.Set(1)
		}
	}(time.Now())

	if err = e.db.Ping(); err != nil {
		if strings.Contains(err.Error(), "sql: database is closed") {
			log.Infoln("Reconnecting to DB")
			e.db = dameng.Connect(e.dsn)
		}
	}
	if err = e.db.Ping(); err != nil {
		log.Errorln("Error pinging dm db:", err)
		//e.db.Close()
		e.up.Set(0)
		return
	} else {
		log.Debugln("Successfully pinged DM database: ")
		e.up.Set(1)
	}

	wg := sync.WaitGroup{}
	for _, metric := range MetricsToScrap.Metric {
		wg.Add(1)
		metric := metric //https://golang.org/doc/faq#closures_and_goroutines

		go func() {
			defer wg.Done()

			log.Debugln("About to scrape metric: ")
			log.Debugln("- Metric MetricsDesc: ", metric.MetricsDesc)
			log.Debugln("- Metric Context: ", metric.Context)
			log.Debugln("- Metric MetricsType: ", metric.MetricsType)
			log.Debugln("- Metric Labels: ", metric.Labels)
			log.Debugln("- Metric FieldToAppend: ", metric.FieldToAppend)
			log.Debugln("- Metric IgnoreZeroResult: ", metric.IgnoreZeroResult)
			log.Debugln("- Metric Request: ", metric.Request)

			if len(metric.Request) == 0 {
				log.Errorln("Error scraping for ", metric.MetricsDesc, ". Did you forget to define request in your toml file?")
			}

			if len(metric.MetricsDesc) == 0 {
				log.Errorln("Error scraping for query", metric.Request, ". Did you forget to define metricsdesc  in your toml file?")
			}

			if err = ScrapeMetric(e.db, ch, metric); err != nil {
				log.Errorln("Error scraping for", metric.Context, "_", metric.MetricsDesc, ":", err)
				e.scrapeErrors.WithLabelValues(metric.Context).Inc()
			} else {
				log.Debugln("Successfully scrapped metric: ", metric.Context)
			}
		}()
	}
	wg.Wait()
}
