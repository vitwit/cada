package relayer_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/require"
)

func (s *RelayerTestSuite) TestHTTPClientHandler_Get() {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(s.T(), "GET", r.Method)

		fmt.Fprintln(w, `{"message": "GET request successful"}`)
	}))
	defer mockServer.Close()

	resp, err := s.httpHandler.Get(mockServer.URL)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)

	expected := `{"message": "GET request successful"}`
	require.JSONEq(s.T(), expected, string(resp))
}

func (s *RelayerTestSuite) TestHTTPClientHandler_GetError() {
	_, err := s.httpHandler.Get("http://invalid-url")
	require.Error(s.T(), err)
	require.Contains(s.T(), err.Error(), "GET request error")
}

func (s *RelayerTestSuite) TestPostRequest() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Require().Equal(http.MethodPost, r.Method)
		s.Require().Equal("/post", r.URL.Path)

		s.Require().Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"response":"success"}`))
		s.Require().NoError(err)
	}))
	defer server.Close()
}
