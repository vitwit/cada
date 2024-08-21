package keeper_test

import (
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/test-go/testify/require"
)

func (s *KeeperTestSuite) TestPreBlocker() {
	// Initialize the keeper, possibly with mocks.
	req := abci.RequestFinalizeBlock{ // Create a sample request for testing.
		Time:            time.Now(),
		Height:          1,
		ProposerAddress: []byte{0x1, 0x2, 0x3, 0x4},
		Txs:             [][]byte{},
	}

	// Test cases
	tests := []struct {
		name    string
		setup   func()
		req     abci.RequestFinalizeBlock
		wantErr bool
	}{
		{
			name: "successful execution",
			setup: func() {
				// Optionally, set up any additional context or mocks here.
			},
			req:     req,
			wantErr: false,
		},
		// {
		// 	name: "error during preblockerPendingBlocks",
		// 	setup: func() {
		// 		// Simulate an error in preblockerPendingBlocks by mocking the Keeper method.
		// 		s.app.AvailBlobKeeper. PreBlockerPendingBlocks = func(ctx sdk.Context, time time.Time, proposerAddr []byte, pendingBlocks *[]int64) error {
		// 			return errors.New("preblocker pending blocks error")
		// 		}
		// 	},
		// 	req:     req,
		// 	wantErr: true,
		// },
		// {
		// 	name: "handling larger tx sizes",
		// 	setup: func() {
		// 		req.Txs = append(req.Txs, []byte("larger transaction"))
		// 		req.Size = func() int { return len(req.Txs[0]) } // Simulate the size of txs.
		// 	},
		// 	req:     req,
		// 	wantErr: false,
		// },
	}

	// Run each test case
	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := s.app.AvailBlobKeeper.PreBlocker(s.sdkCtx, &tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
