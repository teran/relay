package mailgun

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"

	smtp "github.com/emersion/go-smtp"
	mg "gopkg.in/mailgun/mailgun-go.v1"
)

var _ smtp.Backend = &Backend{}
var _ smtp.User = &User{}

// Mail Backend
type Backend struct {
	Domain     string
	privateKey string
	publicKey  string
}

type User struct {
	mailgunClient mg.Mailgun
}

func NewBackend(domain, privateKey, publicKey string) (smtp.Backend, error) {
	if domain == "" || privateKey == "" || publicKey == "" {
		return nil, fmt.Errorf("domain, privateKey, publicKey must not be empty")
	}

	return &Backend{
		Domain:     domain,
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (b *Backend) Login(username, password string) (smtp.User, error) {
	return &User{
		mailgunClient: mg.NewMailgun(b.Domain, b.privateKey, b.publicKey),
	}, nil
}

func (b *Backend) AnonymousLogin() (smtp.User, error) {
	return &User{
		mailgunClient: mg.NewMailgun(b.Domain, b.privateKey, b.publicKey),
	}, nil
}

func (u *User) Send(from string, to []string, r io.Reader) error {
	m, err := mail.ReadMessage(r)
	if err != nil {
		return err
	}

	mBody, err := ioutil.ReadAll(m.Body)
	if err != nil {
		return err
	}

	for _, recipient := range to {
		message := u.mailgunClient.NewMessage(from, m.Header.Get("Subject"), string(mBody), recipient)
		resp, id, err := u.mailgunClient.Send(message)
		if err != nil {
			return err
		}
		log.Printf("ID: %s Resp: %s", id, resp)
	}
	return nil
}

func (u *User) Logout() error {
	return nil
}
