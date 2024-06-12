package mailgun

import (
	"context"
	"io"
	"strings"
	"testing"

	mg "github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMailgunDriver(t *testing.T) {
	r := require.New(t)

	mgMsg := newMgMessage("from@example.com", "test message", "test message", "to@example.org")

	mgImpl := &mailgunImplMock{}
	mgImpl.On(
		"NewMIMEMessage", []byte("test message"), []string{"to@example.org"},
	).Return(mgMsg).Once()
	mgImpl.On("Send", mgMsg).Return("<resp>", "<id>", nil).Once()

	d := New(mgImpl)
	err := d.Send(
		context.Background(),
		"from@example.org", []string{"to@example.org"}, strings.NewReader("test message"),
	)
	r.NoError(err)
}

type mailgunImplMock struct {
	mock.Mock
}

func (m *mailgunImplMock) NewMIMEMessage(body io.ReadCloser, to ...string) *mg.Message {
	data, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	args := m.Called(data, to)
	return args.Get(0).(*mg.Message)
}

func (m *mailgunImplMock) Send(_ context.Context, msg *mg.Message) (string, string, error) {
	args := m.Called(msg)
	return args.String(0), args.String(1), args.Error(2)
}

func newMgMessage(from, subject, text string, to ...string) *mg.Message {
	return (&mg.MailgunImpl{}).NewMessage(from, subject, text, to...)
}
