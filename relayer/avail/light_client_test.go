package avail_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	avail "github.com/vitwit/avail-da-module/relayer/avail"
	httpClient "github.com/vitwit/avail-da-module/relayer/http"
)

func newMockLightClient(serverURL string) *avail.LightClient {
	httpClient := httpClient.NewHandler()
	return avail.NewLightClient(serverURL, httpClient)
}

// Test IsDataAvailable
func TestLightClient_IsDataAvailable(t *testing.T) {
	expectedBlockHeight := 100
	expectedBase64Data := base64.StdEncoding.EncodeToString([]byte("cosmos-block-data"))

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		block := avail.Block{
			Block: 100,
			Extrinsics: []avail.Extrinsics{
				{Data: expectedBase64Data},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(block)
	}))
	defer mockServer.Close()

	lightClient := newMockLightClient(mockServer.URL)

	isAvailable, err := lightClient.IsDataAvailable([]byte("cosmos-block-data"), expectedBlockHeight)
	assert.NoError(t, err, "expected no error")
	assert.True(t, isAvailable, "expected data to be available")
}

func TestLightClient_GetBlock(t *testing.T) {
	expectedBlockHeight := 100
	expectedExtrinsics := []avail.Extrinsics{
		{Data: base64.StdEncoding.EncodeToString([]byte("cosmos-block-data"))},
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		block := avail.Block{
			Block:      int64(expectedBlockHeight),
			Extrinsics: expectedExtrinsics,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(block)
	}))
	defer mockServer.Close()

	lightClient := newMockLightClient(mockServer.URL)

	block, err := lightClient.GetBlock(expectedBlockHeight)
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, expectedBlockHeight, int(block.Block), "expected block height to match")
	assert.Equal(t, expectedExtrinsics, block.Extrinsics, "expected extrinsics to match")
}

func TestLightClient_Submit(t *testing.T) {
	expectedBlockMetaData := avail.BlockMetaData{
		BlockNumber: 100,
		BlockHash:   "testhash",
		Hash:        "datahash",
		Index:       1,
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedBlockMetaData)
	}))
	defer mockServer.Close()

	lightClient := newMockLightClient(mockServer.URL)

	blockMeta, err := lightClient.Submit([]byte("testdata"))
	assert.NoError(t, err, "expected no error")
	assert.Equal(t, expectedBlockMetaData, blockMeta, "expected block metadata to match")
}
