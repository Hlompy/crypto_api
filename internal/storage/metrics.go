package storage

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	UserCoinsGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "user_coins_total",
			Help: "Total number of user-created coins",
		},
	)

	PublishedMessages = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_published_messages_total",
			Help: "Total number of messages published to Kafka",
		},
	)
)

func InitMetrics() {
	prometheus.MustRegister(UserCoinsGauge)
	prometheus.MustRegister(PublishedMessages)

	// Стартуем HTTP endpoint для метрик
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8080", nil)
	}()
}
