package mailgun

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"strings"

	smtp "github.com/emersion/go-smtp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/html"
	mg "gopkg.in/mailgun/mailgun-go.v1"
)

var _ smtp.Backend = &Backend{}
var _ smtp.User = &User{}

// Backend type
type Backend struct {
	Domain                 string
	privateKey             string
	publicKey              string
	metricsMailgunMessages *prometheus.CounterVec
}

// User type
type User struct {
	mailgunClient          mg.Mailgun
	metricsMailgunMessages *prometheus.CounterVec
}

// NewBackend returns new instance of backend
func NewBackend(domain, privateKey, publicKey string) (smtp.Backend, error) {
	if domain == "" || privateKey == "" || publicKey == "" {
		return nil, fmt.Errorf("domain, privateKey, publicKey must not be empty")
	}

	b := &Backend{
		Domain:     domain,
		privateKey: privateKey,
		publicKey:  publicKey,
		metricsMailgunMessages: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "mailgun_messages",
				Help: "A counter for messages sent",
			},
			[]string{"status"},
		),
	}

	prometheus.MustRegister(b.metricsMailgunMessages)

	return b, nil
}

// Login is used to authenticate the user
// In relay there's no need for that at the moment
func (b *Backend) Login(username, password string) (smtp.User, error) {
	return &User{
		mailgunClient:          mg.NewMailgun(b.Domain, b.privateKey, b.publicKey),
		metricsMailgunMessages: b.metricsMailgunMessages,
	}, nil
}

// AnonymousLogin returns anonymouse user object
func (b *Backend) AnonymousLogin() (smtp.User, error) {
	return &User{
		mailgunClient:          mg.NewMailgun(b.Domain, b.privateKey, b.publicKey),
		metricsMailgunMessages: b.metricsMailgunMessages,
	}, nil
}

func (b *Backend) ListenAndServeMetrics(addr string) error {
	s := http.Server{
		Addr:    addr,
		Handler: promhttp.Handler(),
	}
	return s.ListenAndServe()
}

// Send will send email synchronously via Mailgun service
func (u *User) Send(from string, to []string, r io.Reader) error {
	m, err := mail.ReadMessage(r)
	if err != nil {
		return err
	}

	mBody, err := ioutil.ReadAll(m.Body)
	if err != nil {
		return err
	}

	mBodyString := string(mBody)
	var (
		mBodyText string
		mBodyHTML string
	)

	_, err = html.Parse(strings.NewReader(mBodyString))
	if err == nil {
		mBodyHTML = mBodyString
	} else {
		mBodyText = mBodyString
	}

	for _, recipient := range to {
		message := u.mailgunClient.NewMessage(from, m.Header.Get("Subject"), mBodyText, recipient)
		message.SetHtml(mBodyHTML)
		resp, id, err := u.mailgunClient.Send(message)
		if err != nil {
			u.metricsMailgunMessages.WithLabelValues("fail").Inc()
			return err
		}
		u.metricsMailgunMessages.WithLabelValues("success").Inc()
		log.Printf("ID: %s Resp: %s", id, resp)
	}
	return nil
}

// Logout is called after all operations are complete within the session
// Here in relay there's no need to implement anything special for that
func (u *User) Logout() error {
	return nil
}
