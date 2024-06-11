package mailgun

import (
	"github.com/prometheus/client_golang/prometheus"
)

var mgMessagesCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "mailgun_messages",
		Help: "A counter for messages sent",
	},
	[]string{"status"},
)

func init() {
	prometheus.MustRegister(mgMessagesCount)
}
