package driver

import (
	"context"
	"io"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Send(_ context.Context, from string, to []string, body io.Reader) error {
	data, err := io.ReadAll(body)
	if err != nil {
		panic(err)
	}

	args := m.Called(from, to, data)
	return args.Error(0)
}
