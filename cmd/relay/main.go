package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	smtp "github.com/emersion/go-smtp"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/kelseyhightower/envconfig"
	mg "github.com/mailgun/mailgun-go/v4"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/teran/relay/driver"
	"github.com/teran/relay/driver/mailgun"
	"github.com/teran/relay/driver/printer"
	smtpWrapper "github.com/teran/relay/smtp"
)

var (
	appVersion     = "n/a (dev build)"
	buildTimestamp = "undefined"
)

type config struct {
	LogLevel          log.Level `envconfig:"LOG_LEVEL" default:"INFO"`
	Addr              string    `envconfig:"ADDR" default:":25"`
	EnableTLS         bool      `envconfig:"ENABLE_TLS" default:"false"`
	TLSAddr           string    `envconfig:"TLS_ADDR" default:":465"`
	TLSCertificate    string    `envconfig:"TLS_CERTIFICATE"`
	TLSKey            string    `envconfig:"TLS_KEY"`
	Driver            string    `envconfig:"DRIVER" required:"true"`
	AllowInsecureAuth bool      `envconfig:"ALLOW_INSECURE_AUTH"`
	AuthDisabled      bool      `envconfig:"AUTH_DISABLED"`
	Domain            string    `envconfig:"DOMAIN" required:"true"`
	MailgunAPIKey     string    `envconfig:"MAILGUN_API_KEY"`
	MailgunURL        string    `envconfig:"MAILGUN_URL"`
	MaxMessageBytes   int64     `envconfig:"MAX_MESSAGE_BYTES" default:"1048576"`
	MaxRecipients     int       `envconfig:"MAX_RECIPIENTS" default:"5"`
	MetricsAddr       string    `envconfig:"METRICS_ADDR" default:":8081" `
}

func (c config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Driver, validation.Required, validation.In("mailgun", "printer")),
		validation.Field(&c.Domain, validation.Required, is.Domain),
		validation.Field(&c.MailgunAPIKey, validation.When(c.Driver == "mailgun", validation.Required)),
		validation.Field(&c.MailgunURL, validation.When(c.Driver == "mailgun", validation.Required, is.URL)),
	)
}

func main() {
	var cfg config
	envconfig.MustProcess("RELAY", &cfg)

	if err := cfg.Validate(); err != nil {
		panic(err)
	}

	ctx := context.Background()

	log.SetLevel(cfg.LogLevel)

	lf := new(log.TextFormatter)
	lf.FullTimestamp = true
	log.SetFormatter(lf)

	log.WithFields(log.Fields{
		"app_version":     appVersion,
		"build_timestamp": buildTimestamp,
	}).Infof("initializing application")

	dr, err := newDriver(cfg)
	if err != nil {
		panic(err)
	}

	be := smtpWrapper.NewBackend(ctx, dr)

	g, _ := errgroup.WithContext(ctx)

	log.WithFields(log.Fields{
		"addr":   cfg.Addr,
		"domain": cfg.Domain,
	}).Trace("initializing SMTP server ...")

	s, err := newServer(cfg, be)
	if err != nil {
		panic(err)
	}

	g.Go(func() error {
		return s.ListenAndServe()
	})

	if cfg.EnableTLS {
		log.WithFields(log.Fields{
			"addr":   cfg.TLSAddr,
			"domain": cfg.Domain,
		}).Trace("initializing SMTPS server ...")

		ts, err := newTLSServer(cfg, be)
		if err != nil {
			panic(err)
		}

		g.Go(func() error {
			return ts.ListenAndServeTLS()
		})
	}

	g.Go(func() error {
		http.Handle("/metrics", promhttp.Handler())
		return http.ListenAndServe(cfg.MetricsAddr, nil)
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}
}

func newDriver(cfg config) (driver.Driver, error) {
	log.Tracef("mailing driver `%s` requested ...", cfg.Driver)

	switch cfg.Driver {
	case "mailgun":
		mgCli := mg.NewMailgun(cfg.Domain, cfg.MailgunAPIKey)
		mgCli.SetAPIBase(cfg.MailgunURL)

		return mailgun.New(mgCli), nil
	case "printer":
		return printer.New(), nil
	default:
		return nil, errors.Errorf("unexpected driver: `%s`", cfg.Driver)
	}
}

func newServer(cfg config, be smtp.Backend) (*smtp.Server, error) {
	s := smtp.NewServer(be)

	s.Addr = cfg.Addr
	s.Domain = cfg.Domain
	s.MaxMessageBytes = cfg.MaxMessageBytes
	s.MaxRecipients = cfg.MaxRecipients
	s.AllowInsecureAuth = cfg.AllowInsecureAuth
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second

	return s, nil
}

func newTLSServer(cfg config, be smtp.Backend) (*smtp.Server, error) {
	keypair, err := tls.X509KeyPair([]byte(cfg.TLSCertificate), []byte(cfg.TLSKey))
	if err != nil {
		return nil, err
	}

	s, err := newServer(cfg, be)
	if err != nil {
		return nil, err
	}

	s.Addr = cfg.TLSAddr
	s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{keypair}}

	return s, nil
}
