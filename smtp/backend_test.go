package smtp

import (
	"context"
	"net"
	"testing"

	smtp "github.com/emersion/go-smtp"
	"github.com/stretchr/testify/require"

	"github.com/teran/relay/driver"
)

func TestBackend(t *testing.T) {
	r := require.New(t)

	driverMock := driver.NewMock()
	defer driverMock.AssertExpectations(t)

	driverMock.On(
		"Send", "test@example.org", []string{"user@example.org"}, []byte("test message\r\n"),
	).Return(nil).Once()

	l := newLocalListener()

	be := NewBackend(context.Background(), driverMock)
	s := smtp.NewServer(be)
	s.AllowInsecureAuth = true

	defer s.Close()

	go s.Serve(l)

	err := sendEmail(
		l.Addr().String(),
		"test@example.org",
		"user@example.org",
		"test message",
	)
	r.NoError(err)
}

func TestBackendEmptyMessage(t *testing.T) {
	r := require.New(t)

	driverMock := driver.NewMock()
	defer driverMock.AssertExpectations(t)

	l := newLocalListener()

	be := NewBackend(context.Background(), driverMock)
	s := smtp.NewServer(be)
	s.AllowInsecureAuth = true

	defer s.Close()

	go s.Serve(l)

	err := sendEmail(
		l.Addr().String(),
		"",
		"",
		"",
	)
	r.NoError(err)
}

func newLocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	return l
}

func sendEmail(serverAddr, from, to, text string) error {
	c, err := smtp.Dial(serverAddr)
	if err != nil {
		return err
	}

	if from != "" {
		if err := c.Mail(from, nil); err != nil {
			return err
		}
	}

	if to != "" {
		if err := c.Rcpt(to, nil); err != nil {
			return err
		}
	}

	if text != "" {
		wc, err := c.Data()
		if err != nil {
			return err
		}

		if _, err := wc.Write([]byte(text)); err != nil {
			return err
		}

		err = wc.Close()
		if err != nil {
			return err
		}
	}

	err = c.Quit()
	if err != nil {
		return err
	}

	return nil
}
