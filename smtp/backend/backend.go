package backend

import (
	"context"

	smtp "github.com/emersion/go-smtp"

	"github.com/teran/relay/driver"
	"github.com/teran/relay/smtp/session"
)

type Backend = smtp.Backend

type backend struct {
	ctx context.Context
	d   driver.Driver
}

func New(ctx context.Context, d driver.Driver) Backend {
	return &backend{
		ctx: ctx,
		d:   d,
	}
}

func (b *backend) NewSession(c *smtp.Conn) (session.Session, error) {
	return session.New(b.ctx, b.d), nil
}
