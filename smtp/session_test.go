package smtp

import (
	"testing"

	"github.com/teran/relay/driver"

	smtp "github.com/emersion/go-smtp"
	"github.com/stretchr/testify/suite"
)

func (s *sessionTestSuite) TestSendMail() {
	err := s.session.Rcpt("some-addr1", nil)
	s.Require().NoError(err)

	err = s.session.Rcpt("some-addr2", nil)
	s.Require().NoError(err)

	err = s.session.Mail("from-addr", nil)
	s.Require().NoError(err)

	s.driver.On("Send", "from-addr", []string{"some-addr1", "some-addr2"}, []byte{}).Return(nil).Once()

	err = s.session.Logout()
	s.NoError(err)
}

// ========================================================================
// Test suite setup
// ========================================================================
type sessionTestSuite struct {
	suite.Suite

	driver  *driver.Mock
	session smtp.Session
}

func (s *sessionTestSuite) SetupTest() {
	s.driver = driver.NewMock()
	s.session = newSession(s.T().Context(), s.driver)
}

func TestSessionTestSuite(t *testing.T) {
	suite.Run(t, new(sessionTestSuite))
}
