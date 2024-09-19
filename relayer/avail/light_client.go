package avail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	http_client "github.com/vitwit/avail-da-module/relayer/http"
)

// AvailLightClient is a concrete implementation of the availDA interface.
// It facilitates interaction with the Avail Network by utilizing the Avail light client.
//
// Fields:
// - HttpClient: An HTTP client handler used for making requests to the Avail light client.
// - LightClientURL: The URL of the Avail light client that this module communicates with.
type LightClient struct {
	HTTPClient     *http_client.Handler
	LightClientURL string
}

func NewLightClient(lightClientURL string, httpClient *http_client.Handler) *LightClient {
	return &LightClient{
		HTTPClient:     httpClient,
		LightClientURL: lightClientURL,
	}
}

func (lc *LightClient) IsDataAvailable(data []byte, availBlockHeight int) (bool, error) {
	availBlock, err := lc.GetBlock(availBlockHeight)
	if err != nil {
		return false, err
	}

	base64CosmosBlockData := base64.StdEncoding.EncodeToString(data)

	// TODO: any better / optimized way to check if data is really available?
	return isDataIncludedInBlock(availBlock, base64CosmosBlockData), nil
}

func (lc *LightClient) GetBlock(availBlockHeight int) (Block, error) {
	// Construct the URL with the block number
	url := fmt.Sprintf("%s/v2/blocks/%v/data?fields=data", lc.LightClientURL, availBlockHeight)

	// Perform the GET request, returning the body directly
	body, err := lc.HTTPClient.Get(url)
	if err != nil {
		return Block{}, fmt.Errorf("failed to fetch data from the avail: %w", err)
	}

	// Decode the response body into the AvailBlock struct
	var block Block
	err = json.Unmarshal(body, &block)
	if err != nil {
		return Block{}, fmt.Errorf("failed to decode block response: %w", err)
	}

	return block, nil
}

// Submit sends a block of data to the light client for processing.
func (lc *LightClient) Submit(data []byte) (BlockMetaData, error) {
	var blockInfo BlockMetaData

	datab := base64.StdEncoding.EncodeToString(data)
	jsonData := []byte(fmt.Sprintf(`{"data":"%s"}`, datab))
	url := fmt.Sprintf("%s/v2/submit", lc.LightClientURL)

	// Make the POST request
	responseBody, err := lc.HTTPClient.Post(url, jsonData)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return blockInfo, err
	}

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(responseBody, &blockInfo)
	if err != nil {
		return BlockMetaData{}, err
	}

	return blockInfo, nil
}

// bruteforce comparison check
func isDataIncludedInBlock(availBlock Block, base64cosmosData string) bool {
	for _, data := range availBlock.Extrinsics {
		if data.Data == base64cosmosData {
			return true
		}
	}

	return false
}
