package smtp

import (
	"context"

	smtp "github.com/emersion/go-smtp"

	"github.com/teran/relay/driver"
)

type Backend = smtp.Backend

type backend struct {
	ctx context.Context
	d   driver.Driver
}

func NewBackend(ctx context.Context, d driver.Driver) Backend {
	return &backend{
		ctx: ctx,
		d:   d,
	}
}

func (b *backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return newSession(b.ctx, b.d), nil
}
