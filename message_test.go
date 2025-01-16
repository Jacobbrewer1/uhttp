package uhttp

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MessageSuite struct {
	suite.Suite

	w *httptest.ResponseRecorder
}

func TestMessageSuite(t *testing.T) {
	suite.Run(t, new(MessageSuite))
}

func (s *MessageSuite) SetupTest() {
	s.w = httptest.NewRecorder()
}

func (s *MessageSuite) TestSendMessage() {
	SendMessage(s.w, "test")

	s.Equal(200, s.w.Code)

	expectedResponse := "{\"message\":\"test\"}\n"
	s.Equal(expectedResponse, s.w.Body.String())
}

func (s *MessageSuite) TestSendMessageWithStatus() {
	SendMessageWithStatus(s.w, 400, "test")

	s.Equal(400, s.w.Code)

	expectedResponse := "{\"message\":\"test\"}\n"
	s.Equal(expectedResponse, s.w.Body.String())
}
