package session

import (
	"bytes"
	"context"
	"io"

	smtp "github.com/emersion/go-smtp"
	log "github.com/sirupsen/logrus"

	"github.com/teran/relay/driver"
)

var _ Session = (*session)(nil)

type Session = smtp.Session

type session struct {
	from string
	to   []string
	body []byte

	ctx context.Context
	d   driver.Driver
}

func New(ctx context.Context, d driver.Driver) Session {
	return &session{
		ctx: ctx,
		d:   d,
	}
}

func (s *session) Reset() {
}

func (s *session) Logout() error {
	log.WithFields(log.Fields{
		"from": s.from,
		"to":   s.to,
		"body": string(s.body),
	}).Debugf("on logout")
	return s.d.Send(s.ctx, s.from, s.to, bytes.NewReader(s.body))
}

func (s *session) Mail(from string, opts *smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.to = []string{to}
	return nil
}

func (s *session) Data(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	s.body = data

	return nil
}
