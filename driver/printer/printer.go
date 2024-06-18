package printer

import (
	"context"
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/teran/relay/driver"
)

var _ driver.Driver = (*printer)(nil)

type printer struct{}

func New() driver.Driver {
	return &printer{}
}

func (u *printer) Send(ctx context.Context, from string, to []string, r io.Reader) error {
	message, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"driver":  "printer",
		"from":    from,
		"to":      to,
		"message": string(message),
	}).Infof("Fake sending email")

	return nil
}
