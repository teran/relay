package mailgun

import (
	"context"
	"io"

	mg "github.com/mailgun/mailgun-go/v4"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/teran/relay/driver"
)

type mailgun struct {
	client mg.Mailgun
}

func New(client mg.Mailgun) driver.Driver {
	return &mailgun{
		client: client,
	}
}

func (u *mailgun) Send(ctx context.Context, from string, to []string, r io.Reader) error {
	for _, recipient := range to {
		message := u.client.NewMIMEMessage(io.NopCloser(r), recipient)
		resp, id, err := u.client.Send(ctx, message)
		log.WithFields(log.Fields{
			"driver":   "mailgun",
			"id":       id,
			"response": resp,
			"error":    err,
		}).Infof("Attempting to send mail")
		if err != nil {
			mgMessagesCount.WithLabelValues("failed").Inc()
			return errors.Wrap(err, "error sending mail")
		}
		mgMessagesCount.WithLabelValues("success").Inc()
	}
	return nil
}
