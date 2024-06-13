package main

import (
	"context"
	"net/http"
	"time"

	smtp "github.com/emersion/go-smtp"
	"github.com/kelseyhightower/envconfig"
	mg "github.com/mailgun/mailgun-go/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/teran/relay/driver/mailgun"
	smtpWrapper "github.com/teran/relay/smtp"
)

var (
	appVersion     = "n/a (dev build)"
	buildTimestamp = "undefined"
)

type config struct {
	LogLevel          log.Level `envconfig:"LOG_LEVEL" default:"INFO"`
	Addr              string    `envconfig:"ADDR" default:":25"`
	AllowInsecureAuth bool      `envconfig:"ALLOW_INSECURE_AUTH"`
	AuthDisabled      bool      `envconfig:"AUTH_DISABLED"`
	Domain            string    `envconfig:"DOMAIN" required:"true"`
	MailgunAPIKey     string    `envconfig:"MAILGUN_API_KEY" required:"true"`
	MailgunURL        string    `envconfig:"MAILGUN_URL" required:"true"`
	MaxMessageBytes   int64     `default:"1048576" envconfig:"MAX_MESSAGE_BYTES"`
	MaxRecipients     int       `default:"50" envconfig:"MAX_RECIPIENTS"`
	MetricsAddr       string    `envconfig:"METRICS_ADDR" default:":8081" `
}

func main() {
	var cfg config
	envconfig.MustProcess("RELAY", &cfg)

	ctx := context.Background()

	log.SetLevel(cfg.LogLevel)

	lf := new(log.TextFormatter)
	lf.FullTimestamp = true
	log.SetFormatter(lf)

	log.WithFields(log.Fields{
		"app_version":     appVersion,
		"build_timestamp": buildTimestamp,
	}).Infof("initializing application")

	mgCli := mg.NewMailgun(cfg.Domain, cfg.MailgunAPIKey)
	mgCli.SetAPIBase(cfg.MailgunURL)

	dr := mailgun.New(mgCli)
	be := smtpWrapper.NewBackend(ctx, dr)

	s := smtp.NewServer(be)

	s.Addr = cfg.Addr
	s.Domain = cfg.Domain
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = cfg.MaxMessageBytes
	s.MaxRecipients = cfg.MaxRecipients
	s.AllowInsecureAuth = cfg.AllowInsecureAuth

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.ListenAndServe()
	})

	g.Go(func() error {
		http.Handle("/metrics", promhttp.Handler())
		return http.ListenAndServe(cfg.MetricsAddr, nil)
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
