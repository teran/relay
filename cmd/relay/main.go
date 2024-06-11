package main

import (
	"log"

	smtp "github.com/emersion/go-smtp"
	"github.com/kelseyhightower/envconfig"
	mg "github.com/mailgun/mailgun-go/v4"

	"github.com/teran/relay/backend/mailgun"
)

type config struct {
	Addr              string `default:":25"`
	AllowInsecureAuth bool   `default:"false" envconfig:"ALLOW_INSECURE_AUTH"`
	AuthDisabled      bool   `default:"false" envconfig:"AUTH_DISABLED"`
	Domain            string `required:"true"`
	MailgunPrivateKey string `required:"true" envconfig:"MAILGUN_PRIVATE_KEY"`
	MailgunPublicKey  string `required:"true" envconfig:"MAILGUN_PUBLIC_KEY"`
	MaxIdleSeconds    int    `default:"300" envconfig:"MAX_IDLE_SECONDS"`
	MaxMessageBytes   int    `default:"1048576" envconfig:"MAX_MESSAGE_BYTES"`
	MaxRecipients     int    `default:"50" envconfig:"MAX_RECIPIENTS"`
	MetricsAddr       string `default:":8081" envconfig:"METRICS_ADDR"`
}

func main() {
	var cfg config
	err := envconfig.Process("RELAY", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	mgDriver, err := mg.NewMailgunFromEnv()
	if err != nil {
		panic(err)
	}
	be, err := mailgun.NewBackend(mgDriver)
	if err != nil {
		log.Fatal(err)
	}

	s := smtp.NewServer(be)

	s.Addr = cfg.Addr
	s.Domain = cfg.Domain
	s.MaxIdleSeconds = cfg.MaxIdleSeconds
	s.MaxMessageBytes = cfg.MaxMessageBytes
	s.MaxRecipients = cfg.MaxRecipients
	s.AuthDisabled = cfg.AuthDisabled
	s.AllowInsecureAuth = cfg.AllowInsecureAuth

	go func(metricsAddr string) {
		log.Fatal(be.(*mailgun.Backend).ListenAndServeMetrics(metricsAddr))
	}(cfg.MetricsAddr)

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
