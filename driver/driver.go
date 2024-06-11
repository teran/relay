package driver

import (
	"context"
	"io"
)

type Driver interface {
	Send(ctx context.Context, from string, to []string, body io.Reader) error
}
