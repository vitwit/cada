package relayer_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	relayer "github.com/vitwit/avail-da-module/relayer"
)

func (s *RelayerTestSuite) TestSubmitDataToAvailClient_Success() {
	data := []byte("test data")
	blocks := []int64{1, 2, 3}
	//lightClientUrl := "http://localhost:8080"

	blockInfo := relayer.BlockInfo{
		BlockHash:   "hash123",
		BlockNumber: 1,
		Hash:        "somehash",
	}

	// Create a mock HTTP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.T().Errorf("Expected POST method, got %s", r.Method)
		}

		var requestBody map[string]string
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if requestBody["data"] != base64.StdEncoding.EncodeToString(data) {
			http.Error(w, "Unexpected data", http.StatusBadRequest)
			return
		}

		responseBody, _ := json.Marshal(blockInfo)
		w.Write(responseBody)
	}))

	defer mockServer.Close()

	blockInfoResult, err := s.relayer.SubmitDataToAvailClient("seed", 1, data, blocks, mockServer.URL)

	s.Require().NoError(err)
	s.Require().Equal(blockInfo, blockInfoResult)
}

func (s *RelayerTestSuite) TestSubmitDataToAvailClient_HTTPError() {
	data := []byte("test data")
	blocks := []int64{1, 2, 3}
	//lightClientUrl := "http://localhost:8080"

	// Create a mock HTTP server that returns an error
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))

	defer mockServer.Close()

	blockInfo, err := s.relayer.SubmitDataToAvailClient("seed", 1, data, blocks, mockServer.URL)

	s.Require().Error(err)
	s.Require().Equal(relayer.BlockInfo{}, blockInfo)
}

func (s *RelayerTestSuite) TestSubmitDataToAvailClient_UnmarshalError() {
	data := []byte("test data")
	blocks := []int64{1, 2, 3}
	//lightClientUrl := "http://localhost:8080"

	// Create a mock HTTP server that returns invalid JSON
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{invalid json}`))
	}))

	defer mockServer.Close()

	blockInfo, err := s.relayer.SubmitDataToAvailClient("seed", 1, data, blocks, mockServer.URL)

	s.Require().Error(err)
	s.Require().Equal(relayer.BlockInfo{}, blockInfo)
}
